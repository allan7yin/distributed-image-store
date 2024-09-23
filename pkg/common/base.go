package common

import (
	"github.com/google/uuid"
	"time"
)

type Base struct {
	ID              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	DateTimeCreated time.Time `gorm:"autoCreateTime"`
	DateTimeUpdated time.Time `gorm:"autoUpdateTime"`
}
