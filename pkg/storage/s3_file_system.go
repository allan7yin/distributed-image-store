package storage

import (
	"bit-image/pkg/common"
	"context"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/google/uuid"
	"os"
	"sync"
	"time"
)

type S3FileSystem struct {
	s3Client          *s3.Client
	s3TransferManager *manager.Uploader
	resourceLocks     sync.Map
	urlTTL            int
	lockStripeCount   int
}

func NewS3FileSystem(s3Client *s3.Client, s3TransferManager *manager.Uploader) *S3FileSystem {
	return &S3FileSystem{
		s3Client:          s3Client,
		s3TransferManager: s3TransferManager,
		urlTTL:            120000, // TTL in milliseconds
		lockStripeCount:   10,
	}
}

func (fs *S3FileSystem) createFolder(folderName string) error {
	_, err := fs.s3Client.CreateBucket(context.TODO(), &s3.CreateBucketInput{
		Bucket: aws.String(folderName),
	})
	if err != nil {
		fmt.Printf("Failed to create bucket: %v\n", err)
		return err
	}
	return nil
}

func (fs *S3FileSystem) bucketExists(bucketName string) (bool, error) {
	_, err := fs.s3Client.HeadBucket(context.TODO(), &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "404" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fs *S3FileSystem) fileExists(bucketName string, destKey string) (bool, error) {
	_, err := fs.s3Client.HeadObject(context.TODO(), &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(destKey),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) && apiErr.ErrorCode() == "404" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (fs *S3FileSystem) GeneratePresignedURL(expiry time.Duration) (string, uuid.UUID, error) {
	presigner := s3.NewPresignClient(fs.s3Client)

	imageId := uuid.New()
	imageIdString := imageId.String()
	defaultBucketName := os.Getenv("DEFAULT_BUCKET_NAME")

	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(defaultBucketName),
		Key:    aws.String(imageIdString),
	}

	presignedURL, err := presigner.PresignPutObject(context.TODO(), putObjectInput, s3.WithPresignExpires(expiry))
	if err != nil {
		return "", uuid.UUID{}, fmt.Errorf("error presigning request: %w", err)
	}

	return presignedURL.URL, imageId, nil
}

func (fs *S3FileSystem) MoveFileToFolder(file common.File, srcFolderName, destFolderName string) error {
	// Assuming `file` has a `Name` field that represents the file name
	srcKey := fmt.Sprintf("%s/%s", srcFolderName, file.Id)
	destKey := fmt.Sprintf("%s/%s", destFolderName, file.Id)
	bucket := os.Getenv("DEFAULT_BUCKET_NAME")

	ctx := context.TODO()

	found, _ := fs.fileExists(bucket, destKey)
	if found {
		return fmt.Errorf("Cannot move file to " + destFolderName + ". Already exists object with id: " + file.Id)
	}

	// Step 1: Copy the object to the new location
	_, err := fs.s3Client.CopyObject(ctx, &s3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(bucket + "/" + srcKey), // Source should be in the format "bucket/key"
		Key:        aws.String(destKey),
		ACL:        types.ObjectCannedACLPrivate, // Or other ACL as needed
	})
	if err != nil {
		return fmt.Errorf("failed to copy object: %w", err)
	}

	// Step 2: Delete the original object
	_, err = fs.s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete original object: %w", err)
	}

	return nil
}

func (fs *S3FileSystem) moveFilesToFolder(files []common.File, srcFolderName, destFolderName string) error {
	for _, file := range files {
		err := fs.MoveFileToFolder(file, srcFolderName, destFolderName)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fs *S3FileSystem) GetObjectMetadata(key, bucket string) (int64, string, error) {
	ctx := context.TODO()

	result, err := fs.s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return 0, "", fmt.Errorf("failed to get metadata for object %s in bucket %s: %w", key, bucket, err)
	}

	// Extract size and content type from the metadata
	size := result.ContentLength
	contentType := aws.ToString(result.ContentType)

	return *size, contentType, nil
}

func (fs *S3FileSystem) deleteFilesFromFolder(fileIDs []string, folderName string) error {
	// Placeholder: Implement logic for deleting files from a folder
	return nil
}
