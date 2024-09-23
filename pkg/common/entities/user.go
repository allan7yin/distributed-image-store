package entities

import "bit-image/pkg/common"

type User struct {
	common.Base
	imageUploadLimit int `gorm:"not null"`
	imageUploadCount int `gorm:"not null"`
}
