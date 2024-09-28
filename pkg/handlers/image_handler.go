package handlers

import (
	"bit-image/pkg/services"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PresignedURLRequest struct {
	NumImages int `json:"num_images"`
}

type ConfirmUploadsRequest struct {
	ImageUploads []ConfirmUploadRequest `json:"image_uploads"`
}

type ConfirmUploadRequest struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	IsPrivate bool   `json:"is_private"`
}

type PresignedURLResponse struct {
	ImageUploadURLs []services.PresignedURL `json:"image_upload_urls"`
}

type ImageHandler struct {
	ImageService *services.ImageService
}

func NewImageHandler(imageService *services.ImageService) *ImageHandler {
	return &ImageHandler{ImageService: imageService}
}

func (h *ImageHandler) GeneratePresignedURL() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request PresignedURLRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		if request.NumImages <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		urls, _ := h.ImageService.GeneratePresignedURLs(request.NumImages)

		response := PresignedURLResponse{
			ImageUploadURLs: urls,
		}

		c.Writer.Header().Set("Content-Type", "application/json")
		c.Writer.WriteHeader(http.StatusOK)

		encoder := json.NewEncoder(c.Writer)
		encoder.SetEscapeHTML(false)
		err := encoder.Encode(response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write response"})
			return
		}
	}
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
*/
func (h *ImageHandler) ConfirmUploadImage() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request ConfirmUploadsRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		// now, we will take this and create an Image
		h.ImageService.
	}
}
