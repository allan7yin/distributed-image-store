package handlers

import "bit-image/pkg/common/entities"

type ImageConfirmRequest struct {
	ID        string   `json:"id" validate:"required,uuid4"`
	Name      string   `json:"name" validate:"required"`
	Hash      string   `json:"hash" validate:"required"`
	IsPrivate bool     `json:"is_private"`
	Tags      []string `json:"tags" validate:"required,dive"`
}

type ImageUploadService interface {
	GenerateImageUploadUrls(count int, userID string) ([]string, error)
	ConfirmImagesUploaded(cmds[]) ([]entities.Image, error)
	GetImage(userID, imageID string) (Image, error)
	GetAllPublicImages() ([]Image, error)
}


func (c)
