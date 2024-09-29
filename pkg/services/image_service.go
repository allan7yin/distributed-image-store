package services

import (
	"bit-image/internal/s3"
	"bit-image/pkg/common"
	"bit-image/pkg/common/entities"
	"bit-image/pkg/storage/image"
	"fmt"
	"github.com/google/uuid"
	"os"
	"runtime"
	"sync"
	"time"
)

type ImageService struct {
	S3Handler  *s3.Handler
	ImageStore *image.ImageStore
}

type PresignedURL struct {
	URL     string    `json:"url"`
	Method  string    `json:"method"`
	ImageId uuid.UUID `json:"image_id"`
}

type ConfirmUploadRequest struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	IsPrivate bool   `json:"is_private"`
}

func NewImageService(store *image.ImageStore, s3Handler *s3.Handler) *ImageService {
	return &ImageService{
		ImageStore: store,
		S3Handler:  s3Handler,
	}
}

func (svc *ImageService) GeneratePresignedURLs(NumImages int) ([]PresignedURL, error) {
	if svc == nil {
		return nil, fmt.Errorf("s3 handler not set")
	}
	numCores := runtime.NumCPU()
	numWorkers := min(NumImages, numCores)

	urls := make(chan PresignedURL, NumImages)
	errors := make(chan error, NumImages)
	tasks := make(chan struct{}, NumImages)

	var wg sync.WaitGroup

	// Worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range tasks {
				url, imageId, err := svc.S3Handler.GeneratePresignedURL(15 * time.Minute)
				if err != nil {
					errors <- err
					return
				}
				urls <- PresignedURL{
					URL:     url,
					Method:  "PUT",
					ImageId: imageId,
				}
			}
		}()
	}

	for i := 0; i < NumImages; i++ {
		tasks <- struct{}{}
	}
	close(tasks)

	go func() {
		wg.Wait()
		close(urls)
		close(errors)
	}()

	var presignedURLs []PresignedURL
	for url := range urls {
		presignedURLs = append(presignedURLs, url)
	}

	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to generate some presigned URLs")
	}

	return presignedURLs, nil
}

func (svc *ImageService) ConfirmImageUploads(uploadRequests []ConfirmUploadRequest) {
	var wg sync.WaitGroup
	for _, uploadRequest := range uploadRequests {
		wg.Add(1)
		go func(uploadRequest ConfirmUploadRequest) {
			defer wg.Done()
			svc.ConfirmImage(uploadRequest)
		}(uploadRequest)
	}

	wg.Wait()
}

func (svc *ImageService) ConfirmImage(uploadRequest ConfirmUploadRequest) error {
	// Parse UUID from the request
	ImageId, err := uuid.Parse(uploadRequest.Id)
	if err != nil {
		return fmt.Errorf("failed to parse UUID from request ID %s: %w", uploadRequest.Id, err)
	}

	file := common.File{
		Id:   ImageId.String(),
		Hash: uploadRequest.Hash,
	}

	// Obtain the metadata
	imageSize, contentType, err := svc.S3Handler.GetImageMetaData(common.TEMPORARY_STORAGE_FOLDER+"/"+ImageId.String(), os.Getenv("DEFAULT_BUCKET_NAME"))
	if err != nil {
		return fmt.Errorf("failed to get metadata for image with ID %s: %w", ImageId.String(), err)
	}
	fmt.Printf("Image metadata retrieved - Size: %d bytes, Content-Type: %s\n", imageSize, contentType)

	err = svc.S3Handler.MoveFileToFolder(file, common.TEMPORARY_STORAGE_FOLDER, common.PERMANENT_STORAGE_FOLDER)
	if err != nil {
		return fmt.Errorf("failed to move file with ID %s to folder %s: %w", file.Id, common.PERMANENT_STORAGE_FOLDER, err)
	}

	fmt.Printf("Image with ID %s successfully moved to permanent storage folder.\n", file.Id)

	//Placeholder for creating a new Image entity
	NewImage := entities.Image{
		Base: common.Base{
			Id: ImageId,
		},
		Name:      uploadRequest.Name,
		IsPrivate: uploadRequest.IsPrivate,
		ImageMetaData: common.ImageMetaData{
			Hash: uploadRequest.Hash,
		},
	}

	fmt.Println(NewImage)
	return nil
}
