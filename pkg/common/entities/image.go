package entities

import "bit-image/pkg/common"

// Image TO DO: refactor the Tags and Content Labels into structs -> easier when querying by those values in the db
type Image struct {
	Base          common.Base          `gorm:"embedded;not null"`
	Name          string               `gorm:"not null"`
	IsPrivate     bool                 `gorm:"not null"`
	ImageMetaData common.ImageMetaData `gorm:"embedded;not null"`
}
