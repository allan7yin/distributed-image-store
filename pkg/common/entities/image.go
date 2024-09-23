package entities

import (
	"bit-image/pkg/common"
	"fmt"
)

// Image TO DO: refactor the Tags and Content Labels into structs -> easier when querying by those values in the db
type Image struct {
	common.Base
	name          string               `gorm:"not null"`
	tags          []string             `gorm:"type:jsonb;not null"`
	contentLabels []string             `gorm:"type:jsonb;not null"`
	isPrivate     bool                 `gorm:"not null"`
	imageMetaData common.ImageMetaData `gorm:"embedded;not null"`
}

// Define the method on the Image struct
func (image Image) mapToFileID(fileIDPrefix string) string {
	return fmt.Sprintf(fileIDPrefix, image.Base.Id, image.Base.Id)
}

// Example usage in the imageToFileMapper function
func (image Image) imageToFileMapper() common.File {
	file := common.File{
		Id:   image.mapToFileID(fileIDPrefix), // Now calling the method on the struct
		Hash: image.imageMetaData.Hash,
	}

	return file
}
