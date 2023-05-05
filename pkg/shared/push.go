package shared

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push Docker image to S3 compatible storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		image, _ := cmd.Flags().GetString("image")
		bucket, _ := cmd.Flags().GetString("bucket")

		if image == "" || bucket == "" {
			return fmt.Errorf("both image and bucket flags must be provided")
		}

		return pushImageToS3(context.Background(), image, bucket)
	},
}

func init() {
	PushCmd.Flags().StringP("image", "i", "", "Name of the image to upload")
	PushCmd.Flags().StringP("bucket", "b", "", "Name of the S3 bucket")
}

func pushImageToS3(ctx context.Context, image, bucket string) error {
    cfg, err := config.LoadDefaultConfig(ctx)
    if err != nil {
        return fmt.Errorf("unable to load AWS SDK config: %v", err)
    }

    if endpoint, ok := os.LookupEnv("AWS_S3_ENDPOINT"); ok {
        // Use a custom endpoint resolver if AWS_S3_ENDPOINT is set
        cfg.EndpointResolver = aws.EndpointResolverFunc(func(service, region string) (aws.Endpoint, error) {
            return aws.Endpoint{
                URL:           endpoint,
                HostnameImmutable: true,
                SigningRegion: region,
            }, nil
        })
    }

    s3Client := s3.NewFromConfig(cfg)

    // Create a Docker client
    dockerClient, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        return fmt.Errorf("failed to create Docker client: %v", err)
    }
    defer dockerClient.Close()

    // Save the image to a tarball
    imageTarball := fmt.Sprintf("%s.tar", strings.ReplaceAll(image, ":", "_"))
    imageReader, err := dockerClient.ImageSave(ctx, []string{image})
    if err != nil {
        return fmt.Errorf("failed to save Docker image: %v", err)
    }
    defer imageReader.Close()

    // Upload the image to S3
    uploader := manager.NewUploader(s3Client)
    result, err := uploader.Upload(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(imageTarball),
        Body:   imageReader,
    })
    if err != nil {
        return fmt.Errorf("failed to upload image: %v", err)
    }

    fmt.Printf("Successfully uploaded %s to %s\n", image, result.Location)

    return nil
}
