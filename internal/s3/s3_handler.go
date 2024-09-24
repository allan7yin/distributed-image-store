package internal

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"time"
)

type S3Handler struct {
	Client *s3.Client
}

// NewS3Handler creates a new S3Handler with a configured S3 client
func NewS3Handler() (*S3Handler, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(os.Getenv("AWS_REGION")),
	)
	if err != nil {
		log.Fatalf("Unable to load AWS config: %v", err)
		return nil, err
	}

	s3Client := s3.NewFromConfig(cfg)
	return &S3Handler{
		Client: s3Client,
	}, nil
}

// GeneratePresignedURL generates a presigned URL for an S3 object
func (handler *S3Handler) GeneratePresignedURL(bucket, key string, expiry time.Duration) (string, error) {
	presigner := s3.NewPresignClient(handler.Client)

	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	presignedURL, err := presigner.PresignGetObject(context.TODO(), getObjectInput, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", fmt.Errorf("error presigning request: %w", err)
	}

	return presignedURL.URL, nil
}
