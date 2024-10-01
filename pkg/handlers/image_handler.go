package handlers

import (
	"bit-image/pkg/services"
	"bit-image/pkg/storage"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PresignedURLRequest struct {
	NumImages int `json:"num_images"`
}

type ConfirmUploadsRequest struct {
	ImageUploads []services.ConfirmUploadRequest `json:"image_uploads"`
}

type PresignedURLResponse struct {
	ImageUploadURLs []services.PresignedURL `json:"image_upload_urls"`
}

type ImageHandler struct {
	ImageService *services.ImageService
	ImageStore   *storage.UserStore
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

func (h *ImageHandler) ConfirmImageUploads() gin.HandlerFunc {
	return func(c *gin.Context) {
		var request ConfirmUploadsRequest
		if err := c.ShouldBindJSON(&request); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
			return
		}

		errors := h.ImageService.ConfirmImageUploads(request.ImageUploads)

		if len(errors) > 0 {
			var errorMessages []string
			for _, err := range errors {
				errorMessages = append(errorMessages, err.Error())
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Some image uploads failed to confirm",
				"errors":  errorMessages,
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "All image uploads confirmed successfully"})
	}
}
