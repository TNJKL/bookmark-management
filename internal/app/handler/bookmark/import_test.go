package bookmark

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/service/queue"
	"github.com/TNJKL/bookmark-management/internal/app/service/queue/mocks"
	"github.com/TNJKL/bookmark-management/pkg/csv"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_ImportBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		setupTestRequest   func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer)
		mockQueueSetup     func(ctx context.Context) *mocks.Service
		fileContent        string
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupTestRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/test", body)
				ctx.Set("claims", jwt.MapClaims{"sub": "1d04d961-3a76-4a9f-9758-b03b196b4867"})
				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},
			mockQueueSetup: func(ctx context.Context) *mocks.Service {
				serviceMock := mocks.NewService(t)
				serviceMock.On("SendImportBookmarkJob", ctx, "1d04d961-3a76-4a9f-9758-b03b196b4867", []*queue.ImportBookmarkInput{{Description: "Google", URL: "https://google.com"}, {Description: "Facebook", URL: "https://facebook.com"}, {Description: "Amazon", URL: "https://amazon.com"}}).Return(nil)
				return serviceMock
			},
			fileContent:        "description,url\nGoogle,https://google.com\nFacebook,https://facebook.com\nAmazon,https://amazon.com",
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `{"message":"Imported bookmarks successfully"}`,
		},
		{
			name: "error case - invalid csv",
			setupTestRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/test", body)
				ctx.Set("claims", jwt.MapClaims{"sub": "1d04d961-3a76-4a9f-9758-b03b196b4867"})
				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},
			mockQueueSetup: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			fileContent:        "description,url\nGoogle,1234567NgayTroi\nFacebook,https://facebook.com\nAmazon,https://amazon.com",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Input error","details":["URL is invalid (url)"]}`,
		},
		{
			name: "error case - no file attached",
			setupTestRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/test", nil)
				ctx.Set("claims", jwt.MapClaims{"sub": "1d04d961-3a76-4a9f-9758-b03b196b4867"})
				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},
			mockQueueSetup: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			fileContent:        "",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message":"Invalid file"}`,
		},
		{
			name: "error case - queue service error",
			setupTestRequest: func(ctx *gin.Context, body *bytes.Buffer, writer *multipart.Writer) {
				ctx.Request = httptest.NewRequest(http.MethodPost, "/test", body)
				ctx.Set("claims", jwt.MapClaims{"sub": "1d04d961-3a76-4a9f-9758-b03b196b4867"})
				ctx.Request.Header.Set("Content-Type", writer.FormDataContentType())
			},
			mockQueueSetup: func(ctx context.Context) *mocks.Service {
				serviceMock := mocks.NewService(t)
				serviceMock.On("SendImportBookmarkJob", ctx, "1d04d961-3a76-4a9f-9758-b03b196b4867", []*queue.ImportBookmarkInput{{Description: "Google", URL: "https://google.com"}, {Description: "Facebook", URL: "https://facebook.com"}, {Description: "Amazon", URL: "https://amazon.com"}}).Return(errTest)
				return serviceMock
			},
			fileContent:        "description,url\nGoogle,https://google.com\nFacebook,https://facebook.com\nAmazon,https://amazon.com",
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `{"message":"Processing error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			writer, body := csv.CreateTestMultipartRequest(t, tc.fileContent)
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			tc.setupTestRequest(ctx, body, writer)
			handler := NewHandler(nil, tc.mockQueueSetup(ctx))
			handler.ImportBookmarks(ctx)
			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			assert.Equal(t, tc.expectedResponse, rec.Body.String())
		})
	}
}
