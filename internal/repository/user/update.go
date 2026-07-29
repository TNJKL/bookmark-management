package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

// UpdateUser updates the display_name and email of an existing user in the database
func (r *sqlRepository) UpdateUser(ctx context.Context, currentuser *model.User) error {
	err := r.db.WithContext(ctx).Model(currentuser).Updates(model.User{
		DisplayName: currentuser.DisplayName,
		Email:       currentuser.Email,
	}).Error
	if err != nil {
		return dbutils.CatchDBError(err)
	}
	return nil
}
