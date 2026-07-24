package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
)

// CreateUser hashes the password and calls the repository to save a new user
func (s *service) CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error) {
	// hash password
	pwdHash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}
	// create user model
	newUser := &model.User{
		Username:    username,
		Password:    pwdHash,
		Email:       email,
		DisplayName: displayName,
	}
	// call repo to create user
	res, err := s.repo.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	//return user
	return res, nil
}
