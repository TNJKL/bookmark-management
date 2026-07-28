package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
)

func (s *service) GetSelfInfo(ctx context.Context, uid string) (*model.User, error) {
	return s.repo.GetUserByID(ctx, uid)
}

func (s *service) UpdateSelfInfo(ctx context.Context, uid, displayName, email string) error {
	user, err := s.repo.GetUserByID(ctx, uid)
	if err != nil {
		return err
	}
	user.DisplayName = displayName
	user.Email = email

	return s.repo.UpdateUser(ctx, user)
}
