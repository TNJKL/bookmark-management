package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"gorm.io/gorm"
)

// CreateBookmark inserts a new bookmark record into the database
func (r *bookmarkRepo) CreateBookmark(ctx context.Context, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := r.db.WithContext(ctx).Create(bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}

// CreateBookmarkTx inserts a bookmark within an existing transaction
func (r *bookmarkRepo) CreateBookmarkTx(ctx context.Context, tx *gorm.DB, bookmark *model.Bookmark) (*model.Bookmark, error) {
	err := tx.WithContext(ctx).Create(bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}

// UpdateCodeTx updates the code column of a bookmark within an existing transaction
func (r *bookmarkRepo) UpdateCodeTx(ctx context.Context, tx *gorm.DB, id string, code string) error {
	return tx.WithContext(ctx).Model(&model.Bookmark{}).Where("id = ?", id).Update("code", code).Error
}
