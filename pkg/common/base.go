package common

import (
	"github.com/google/uuid"
	"time"
)

type Base struct {
	Id              uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	dateTimeCreated time.Time `gorm:"autoCreateTime"`
	dateTimeUpdated time.Time `gorm:"autoUpdateTime"`
}
