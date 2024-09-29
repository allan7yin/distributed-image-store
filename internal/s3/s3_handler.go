package s3

import (
	"bit-image/pkg/common"
	"bit-image/pkg/storage"
	"github.com/google/uuid"
	"time"
)

type Handler struct {
	FileSystem *storage.S3FileSystem
}

func (handler *Handler) GeneratePresignedURL(expiry time.Duration) (string, uuid.UUID, error) {
	// Delegate to FileSystem's GeneratePresignedURL method
	presignedURL, imageId, err := handler.FileSystem.GeneratePresignedURL(expiry)
	if err != nil {
		return "", uuid.UUID{}, err
	}
	return presignedURL, imageId, nil
}

func (handler *Handler) MoveFileToFolder(file common.File, src, dest string) error {
	err := handler.FileSystem.MoveFileToFolder(file, src, dest)
	if err != nil {
		return err
	}
	return nil
}

func (handler *Handler) GetImageMetaData(imageId, bucket string) (int64, string, error) {
	imageSize, contentType, err := handler.FileSystem.GetObjectMetadata(imageId, bucket)
	if err != nil {
		return 0, "", err // make this better later
	}
	return imageSize, contentType, nil
}
