package storage

import (
	"bit-image/pkg/common/entities"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

// AddUser adds a user to the database
func (s *UserStore) AddUser(user entities.User) error {
	// Start a new transaction
	tx := s.db.Begin()

	// Directly use the User entity (no mapper needed)
	if err := tx.Create(&user).Error; err != nil {
		// Rollback transaction if there's an error
		tx.Rollback()
		return err
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DeleteUserByID deletes a user by ID
func (s *UserStore) DeleteUserByID(userID string) error {
	tx := s.db.Begin()

	if err := tx.Delete(&entities.User{}, "id = ?", userID).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

// DoesUserExist checks if a user exists in the database
func (s *UserStore) DoesUserExist(userID string) (bool, error) {
	var user entities.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
