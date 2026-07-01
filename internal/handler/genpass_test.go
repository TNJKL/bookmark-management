package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testErr = errors.New("test error")

func TestGenPass_GeneratePassword(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMockSvc     func() *mocks.GenPass
		setupTestRequest func(ctx *gin.Context)

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "success",
			setupMockSvc: func() *mocks.GenPass {
				mockSvc := mocks.NewGenPass(t)
				mockSvc.On("GeneratePassword", passwordLength).Return("123456789012", nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"password":"123456789012"}`,
		},

		{
			name: "error case",
			setupMockSvc: func() *mocks.GenPass {
				mockSvc := mocks.NewGenPass(t)
				mockSvc.On("GeneratePassword", passwordLength).Return("", testErr)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/genpass", nil)
			},

			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupTestRequest(ctx)

			mockSvc := tc.setupMockSvc()
			genPassHandler := NewGenPass(mockSvc)
			genPassHandler.GeneratePassword(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}
}
