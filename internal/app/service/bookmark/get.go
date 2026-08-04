package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
)

// GetBookmarksResult represents the result containing a list of bookmarks and the total count
type GetBookmarksResult struct {
	Bookmarks []*model.Bookmark `json:"bookmarks"`
	Total     int64             `json:"total"`
}

// GetBookmarks retrieves paginated bookmarks and total count for a specific user
func (s *bookmarkService) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error) {
	offset := (page - 1) * limit
	bookmarks, count, err := s.repo.GetBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return &GetBookmarksResult{
		Bookmarks: bookmarks,
		Total:     count,
	}, nil
}
