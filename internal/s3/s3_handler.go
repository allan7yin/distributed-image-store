package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"os"
	"time"
)

type Handler struct {
	Client *s3.Client
}

func (handler *Handler) GeneratePresignedURL(expiry time.Duration) (string, uuid.UUID, error) {
	presigner := s3.NewPresignClient(handler.Client)

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
