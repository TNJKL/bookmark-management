package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/repository/ratelimit/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimit(t *testing.T) {
	errTest := errors.New("redis connection error")
	t.Parallel()

	testCases := []struct {
		name                 string
		claims               jwt.MapClaims
		mockRateLimitRepo    func() *mocks.Repo
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "happy path",
			claims: jwt.MapClaims{
				"sub": "deb745af-1a62-4efa-99a0-f06b274bd993",
			},
			mockRateLimitRepo: func() *mocks.Repo {
				mockRepo := mocks.NewRepo(t)
				key := "rate_limit:deb745af-1a62-4efa-99a0-f06b274bd993"
				mockRepo.On("GetCurrentRateLimit", mock.AnythingOfType("*gin.Context"), key).Return(15, nil)
				mockRepo.On("IncreaseRateLimit", mock.AnythingOfType("*gin.Context"), key, 10*time.Second).Return()
				return mockRepo
			},
			expectedStatusCode: http.StatusOK,
		},
		{
			name: "error case - rate limit exceeded",
			claims: jwt.MapClaims{
				"sub": "deb745af-1a62-4efa-99a0-f06b274bd993",
			},
			mockRateLimitRepo: func() *mocks.Repo {
				mockRepo := mocks.NewRepo(t)
				key := "rate_limit:deb745af-1a62-4efa-99a0-f06b274bd993"
				mockRepo.On("GetCurrentRateLimit", mock.AnythingOfType("*gin.Context"), key).Return(20, nil)
				return mockRepo
			},
			expectedStatusCode:   http.StatusTooManyRequests,
			expectedResponseBody: `{"message":"rate limit exceeded"}`,
		},

		{
			name:   "error case - missing user claims",
			claims: nil,
			mockRateLimitRepo: func() *mocks.Repo {
				return mocks.NewRepo(t)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},

		{
			name: "edge case - redis repo error fallback",
			claims: jwt.MapClaims{
				"sub": "deb745af-1a62-4efa-99a0-f06b274bd993",
			},
			mockRateLimitRepo: func() *mocks.Repo {
				mockRepo := mocks.NewRepo(t)
				key := "rate_limit:deb745af-1a62-4efa-99a0-f06b274bd993"
				mockRepo.On("GetCurrentRateLimit", mock.AnythingOfType("*gin.Context"), key).Return(0, errTest)
				mockRepo.On("IncreaseRateLimit", mock.AnythingOfType("*gin.Context"), key, 10*time.Second).Return()
				return mockRepo
			},
			expectedStatusCode: http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			_, api := gin.CreateTestContext(rec)
			if tc.claims != nil {
				api.Use(func(ctx *gin.Context) {
					ctx.Set("claims", tc.claims)
					ctx.Next()
				})
			}
			mockRepo := tc.mockRateLimitRepo()
			middleware := NewRateLimit(mockRepo)
			api.Use(middleware.RateLimit())
			api.GET("/ping", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
			})

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			api.ServeHTTP(rec, req)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponseBody != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedResponseBody)
			}
		})
	}
}
