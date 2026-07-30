package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
)

// Repository defines the database operations for user data management
//
//go:generate mockery --name Repository --filename sqlRepository.go
type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
	GetUserByUsername(ctx context.Context, username string) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
}
