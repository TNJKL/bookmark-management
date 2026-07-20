package handler

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/service/mocks"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestShortenURL_ShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		setupMockSvc       func(ctx context.Context) *mocks.ShortenURL
		setupTestRequest   func(ctx *gin.Context)
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenURL {
				mockSvc := mocks.NewShortenURL(t)
				mockSvc.On("CreateShortenLink", ctx, "https://google.com", int64(50000)).Return("Songoku", nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"https://google.com","exp":50000}`)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"code":"Songoku"`,
		},

		{
			name: "invalid input",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenURL {
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
			setupMockSvc: func(ctx context.Context) *mocks.ShortenURL {
				mockSvc := mocks.NewShortenURL(t)
				mockSvc.On("CreateShortenLink", ctx, "https://google.com", int64(50000)).Return("", testError)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				body := bytes.NewBufferString(`{"url":"https://google.com","exp":50000}`)
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
			mockSvc := tc.setupMockSvc(ctx)
			shortenHandler := NewShortenURL(mockSvc)
			shortenHandler.ShortenLink(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tc.expectedResponse)
		})
	}

}

func TestShortenURL_Redirect(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		setupMockSvc       func(ctx context.Context) *mocks.ShortenURL
		setupTestRequest   func(ctx *gin.Context)
		expectedStatusCode int
		expectedURL        string
	}{
		{
			name: "happy path",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenURL {
				mockSvc := mocks.NewShortenURL(t)
				mockSvc.On("GetLinkFromCode", ctx, "Songoku").Return("https://google.com", nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/Songoku", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "Songoku"}}
			},
			expectedStatusCode: http.StatusFound,
			expectedURL:        "https://google.com",
		},
		{
			name: "errros case",
			setupMockSvc: func(ctx context.Context) *mocks.ShortenURL {
				mockSvc := mocks.NewShortenURL(t)
				mockSvc.On("GetLinkFromCode", ctx, "Songoku").Return("", testError)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/links/redirect/Songoku", nil)
				ctx.Params = gin.Params{{Key: "code", Value: "Songoku"}}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedURL:        "",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupTestRequest(ctx)
			mockSvc := tc.setupMockSvc(ctx)
			shortenHandler := NewShortenURL(mockSvc)
			shortenHandler.Redirect(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedURL, rec.Header().Get("Location"))
		})
	}

}
