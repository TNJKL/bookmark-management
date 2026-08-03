package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

// UpdateBookmark updates description and url of a bookmark for a user
func (r *bookmarkRepo) UpdateBookmark(ctx context.Context, id, userID, description, url string) error {
	res := r.db.WithContext(ctx).Model(&model.Bookmark{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]interface{}{
			"description": description,
			"url":         url,
		})

	if res.Error != nil {
		return dbutils.CatchDBError(res.Error)
	}
	if res.RowsAffected == 0 {
		return dbutils.ErrRecordNotFound
	}
	return nil
}
