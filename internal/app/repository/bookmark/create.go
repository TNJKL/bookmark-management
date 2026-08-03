package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

// CreateBookmark inserts a new bookmark record into the database
func (r *bookmarkRepo) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := r.db.WithContext(ctx).Create(bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}
