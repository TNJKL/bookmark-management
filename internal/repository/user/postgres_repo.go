package user

import "gorm.io/gorm"

// sqlRepository implements the Repository interface using GORM
type sqlRepository struct {
	db *gorm.DB
}

// NewSQLRepository creates a new instance of Repository with the GORM database connection
func NewSQLRepository(db *gorm.DB) Repository {
	return &sqlRepository{db: db}
}
