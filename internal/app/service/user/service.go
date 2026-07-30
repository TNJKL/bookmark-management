package user

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/app/repository/user"
	"github.com/TNJKL/bookmark-management/pkg/jwtutils"
	"github.com/TNJKL/bookmark-management/pkg/utils"
)

// Service defines the business logic operations for user management
//
//go:generate mockery --name Service --filename serivce.go
type Service interface {
	CreateUser(ctx context.Context, username, password, displayName, email string) (*model.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	GetSelfInfo(ctx context.Context, uid string) (*model.User, error)
	UpdateSelfInfo(ctx context.Context, uid, displayName, email string) error
}

// service implements the Service interface using the repository and hasher dependencies
type service struct {
	repo         user.Repository
	hasher       utils.Hasher
	jwtGenerator jwtutils.JWTGenerator
}

// NewService creates a new instance of Service.
func NewService(repo user.Repository, hasher utils.Hasher, jwtGen jwtutils.JWTGenerator) Service {
	return &service{repo: repo, hasher: hasher, jwtGenerator: jwtGen}
}
