package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/service/user/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Register(t *testing.T) {
	t.Parallel()
	username := "testuser"
	password := "Kakarot1996@"
	displayName := "Test User"
	email := "test@gmail.com"

	createdUser := &model.User{
		ID:          "deb745af-1a62-4efa-99a0-f06b274bd990",
		Username:    username,
		DisplayName: displayName,
		Email:       email,
	}

	testCases := []struct {
		name               string
		setupMockSvc       func(ctx context.Context) *mocks.Service
		setupTestRequest   func(ctx *gin.Context)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateUser", ctx, username, password, displayName, email).Return(createdUser, nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := registerInputBody{
					Username:    username,
					Password:    password,
					DisplayName: displayName,
					Email:       email,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"Register an user successfully!"`,
		},
		{
			name: "invalid input: password too short",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := registerInputBody{
					Username:    username,
					Password:    "shortpw",
					DisplayName: displayName,
					Email:       email,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "duplicate username",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateUser", ctx, username, password, displayName, email).Return(nil, dbutils.ErrDuplicationUsername)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := registerInputBody{
					Username:    username,
					Password:    password,
					DisplayName: displayName,
					Email:       email,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `"message":"Username already taken"`,
		},
		{
			name: "duplicate email",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateUser", ctx, username, password, displayName, email).Return(nil, dbutils.ErrDuplicationEmail)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := registerInputBody{
					Username:    username,
					Password:    password,
					DisplayName: displayName,
					Email:       email,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `"message":"Email already taken"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			mockSvc := tc.setupMockSvc(ctx)
			tc.setupTestRequest(ctx)
			userHandler := NewHandler(mockSvc)
			userHandler.Register(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)
		})
	}
}
