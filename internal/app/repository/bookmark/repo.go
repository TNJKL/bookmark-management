package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"gorm.io/gorm"
)

//go:generate mockery --name Repository --filename bookmarkRepo.go
type Repository interface {
	CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error)
	UpdateBookmark(ctx context.Context, id, userID, description, url string) error
	DeleteBookmark(ctx context.Context, id, userID string) error
}

// bookmarkRepo is the concrete implementation of the Repository interface using GORM
type bookmarkRepo struct {
	db *gorm.DB
}

// NewRepository creates a new instance of bookmarkRepo with the provided GORM database connection
func NewRepository(db *gorm.DB) Repository {
	return &bookmarkRepo{db: db}
}
