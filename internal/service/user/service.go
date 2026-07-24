package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/repository/user"
	"github.com/TNJKL/bookmark-management/pkg/utils"
)

// Service defines the business logic operations for user management
//
//go:generate mockery --name Service --filename serivce.go
type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
}

// service implements the Service interface using the repository and hasher dependencies
type service struct {
	repo   user.Repository
	hasher utils.Hasher
}

// NewService creates a new instance of Service.
func NewService(repo user.Repository, hasher utils.Hasher) Service {
	return &service{repo: repo, hasher: hasher}
}
