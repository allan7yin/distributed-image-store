package services

import (
	"bit-image/internal/s3"
	"fmt"
	"github.com/google/uuid"
	"sync"
	"time"
)

/*
type ImageUploadService interface {
	GenerateImageUploadUrls(count int, userID string) ([]string, error)
	//ConfirmImagesUploaded(cmds[]) ([]entities.Image, error)
	//GetImage(userID, imageID string) (Image, error)
	//GetAllPublicImages() ([]Image, error)
}
*/

type ImageService struct {
	S3Handler *s3.Handler
}

type PresignedURL struct {
	URL     string    `json:"url"`
	Method  string    `json:"method"`
	ImageId uuid.UUID `json:"image_id"`
}

func NewImageService(s3Handler *s3.Handler) *ImageService {
	return &ImageService{S3Handler: s3Handler}
}

func (svc *ImageService) GeneratePresignedURLs(NumImages int) ([]PresignedURL, error) {
	if svc == nil {
		return nil, fmt.Errorf("s3 handler not set")
	}

	numWorkers := min(NumImages, 20)

	urls := make(chan PresignedURL, NumImages)
	errors := make(chan error, NumImages)
	tasks := make(chan struct{}, NumImages)

	var wg sync.WaitGroup

	// Worker pool
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range tasks { // Each worker picks up tasks from the 'tasks' channel
				url, imageId, err := svc.S3Handler.GeneratePresignedURL(15 * time.Minute)
				if err != nil {
					errors <- err
					return
				}
				urls <- PresignedURL{
					URL:     url,
					Method:  "PUT", // Assuming PUT for uploads
					ImageId: imageId,
				}
			}
		}()
	}

	// Send tasks to workers
	for i := 0; i < NumImages; i++ {
		tasks <- struct{}{}
	}
	close(tasks) // Close the tasks channel after sending all tasks

	// Wait for all workers to finish and then close the urls channel
	go func() {
		wg.Wait()
		close(urls)
		close(errors)
	}()

	// Collect presigned URLs
	var presignedURLs []PresignedURL
	for url := range urls {
		presignedURLs = append(presignedURLs, url)
	}

	// Handle errors
	if len(errors) > 0 {
		return nil, fmt.Errorf("failed to generate some presigned URLs")
	}

	return presignedURLs, nil
}
