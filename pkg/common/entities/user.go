package entities

import (
	"bit-image/pkg/common"
)

type User struct {
	Base             common.Base `gorm:"embedded;not null"`
	imageUploadLimit int         `gorm:"not null"`
	imageUploadCount int         `gorm:"not null"`
}
