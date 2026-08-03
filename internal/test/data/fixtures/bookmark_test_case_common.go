package fixtures

import (
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

// BookmarkCommonTestDB provides test database fixtures with seeded user and bookmark records
type BookmarkCommonTestDB struct {
	UserCommonTestDB
}

// Migrate auto-migrates User and Bookmark models into the test database
func (b *BookmarkCommonTestDB) Migrate() error {
	return b.db.AutoMigrate(&model.User{}, &model.Bookmark{})
}

// GenerateData seeds the test database with initial user and bookmark data
func (b *BookmarkCommonTestDB) GenerateData() error {
	err := b.UserCommonTestDB.GenerateData()
	if err != nil {
		return err
	}
	bookmarks := []*model.Bookmark{
		{
			Base:        GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd993"),
			Description: "bookmark1",
			URL:         "https://google.com",
			Code:        "123456",
			UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
		},
		{
			Base:        GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd992"),
			Description: "bookmark2",
			URL:         "https://google.com",
			Code:        "123457",
			UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
		},
		{
			Base:        GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd994"),
			Description: "bookmark3",
			URL:         "https://github.com",
			Code:        "123458",
			UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
		},
		{
			Base:        GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd995"),
			Description: "bookmark4",
			URL:         "https://youtube.com",
			Code:        "123459",
			UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
		},
		{
			Base:        GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd998"),
			Description: "bookmark5",
			URL:         "https://golang.org",
			Code:        "123460",
			UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
		},
	}
	return b.db.Session(&gorm.Session{SkipHooks: true}).CreateInBatches(bookmarks, 10).Error
}
