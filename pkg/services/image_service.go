package services

import (
	"bit-image/internal/s3"
	"bit-image/pkg/common"
	"bit-image/pkg/common/entities"
	"bit-image/pkg/storage/image"
	"fmt"
	"github.com/google/uuid"
	"log"
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

func (svc *ImageService) ConfirmImageUploads(uploadRequests []ConfirmUploadRequest) []error {
	// idiomatic go -> handle errors with a wait channel
	var wg sync.WaitGroup
	errChan := make(chan error, len(uploadRequests))

	for _, uploadRequest := range uploadRequests {
		wg.Add(1)
		go func(uploadRequest ConfirmUploadRequest) {
			defer wg.Done()

			// Call ConfirmImage and send error to channel if any
			if err := svc.ConfirmImage(uploadRequest); err != nil {
				errChan <- fmt.Errorf("failed to confirm upload for request ID %s: %w", uploadRequest.Id, err)
			} else {
				errChan <- nil
			}
		}(uploadRequest)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		return errs
	}
	return nil
}

func (svc *ImageService) ConfirmImage(uploadRequest ConfirmUploadRequest) error {
	imageID, err := uuid.Parse(uploadRequest.Id)
	if err != nil {
		return fmt.Errorf("failed to parse UUID from request ID %s: %w", uploadRequest.Id, err)
	}

	imageSize, contentType, err := svc.S3Handler.GetImageMetaData(common.TEMPORARY_STORAGE_FOLDER+"/"+imageID.String(), os.Getenv("DEFAULT_BUCKET_NAME"))
	if err != nil {
		return fmt.Errorf("failed to get metadata for image with ID %s: %w", imageID.String(), err)
	}
	fmt.Printf("Image metadata retrieved - Size: %d bytes, Content-Type: %s\n", imageSize, contentType)

	tx, commit, rollback, err := svc.ImageStore.DBHandler.OpenTransaction()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	file := common.File{
		Id:   imageID.String(),
		Hash: uploadRequest.Hash,
	}
	if err = svc.S3Handler.MoveFileToFolder(file, common.TEMPORARY_STORAGE_FOLDER, common.PERMANENT_STORAGE_FOLDER); err != nil {
		return fmt.Errorf("failed to move file with ID %s to folder %s: %w", file.Id, common.PERMANENT_STORAGE_FOLDER, err)
	}
	fmt.Printf("Image with ID %s successfully moved to permanent storage folder.\n", file.Id)

	newImage := entities.Image{
		Base: common.Base{
			Id: imageID,
		},
		Name:      uploadRequest.Name,
		IsPrivate: uploadRequest.IsPrivate,
		Path:      common.PERMANENT_STORAGE_FOLDER + "/" + imageID.String(),
		ImageMetaData: common.ImageMetaData{
			Hash:     uploadRequest.Hash,
			FileSize: float64(imageSize),
			Format:   contentType,
		},
	}

	if err = svc.ImageStore.AddImageWithTransaction(tx, newImage); err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			log.Printf("failed to rollback transaction: %v", rollbackErr)
		}

		if moveBackErr := svc.S3Handler.MoveFileToFolder(file, common.PERMANENT_STORAGE_FOLDER, common.TEMPORARY_STORAGE_FOLDER); moveBackErr != nil {
			fmt.Printf("Warning: failed to move file back to temporary folder after DB insert failure: %v\n", moveBackErr)
		}
		return fmt.Errorf("failed to save image metadata to database: %w", err)
	}

	if err = commit(); err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			log.Printf("failed to rollback transaction: %v", rollbackErr)
		}

		if moveBackErr := svc.S3Handler.MoveFileToFolder(file, common.PERMANENT_STORAGE_FOLDER, common.TEMPORARY_STORAGE_FOLDER); moveBackErr != nil {
			fmt.Printf("Warning: failed to move file back to temporary folder after commit failure: %v\n", moveBackErr)
		}
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
