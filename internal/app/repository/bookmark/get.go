package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

// GetBookmarks retrieves a paginated list of bookmarks and total count for a user
func (r *bookmarkRepo) GetBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, int64, error) {
	bookmarks, err := r.getBookmarks(ctx, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := r.countBookmarks(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	return bookmarks, count, nil
}

// getBookmarks queries a slice of bookmarks from the database with pagination parameters
func (r *bookmarkRepo) getBookmarks(ctx context.Context, userID string, limit, offset int) ([]*model.Bookmark, error) {
	bookmarks := make([]*model.Bookmark, limit)
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at ASC").Offset(offset).Limit(limit).Find(&bookmarks).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmarks, nil
}

// countBookmarks counts the total number of bookmarks belonging to a user in the database
func (r *bookmarkRepo) countBookmarks(ctx context.Context, userID string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&model.Bookmark{}).Where("user_id = ?", userID).Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}
