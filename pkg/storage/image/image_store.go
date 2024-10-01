package image

import (
	"bit-image/internal/postrges"
	"bit-image/pkg/common/entities"
	"fmt"
	"gorm.io/gorm"
)

type ImageStore struct {
	DBHandler *postrges.ConnectionHandler
}

func NewImageStore(dbHandler *postrges.ConnectionHandler) *ImageStore {
	return &ImageStore{
		DBHandler: dbHandler,
	}
}

// TODO: might not need, this is just a transaction for a single record write
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

func (store *ImageStore) AddImageWithTransaction(tx *gorm.DB, image entities.Image) error {
	if err := tx.Create(&image).Error; err != nil {
		return fmt.Errorf("failed to insert image: %w", err)
	}
	return nil
}
