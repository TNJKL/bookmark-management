package bookmark

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestBookmarkHandler_UpdateBookmark(t *testing.T) {
	t.Parallel()
	userID := "uuid-123"
	bookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	description := "Google"
	url := "https://www.google.com"

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
				mockSvc.On("UpdateBookmark", ctx, bookmarkID, userID, description, url).Return(nil)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Params = gin.Params{{Key: "id", Value: bookmarkID}}
				input := updateBookmarkBody{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, rec *httptest.ResponseRecorder) {
				var msg response.Message
				err := json.Unmarshal(rec.Body.Bytes(), &msg)
				assert.NoError(t, err)
				assert.Equal(t, "Success", msg.Message)
			},
		},
		{
			name: "invalid input - invalid url",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Params = gin.Params{{Key: "id", Value: bookmarkID}}
				input := updateBookmarkBody{
					Description: description,
					Url:         "invalid-url",
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "invalid input - invalid uuid param",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Params = gin.Params{{Key: "id", Value: "not-a-uuid"}}
				input := updateBookmarkBody{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/not-a-uuid", bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "bookmark not found",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("UpdateBookmark", ctx, bookmarkID, userID, description, url).Return(dbutils.ErrRecordNotFound)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Params = gin.Params{{Key: "id", Value: bookmarkID}}
				input := updateBookmarkBody{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(bodyBytes))
				ctx.Request.Header.Set("Content-Type", "application/json")
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `"message":"Bookmark not found"`,
		},
		{
			name: "internal server error",
			setupMockSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("UpdateBookmark", ctx, bookmarkID, userID, description, url).Return(errTest)
				return mockSvc
			},
			setupTestRequest: func(ctx *gin.Context) {
				ctx.Set("claims", jwt.MapClaims{"sub": userID})
				ctx.Params = gin.Params{{Key: "id", Value: bookmarkID}}
				input := updateBookmarkBody{
					Description: description,
					Url:         url,
				}
				bodyBytes, _ := json.Marshal(input)
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(bodyBytes))
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
				ctx.Request = httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, nil)
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
			bookmarkHandler.UpdateBookmark(ctx)

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
