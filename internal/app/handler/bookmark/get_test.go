package bookmark

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/handler/dto"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	serviceBookmark "github.com/TNJKL/bookmark-management/internal/app/service/bookmark"
	"github.com/TNJKL/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_GetBookmarks(t *testing.T) {
	t.Parallel()
	userID := "uuid-123"
	errTest := errors.New("test error")

	mockBookmarks := []*model.Bookmark{
		{
			Base: model.Base{
				ID: "bookmark-1",
			},
			Description: "Google Search",
			URL:         "https://google.com",
			Code:        "12345678",
			UserID:      userID,
		},
		{
			Base: model.Base{
				ID: "bookmark-2",
			},
			Description: "GitHub",
			URL:         "https://github.com",
			Code:        "87654321",
			UserID:      userID,
		},
	}

	testCases := []struct {
		name               string
		setupMockSvc       func(ctx context.Context) *mocks.Service
		setupTestRequest   func(ctx *gin.Context)
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, rec *httptest.ResponseRecorder)
	}{
		{
			name: "happy path",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("GetBookmarks", ctx, userID, 1, 10).Return(&serviceBookmark.GetBookmarksResult{
					Bookmarks: mockBookmarks,
					Total:     2,
				}, nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp dto.SuccessResponse[[]*model.Bookmark]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Len(t, resp.Data, 2)
				assert.Equal(t, "bookmark-1", resp.Data[0].ID)
				assert.Equal(t, "Google Search", resp.Data[0].Description)
				assert.Equal(t, "https://google.com", resp.Data[0].URL)

				assert.NotNil(t, resp.Pagination)
				assert.Equal(t, 1, resp.Pagination.Page)
				assert.Equal(t, 10, resp.Pagination.Limit)
				assert.Equal(t, int64(2), resp.Pagination.Total)
			},
		},
		{
			name: "invalid input - page less than 1",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=0&limit=10", nil)
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "internal server error",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("GetBookmarks", ctx, userID, 1, 10).Return(nil, errTest)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   `"message":"Processing error"`,
		},
		{
			name: "unauthorized - no token",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Request = httptest.NewRequest(http.MethodGet, "/v1/bookmarks", nil)
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Invalid token"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			rec := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(rec)
			mockSvc := tc.setupMockSvc(ctx)
			tc.setupTestRequest(ctx)

			bookmarkHandler := NewHandler(mockSvc)
			bookmarkHandler.GetBookmarks(ctx)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, rec.Body.String(), tc.expectedResponse)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, rec)
			}
		})
	}
}
