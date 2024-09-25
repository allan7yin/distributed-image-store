package s3

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/wire"
)

// for wiring up project with DI

// NewS3Client creates an S3 client from the AWS configuration
func NewS3Client() (*s3.Client, error) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	// Create the S3 client
	return s3.NewFromConfig(cfg), nil
}

// NewS3Handler creates a new S3Handler with the S3 client
func NewS3Handler(client *s3.Client) *Handler {
	return &Handler{
		Client: client,
	}
}

// ProviderSet set for wire
var ProviderSet = wire.NewSet(NewS3Client, NewS3Handler)
