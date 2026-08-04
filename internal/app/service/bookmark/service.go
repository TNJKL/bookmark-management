package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/app/repository/bookmark"
	"github.com/TNJKL/bookmark-management/pkg/utils"
)

// Service defines business logic operations for managing bookmarks
//
//go:generate mockery --name Service --filename bookmarkService.go
type Service interface {
	CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error)
	GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error)
	UpdateBookmark(ctx context.Context, id, userID, description, url string) error
	DeleteBookmark(ctx context.Context, id, userID string) error
}

// bookmarkService implements the Service interface
type bookmarkService struct {
	repo   bookmark.Repository
	keyGen utils.KeyGenerator
}

// NewService creates a new instance of bookmarkService with provided repository and key generator
func NewService(repo bookmark.Repository, keyGen utils.KeyGenerator) Service {
	return &bookmarkService{
		repo:   repo,
		keyGen: keyGen,
	}
}
