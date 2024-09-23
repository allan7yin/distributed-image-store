package common

// User represents a user entity in the system
type User struct {
	Base
	ImageUploadLimit int `gorm:"not null"`
	ImageUploadCount int `gorm:"not null"`
}
