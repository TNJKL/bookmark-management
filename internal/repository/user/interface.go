package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
)

// Repository defines the database operations for user data management
//
//go:generate mockery --name Repository --filename sqlRepository.go
type Repository interface {
	CreateUser(ctx context.Context, newUser *model.User) (*model.User, error)
}
