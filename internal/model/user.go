package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the user account data model in the database
type User struct {
	ID          string    `json:"id" gorm:"type:uuid;primaryKey"`
	DisplayName string    `json:"display_name"`
	Username    string    `json:"username" gorm:"unique"`
	Password    string    `json:"-"`
	Email       string    `json:"email" gorm:"unique"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BeforeCreate is a GORM hook that automatically generates a new UUID for the user ID if empty
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == "" {
		u.ID = uuid.New().String()
	}
	return
}
