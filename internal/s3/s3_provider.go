package s3

import (
	"bit-image/pkg/storage"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/wire"
)

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

// NewS3FileSystem creates a new S3FileSystem with the S3 client
func NewS3FileSystem(s3Client *s3.Client) *storage.S3FileSystem {
	return storage.NewS3FileSystem(s3Client, nil) // Assuming s3TransferManager is not required or is nil for now
}

// NewHandler creates a new Handler with the S3FileSystem
func NewHandler(fileSystem *storage.S3FileSystem) *Handler {
	return &Handler{
		FileSystem: fileSystem,
	}
}

// ProviderSet set for wire
var ProviderSet = wire.NewSet(NewS3Client, NewS3FileSystem, NewHandler)
