package entities

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base embeds common fields ที่ทุก entity ใช้ร่วมกัน
type Base struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey"                    json:"id"`
	CreatedAt time.Time      `gorm:"not null;default:now()"                  json:"created_at"`
	UpdatedAt time.Time      `gorm:"not null;default:now()"                  json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"                                   json:"deleted_at,omitempty"`
}

// BeforeCreate — auto-generate UUID ก่อน insert
func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}
