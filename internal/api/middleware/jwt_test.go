package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/pkg/jwtutils/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTAuth(t *testing.T) {
	errTest := errors.New("test error")
	t.Parallel()
	testCases := []struct {
		name               string
		inputAuthHeader    string
		claims             jwt.MapClaims
		mockJWTValidator   func(claims jwt.MapClaims) *mocks.JWTValidator
		expectedStatusCode int
	}{
		{
			name:            "happy path",
			inputAuthHeader: "Bearer test",
			claims: jwt.MapClaims{
				"sub": "deb745af-1a62-4efa-99a0-f06b274bd993",
			},
			mockJWTValidator: func(claims jwt.MapClaims) *mocks.JWTValidator {
				mock := mocks.NewJWTValidator(t)
				mock.On("ValidateJWT", "test").Return(claims, nil)
				return mock
			},
			expectedStatusCode: http.StatusOK,
		},

		{
			name:            "error case - invalid token",
			inputAuthHeader: "Bearer test",
			mockJWTValidator: func(claims jwt.MapClaims) *mocks.JWTValidator {
				mock := mocks.NewJWTValidator(t)
				mock.On("ValidateJWT", "test").Return(nil, errTest)
				return mock
			},
			expectedStatusCode: http.StatusUnauthorized,
		},

		{
			name:            "error case - missing header",
			inputAuthHeader: "",
			mockJWTValidator: func(claims jwt.MapClaims) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},

		{
			name:            "error case - wrong header format",
			inputAuthHeader: "bruh",

			mockJWTValidator: func(claims jwt.MapClaims) *mocks.JWTValidator {
				return mocks.NewJWTValidator(t)
			},
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			_, api := gin.CreateTestContext(rec)
			jwtValidator := tc.mockJWTValidator(tc.claims)
			middleware := NewJWTAuth(jwtValidator)

			api.Use(middleware.JWTAuth())
			api.GET("/ping", func(ctx *gin.Context) {
				ctx.JSON(http.StatusOK, gin.H{"message": "pong"})
			})

			req := httptest.NewRequest(http.MethodGet, "/ping", nil)
			if tc.inputAuthHeader != "" {
				req.Header.Set("Authorization", tc.inputAuthHeader)
			}
			api.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
		})
	}
}
