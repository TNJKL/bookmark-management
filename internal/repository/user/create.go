package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
)

// CreateUser inserts a new user record into the database and returns the created user
func (r *sqlRepository) CreateUser(ctx context.Context, newUser *model.User) (*model.User, error) {
	err := r.db.WithContext(ctx).Create(newUser).Error
	if err != nil {
		return nil, dbutils.CatchDBError(err)
	}
	return newUser, err
}
