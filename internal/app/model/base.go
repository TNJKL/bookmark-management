package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Base contains common fields embedded into application models (ID, CreatedAt, UpdatedAt, DeletedAt)
type Base struct {
	ID        string    `json:"id" gorm:"type:uuid;primaryKey;column:id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	DeletedAt time.Time `json:"-"`
}

// BeforeCreate is a GORM hook that automatically generates a new UUID for the user ID if empty
func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == "" {
		b.ID = uuid.New().String()
	}
	return
}
