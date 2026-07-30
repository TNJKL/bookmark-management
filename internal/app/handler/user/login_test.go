package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/service/user"
	"github.com/TNJKL/bookmark-management/internal/app/service/user/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_Login(t *testing.T) {
	t.Parallel()
	username := "testuser"
	password := "Kakarot1996@"
	jwtToken := "mock-jwt-token"
	errTest := errors.New("test error")

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
				mockSvc.On("Login", ctx, username, password).Return(jwtToken, nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := loginInputBody{
					Username: username,
					Password: password,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Logged in successfully","data":"mock-jwt-token"`,
		},

		{
			name: "invalid input - password too short",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := loginInputBody{
					Username: username,
					Password: "short-password",
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},

		{
			name: "invalid credentials",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("Login", ctx, username, password).Return("", user.ErrInvalidCredentials)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := loginInputBody{
					Username: username,
					Password: password,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"invalid credentials"`,
		},

		{
			name: "internal server error",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("Login", ctx, username, password).Return("", errTest)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				input := loginInputBody{
					Username: username,
					Password: password,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `"message":"Processing error"`,
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
			userHandler.Login(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)
		})
	}
}
