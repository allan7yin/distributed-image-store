package common

// TO DO: refactor the Tags and Content Labels into structs -> easier when querying by those values in the db
type Image struct {
	Base
	Name          string        `gorm:"not null"`
	Tags          []string      `gorm:"not null"`
	ContentLabels []string      `gorm:"not null"`
	IsPrivate     bool          `gorm:"not null"`
	ImageMetaData ImageMetaData `gorm:"embedded;not null"`
}
