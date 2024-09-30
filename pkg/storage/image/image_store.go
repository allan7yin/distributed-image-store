package image

import (
	"bit-image/internal/postrges"
	"bit-image/pkg/common/entities"
	"fmt"
)

type ImageStore struct {
	DBHandler *postrges.ConnectionHandler
}

func NewImageStore(dbHandler *postrges.ConnectionHandler) *ImageStore {
	return &ImageStore{
		DBHandler: dbHandler,
	}
}

func (store *ImageStore) AddImage(image entities.Image) error {
	tx, commit, rollback, err := store.DBHandler.OpenTransaction()
	if err != nil {
		return err
	}

	// Insert the entity within the transaction
	if err := tx.Create(&image).Error; err != nil {
		if rollbackErr := rollback(); rollbackErr != nil {
			return fmt.Errorf("insert error: %v, rollback error: %v", err, rollbackErr)
		}
		return err
	}

	// Commit the transaction if everything is successful
	if err := commit(); err != nil {
		return fmt.Errorf("commit error: %v", err)
	}
	return nil
}
