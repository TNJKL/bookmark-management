package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/handler/dto"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

func TestBookmarkHandler_CreateBookmark(t *testing.T) {
	t.Parallel()
	userID := "uuid-123"
	url := "https://google.com"
	description := "Google Search"
	code := "12345678"

	createdBookmark := &model.Bookmark{
		Base: model.Base{
			ID: "bookmark-1",
		},
		Description: description,
		URL:         url,
		Code:        code,
		UserID:      userID,
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
				mockSvc.On("CreateBookmark", ctx, description, url, userID).Return(createdBookmark, nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				input := createBookmarkInput{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var resp dto.SuccessResponse[*model.Bookmark]
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)

				assert.Equal(t, "Create a bookmark successfully!", resp.Message)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, "bookmark-1", resp.Data.ID)
				assert.Equal(t, description, resp.Data.Description)
				assert.Equal(t, url, resp.Data.URL)
				assert.Equal(t, code, resp.Data.Code)
			},
		},
		{
			name: "invalid input - invalid url",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				input := createBookmarkInput{
					Description: description,
					Url:         "invalid-url",
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "duplicate bookmark code",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateBookmark", ctx, description, url, userID).Return(nil, dbutils.ErrUniqueConstraint)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				input := createBookmarkInput{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `"message":"Bookmark already exists"`,
		},
		{
			name: "internal server error",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateBookmark", ctx, description, url, userID).Return(nil, errTest)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				input := createBookmarkInput{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
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
				ctx.Request = httptest.NewRequest(http.MethodPost, "/v1/bookmarks", nil)
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

			bookmarkHandler := NewHandler(mockSvc, nil)
			bookmarkHandler.CreateBookmark(ctx)

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
