package commands

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/docker/docker/client"

	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull a Docker image from S3 compatible storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		image, _ := cmd.Flags().GetString("image")
		bucket, _ := cmd.Flags().GetString("bucket")

		if image == "" || bucket == "" {
			return fmt.Errorf("both image and bucket flags must be provided")
		}

		return pullImageFromS3(context.Background(), image, bucket)
	},
}

func init() {
	PullCmd.Flags().StringP("image", "i", "", "Name of the image to download")
	PullCmd.Flags().StringP("bucket", "b", "", "Name of the S3 bucket")
}

func pullImageFromS3(ctx context.Context, image, bucket string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	if endpoint, ok := os.LookupEnv("AWS_S3_ENDPOINT"); ok {
		// Use a custom endpoint resolver if AWS_S3_ENDPOINT is set
		cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
			return aws.Endpoint{
				URL:               endpoint,
				HostnameImmutable: true,
				SigningRegion:     region,
			}, nil
		})
	}

	s3Client := s3.NewFromConfig(cfg)

	// Extract the image name and digest from the input
	imageName := strings.SplitN(image, ":", 4)[0]
	digest := strings.SplitN(image, ":", 4)[3]
	newTag := fmt.Sprintf("%s:%s", imageName, digest)

	// Download the image from S3
	downloader := manager.NewDownloader(s3Client)
	tmpFile, err := os.CreateTemp("", "image-")
	if err != nil {
		return fmt.Errorf("failed to create temporary file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	n, err := downloader.Download(ctx, tmpFile, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(image),
	})
	if err != nil {
		return fmt.Errorf("failed to download image: %v", err)
	}
	fmt.Printf("Downloaded %s (%d bytes)\n", image, n)

	// Load the image into Docker
	dockerClient, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %v", err)
	}
	defer dockerClient.Close()

	imageBytes, err := ioutil.ReadFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read image contents: %v", err)
	}

	imageLoadResponse, err := dockerClient.ImageLoad(ctx, bytes.NewReader(imageBytes), false)
	if err != nil {
		return fmt.Errorf("failed to load image: %v", err)
	}
	defer imageLoadResponse.Body.Close()

	// Tag the image with the new tag
	dockerClient.ImageTag(ctx, image, newTag)
	fmt.Printf("Tagged image %s as %s\n", image, newTag)

	fmt.Printf("Loaded image %s\n", newTag)
	return nil
}
