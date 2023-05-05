package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: %s <image_name:tag> <bucket>", os.Args[0])
	}

	imageName := os.Args[1]
	bucket := os.Args[2]

	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	endpoint := os.Getenv("AWS_S3_ENDPOINT")
	region := os.Getenv("AWS_REGION")

	if accessKey == "" || secretKey == "" || endpoint == "" || region == "" {
		log.Fatalf("Please set AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_S3_ENDPOINT, and AWS_REGION environment variables")
	}

	// Save the image to a tarball
	imageTarball := fmt.Sprintf("%s.tar", strings.ReplaceAll(imageName, ":", "_"))
	cmd := exec.Command("docker", "save", "-o", imageTarball, imageName)
	err := cmd.Run()
	if err != nil {
		log.Fatalf("Failed to save Docker image: %v", err)
	}
	defer os.Remove(imageTarball)

	creds := credentials.NewStaticCredentials(accessKey, secretKey, "")
	_, err = creds.Get()
	if err != nil {
		log.Fatalf("Failed to load credentials: %v", err)
	}

	awsConfig := &aws.Config{
		Region:           aws.String(region),
		Endpoint:         aws.String(endpoint),
		S3ForcePathStyle: aws.Bool(true),
		Credentials:      creds,
	}

	sess := session.Must(session.NewSession(awsConfig))

	uploader := s3manager.NewUploader(sess)
	file, err := os.Open(imageTarball)
	if err != nil {
		log.Fatalf("Failed to open image tarball: %v", err)
	}
	defer file.Close()

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(imageTarball),
		Body:   file,
	})
	if err != nil {
		log.Fatalf("Failed to upload image: %v", err)
	}

	fmt.Printf("Successfully uploaded %s to %s\n", imageName, result.Location)
}
