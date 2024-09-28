package repositories

import (
	"bit-image/internal/postrges"
	"bit-image/pkg/common/entities"
	"github.com/google/uuid"
)

type ImageRepository struct {
	handler *postrges.ConnectionHandler
}

func NewImageRepository(handler *postrges.ConnectionHandler) *ImageRepository {
	return &ImageRepository{}
}

func (repo *ImageRepository) InsertImage(req ConfirmUploadRequest) (entities.Image, error) {
	tx, commit, rollback, err := repo.handler.OpenTransaction()
	if err != nil {
		return entities.Image{}, err
	}

	defer func() {
		if r := recover(); r != nil {
			rollback()
		}
	}()

	// Convert ConfirmUploadRequest to entities.Image
	newImage := entities.Image{
		Base: entities.Base{
			Id: uuid.New(), // Or parse the provided `req.Id` if necessary
		},
		Name:      req.Name,
		IsPrivate: req.IsPrivate,
		ImageMetaData: entities.ImageMetaData{
			Hash: req.Hash,
			// Fill in other metadata as needed
		},
	}

	// Insert image using the transaction
	if err := tx.Create(&newImage).Error; err != nil {
		rollback() // Rollback on error
		return entities.Image{}, err
	}

	// Commit the transaction if successful
	if err := commit(); err != nil {
		return entities.Image{}, err
	}

	return newImage, nil
}

/*
type ConfirmUploadRequest struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Hash      string `json:"hash"`
	IsPrivate bool   `json:"is_private"`
	//Tags      []Tag  `json:"tags"`
}
*/
