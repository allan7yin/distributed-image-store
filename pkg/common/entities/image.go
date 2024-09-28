package entities

import (
	"bit-image/pkg/common"
	"fmt"
)

// Image TO DO: refactor the Tags and Content Labels into structs -> easier when querying by those values in the db
type Image struct {
	common.Base
	Name string `gorm:"not null"`
	//tags          []string             `gorm:"type:jsonb;not null"`
	//contentLabels []string             `gorm:"type:jsonb;not null"`
	IsPrivate     bool                 `gorm:"not null"`
	ImageMetaData common.ImageMetaData `gorm:"embedded;not null"`
}

// Define the method on the Image struct
func (image Image) mapToFileID(fileIDPrefix string) string {
	return fmt.Sprintf(fileIDPrefix, image.Base.Id, image.Base.Id)
}
