package fixtures

import (
	"github.com/TNJKL/bookmark-management/internal/model"
	"gorm.io/gorm"
)

// UserCommonTestDB implements the Fixture interface to set up the user database for integration tests
type UserCommonTestDB struct {
	base
}

// Migrate creates the user table schema in the test database
func (u *UserCommonTestDB) Migrate() error {
	return u.db.AutoMigrate(&model.User{})
}

// GenerateData inserts a set of common sample users into the test database
func (u *UserCommonTestDB) GenerateData() error {
	db := u.db.Session(&gorm.Session{SkipHooks: true})

	users := []*model.User{
		{
			ID:          "deb745af-1a62-4efa-99a0-f06b274bd990",
			DisplayName: "John Doe",
			Username:    "johndoe",
			Password:    "password123",
			Email:       "johndoe@example.com",
		},
		{
			ID:          "deb745af-1a62-4efa-99a0-f06b274bd991",
			DisplayName: "Jane Doe",
			Username:    "janedoe",
			Password:    "password123",
			Email:       "janedoe@example.com",
		},
	}
	return db.CreateInBatches(users, 10).Error
}
