package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestShortenURL_ShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		setupMockSvc       func() *mocks.ShortenURL
		setupTestRequest   func(ctx *gin.Context)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupMockSvc: func() *mocks.ShortenURL {
				mockSvc := mocks.NewShortenURL(t)
				mockSvc.On("CreateShortenLink", mock.Anything, mock.Anything, mock.Anything).Return("Songoku", nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"https://google.com","exp":100000000}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"code":"Songoku"`,
		},

		{
			name: "invalid input",
			setupMockSvc: func() *mocks.ShortenURL {
				//ko cần setup gì ShouldByJSON lỗi thì ko gọi tới service
				return mocks.NewShortenURL(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString("invalid json body")
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"error":"Invalid input"}`,
		},
		{
			name: "service error",
			setupMockSvc: func() *mocks.ShortenURL {
				mockSvc := mocks.NewShortenURL(t)
				mockSvc.On("CreateShortenLink", mock.Anything, mock.Anything, mock.Anything).Return("", testError)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"https://google.com","exp":100000000}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
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
			shortenHandler := NewShortenURL(mockSvc)
			shortenHandler.ShortenLink(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)
		})
	}

}
