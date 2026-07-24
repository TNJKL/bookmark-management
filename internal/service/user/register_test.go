package user

import (
	"context"
	"errors"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/model"
	repoMocks "github.com/TNJKL/bookmark-management/internal/repository/user/mocks"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

func TestService_CreateUser(t *testing.T) {
	t.Parallel()
	inputUsername := "testuser"
	inputPassword := "password123"
	inputDisplayName := "Test User"
	inputEmail := "testuser@example.com"
	hashedPassword := "hashed_password123"

	expectedUser := &model.User{
		Username:    inputUsername,
		Password:    hashedPassword,
		Email:       inputEmail,
		DisplayName: inputDisplayName,
	}

	testCases := []struct {
		name            string
		setupMockHasher func() *mocks.Hasher
		setupMockRepo   func(ctx context.Context) *repoMocks.Repository
		expectedUser    *model.User
		expectedErr     error
	}{
		{
			name: "happy path",
			setupMockHasher: func() *mocks.Hasher {
				mockHasher := mocks.NewHasher(t)
				mockHasher.On("Hash", inputPassword).Return(hashedPassword, nil)
				return mockHasher
			},
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateUser", ctx, expectedUser).Return(expectedUser, nil)
				return mockRepo
			},
			expectedUser: expectedUser,
			expectedErr:  nil,
		},
		{
			name: "hash password fail",
			setupMockHasher: func() *mocks.Hasher {
				mockHasher := mocks.NewHasher(t)
				mockHasher.On("Hash", inputPassword).Return("", utils.ErrHashFailed)
				return mockHasher
			},
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				return repoMocks.NewRepository(t)
			},
			expectedUser: nil,
			expectedErr:  utils.ErrHashFailed,
		},
		{
			name: "repo create user fail",
			setupMockHasher: func() *mocks.Hasher {
				mockHasher := mocks.NewHasher(t)
				mockHasher.On("Hash", inputPassword).Return(hashedPassword, nil)
				return mockHasher
			},
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateUser", ctx, expectedUser).Return(nil, errTest)
				return mockRepo
			},
			expectedUser: nil,
			expectedErr:  errTest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockHasher := tc.setupMockHasher()
			mockRepo := tc.setupMockRepo(ctx)

			svc := NewService(mockRepo, mockHasher)

			user, err := svc.CreateUser(ctx, inputUsername, inputPassword, inputDisplayName, inputEmail)
			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
