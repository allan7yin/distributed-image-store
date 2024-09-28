package services

import (
	"bit-image/internal/s3"
	"bit-image/pkg/common"
	"bit-image/pkg/common/entities"
	"fmt"
	"github.com/google/uuid"
	"runtime"
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

type ConfirmUploadRequest struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	IsPrivate bool   `json:"is_private"`
}

func NewImageService(s3Handler *s3.Handler) *ImageService {
	return &ImageService{S3Handler: s3Handler}
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

// next step -> do ConfirmImagesUploaded
/*
Rundown:
Request body will include:
- image id
- name
- hash
- is_private (determines if the s3 object is public or not)
- tags -> tags for the image the user has assigned


steps:
- move image in s3 to a "TEMP_HOME" folder
- then, retrieve the metadata for that image -> size, type, etc.
- then, insert this into Image entity and insert into database
*/

func (svc *ImageService) MoveToFolder(ImageId int) ([]PresignedURL, error) {}

func (svc *ImageService) CreateImage(req ConfirmUploadRequest) {
	// construct the Image entity from the request
	ImageId, err := uuid.Parse(req.Id)
	if err != nil {
		fmt.Println("Invalid UUID string:", err)
		return
	}

	// need to first retreiev the other image data

	NewImage := entities.Image{
		Base: common.Base{
			Id: ImageId,
		},
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
		ImageMetaData: common.ImageMetaData{
			Hash: req.Hash,
		}
	}
}
