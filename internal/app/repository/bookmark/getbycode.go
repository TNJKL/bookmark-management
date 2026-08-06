package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

// GetByCode retrieves a single Bookmark record from the database matching the provided short code.
func (r *bookmarkRepo) GetByCode(ctx context.Context, code string) (*model.Bookmark, error) {
	bookmark := &model.Bookmark{}
	err := r.db.WithContext(ctx).Where("code = ?", code).First(bookmark).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return bookmark, nil
}
