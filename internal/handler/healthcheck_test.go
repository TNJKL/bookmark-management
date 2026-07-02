package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var testError = errors.New("test error")

func TestHealthCheck_HealthCheck(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupMockSvc     func() *mocks.HealthChecker
		setupTestRequest func(ctx *gin.Context)

		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupMockSvc: func() *mocks.HealthChecker {
				mockSvc := mocks.NewHealthChecker(t)
				mockSvc.On("HealthCheck").
					Return(model.HealthCheckResponse{
						"OK", "bookmark_service", "67026e45-34fd-449c-aa29-d18c7686ab00",
					}, nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"OK","service_name":"bookmark_service","instance_id":"67026e45-34fd-449c-aa29-d18c7686ab00"}`,
		},
		{
			name: "error case",
			setupMockSvc: func() *mocks.HealthChecker {
				mockSvc := mocks.NewHealthChecker(t)
				mockSvc.On("HealthCheck").
					Return(model.HealthCheckResponse{
						"", "", "",
					}, testError)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/healthcheck", nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"error":"Internal Server Error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			//gọi Parallel để chạy song song các test cases
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)

			tc.setupTestRequest(ctx)

			mockSvc := tc.setupMockSvc()

			healthCheckHandler := NewHealthCheck(mockSvc)
			healthCheckHandler.HealthCheck(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())

		})
	}

}
