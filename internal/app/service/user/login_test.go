package user

import (
	"context"
	"errors"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/user/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	jwtMocks "github.com/TNJKL/bookmark-management/pkg/jwtutils/mocks"
	"github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Login(t *testing.T) {
	t.Parallel()
	inputUsername := "testuser"
	inputPassword := "password123"
	hashedPassword := "hashed_password123"
	jwtToken := "mock-jwt-token-xyz"

	errTest := errors.New("test error")

	testUser := &model.User{
		ID:          "uuid-123",
		Username:    inputUsername,
		Password:    hashedPassword,
		Email:       "test@example.com",
		DisplayName: "Test User",
	}

	testCases := []struct {
		name            string
		setupMockRepo   func(ctx context.Context) *repoMocks.Repository
		setupMockHasher func() *mocks.Hasher
		setupMockJWTGen func() *jwtMocks.JWTGenerator
		expectedToken   string
		expectedErr     error
	}{
		{
			name: "happy path",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByUsername", ctx, inputUsername).Return(testUser, nil)
				return mockRepo
			},
			setupMockHasher: func() *mocks.Hasher {
				mockHasher := mocks.NewHasher(t)
				mockHasher.On("Compare", hashedPassword, inputPassword).Return(true)
				return mockHasher
			},
			setupMockJWTGen: func() *jwtMocks.JWTGenerator {
				mockJWTGen := jwtMocks.NewJWTGenerator(t)
				mockJWTGen.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == "uuid-123" && claims["email"] == "test@example.com"
				})).Return(jwtToken, nil)
				return mockJWTGen
			},
			expectedToken: jwtToken,
			expectedErr:   nil,
		},

		{
			name: "user not found",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByUsername", ctx, inputUsername).Return(nil, dbutils.ErrRecordNotFound)
				return mockRepo
			},
			setupMockHasher: func() *mocks.Hasher {
				return mocks.NewHasher(t)
			},
			setupMockJWTGen: func() *jwtMocks.JWTGenerator {
				return jwtMocks.NewJWTGenerator(t)
			},
			expectedToken: "",
			expectedErr:   ErrInvalidCredentials,
		},

		{
			name: "wrong password",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByUsername", ctx, inputUsername).Return(testUser, nil)
				return mockRepo
			},
			setupMockHasher: func() *mocks.Hasher {
				mockHasher := mocks.NewHasher(t)
				mockHasher.On("Compare", hashedPassword, inputPassword).Return(false)
				return mockHasher
			},
			setupMockJWTGen: func() *jwtMocks.JWTGenerator {
				return jwtMocks.NewJWTGenerator(t)
			},
			expectedToken: "",
			expectedErr:   ErrInvalidCredentials,
		},

		{
			name: "JWT generation fail",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetUserByUsername", ctx, inputUsername).Return(testUser, nil)
				return mockRepo
			},
			setupMockHasher: func() *mocks.Hasher {
				mockHasher := mocks.NewHasher(t)
				mockHasher.On("Compare", hashedPassword, inputPassword).Return(true)
				return mockHasher
			},
			setupMockJWTGen: func() *jwtMocks.JWTGenerator {
				mockJWTGen := jwtMocks.NewJWTGenerator(t)
				mockJWTGen.On("GenerateJWT", mock.MatchedBy(func(claims jwt.MapClaims) bool {
					return claims["sub"] == "uuid-123" && claims["email"] == "test@example.com"
				})).Return("", errTest)
				return mockJWTGen
			},
			expectedToken: "",
			expectedErr:   errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)
			mockHasher := tc.setupMockHasher()
			mockJWTGen := tc.setupMockJWTGen()

			svc := NewService(mockRepo, mockHasher, mockJWTGen)

			token, err := svc.Login(ctx, inputUsername, inputPassword)
			assert.Equal(t, tc.expectedToken, token)
			assert.ErrorIs(t, tc.expectedErr, err)
		})
	}
}
