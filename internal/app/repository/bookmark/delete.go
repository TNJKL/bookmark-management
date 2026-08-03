package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

//Ở đây e k biết rõ yêu cầu là xóa hẳn hay là soft delete bookmark nên e làm tạm xóa hẳn luôn nha anh Duy Anh

// DeleteBookmark deletes a bookmark belonging to a specific user by id
func (r *bookmarkRepo) DeleteBookmark(ctx context.Context, id, userID string) error {
	res := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Bookmark{})
	if res.Error != nil {
		return dbutils.CatchDBError(res.Error)
	}
	if res.RowsAffected == 0 {
		return dbutils.ErrRecordNotFound
	}
	return nil
}
