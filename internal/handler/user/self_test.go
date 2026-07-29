package user

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/service/user/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestUserHandler_GetSelfInfo(t *testing.T) {
	t.Parallel()
	uid := "uuid-123"
	errTest := errors.New("test error")
	testUser := &model.User{
		ID:          uid,
		Username:    "testuser",
		DisplayName: "Test User",
		Email:       "test@gmail.com",
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
				mockSvc.On("GetSelfInfo", ctx, uid).Return(testUser, nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"sub": uid,
				})
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"id":"uuid-123","display_name":"Test User","username":"testuser","email":"test@gmail.com"`,
		},

		{
			name: "Unauthorized - missing claims",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Invalid token"`,
		},

		{
			name: "Internal server error",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("GetSelfInfo", ctx, uid).Return(nil, errTest)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{
					"sub": uid,
				})
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
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
			userHandler.GetSelfInfo(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)
		})
	}
}

func TestUserHandler_UpdateSelfInfo(t *testing.T) {
	t.Parallel()
	uid := "uuid-123"
	newDisplayName := "New Name"
	newEmail := "new@example.com"
	//errTest := errors.New("test error")

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
				mockSvc.On("UpdateSelfInfo", ctx, uid, newDisplayName, newEmail).Return(nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": uid})
				input := updateSelfInfoBody{
					DisplayName: newDisplayName,
					Email:       newEmail,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Edit current user successfully!"`,
		},

		{
			name: "invalid input - wrong email format",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": uid})
				input := updateSelfInfoBody{
					DisplayName: newDisplayName,
					Email:       "bla-bruh-hahah",
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},

		{
			name: "unauthorized - no token",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Invalid token"`,
		},

		{
			name: "email already taken",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("UpdateSelfInfo", ctx, uid, newDisplayName, newEmail).Return(dbutils.ErrDuplicationEmail)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": uid})
				input := updateSelfInfoBody{
					DisplayName: newDisplayName,
					Email:       newEmail,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/self/info", bytes.NewBuffer(bodyBytes))
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
			userHandler.UpdateSelfInfo(ctx)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)
		})
	}
}
