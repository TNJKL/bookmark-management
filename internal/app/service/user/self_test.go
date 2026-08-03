package user

import (
	"context"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/user/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	jwtMocks "github.com/TNJKL/bookmark-management/pkg/jwtutils/mocks"
	"github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetSelfInfo(t *testing.T) {
	t.Parallel()
	inputID := "uuid-123"
	//errTest := errors.New("test error")

	testUser := &model.User{
		Username:    "testuser",
		Email:       "test@example.com",
		DisplayName: "Test User",
	}
	testUser.ID = inputID

	testCases := []struct {
		name          string
		setupMockRepo func(ctx context.Context) *repoMocks.Repository
		expectedUser  *model.User
		expectedErr   error
	}{
		{
			name: "happy path",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByID", ctx, inputID).Return(testUser, nil)
				return mockRepo
			},
			expectedUser: testUser,
			expectedErr:  nil,
		},

		{
			name: "user not found",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByID", ctx, inputID).Return(nil, dbutils.ErrRecordNotFound)
				return mockRepo
			},
			expectedUser: nil,
			expectedErr:  dbutils.ErrRecordNotFound,
		},

		{
			name: "database error",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByID", ctx, inputID).Return(nil, errTest)
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

			mockRepo := tc.setupMockRepo(ctx)
			mockHasher := mocks.NewHasher(t)
			mockJWTGen := jwtMocks.NewJWTGenerator(t)

			svc := NewService(mockRepo, mockHasher, mockJWTGen)

			user, err := svc.GetSelfInfo(ctx, inputID)
			assert.Equal(t, tc.expectedUser, user)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}

func TestService_UpdateSelfInfo(t *testing.T) {
	t.Parallel()
	inputUID := "uuid-123"
	newDisplayName := "New Name"
	newEmail := "new@example.com"
	//errTest := errors.New("test error")

	testUser := &model.User{
		Username:    "testuser",
		Email:       "old@example.com",
		DisplayName: "Old Name",
	}
	testUser.ID = inputUID

	testCases := []struct {
		name          string
		setupMockRepo func(ctx context.Context) *repoMocks.Repository
		expectedErr   error
	}{
		{
			name: "happy path",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByID", ctx, inputUID).Return(testUser, nil)
				mockRepo.On("UpdateUser", ctx, mock.MatchedBy(func(u *model.User) bool {
					return u.DisplayName == newDisplayName && u.Email == newEmail
				})).Return(nil)
				return mockRepo
			},
			expectedErr: nil,
		},

		{
			name: "user not found",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByID", ctx, inputUID).Return(nil, dbutils.ErrRecordNotFound)
				return mockRepo
			},
			expectedErr: dbutils.ErrRecordNotFound,
		},

		{
			name: "repo update fail",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByID", ctx, inputUID).Return(testUser, nil)
				mockRepo.On("UpdateUser", ctx, testUser).Return(dbutils.ErrDuplicationEmail)
				return mockRepo
			},
			expectedErr: dbutils.ErrDuplicationEmail,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)
			mockHasher := mocks.NewHasher(t)
			mockJWTGen := jwtMocks.NewJWTGenerator(t)

			svc := NewService(mockRepo, mockHasher, mockJWTGen)
			err := svc.UpdateSelfInfo(ctx, inputUID, newDisplayName, newEmail)
			assert.ErrorIs(t, tc.expectedErr, err)
		})
	}
}
