package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/app/handler/dto"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/csv"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type createBookmarkInputBody struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

func TestCreateBookmarkEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, db *gorm.DB, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				input := createBookmarkInputBody{
					Description: "My Custom Bookmark",
					Url:         "https://custom-bookmark.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, db *gorm.DB, recorder *httptest.ResponseRecorder) {
				var resp dto.SuccessResponse[*model.Bookmark]
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Create a bookmark successfully!", resp.Message)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, "My Custom Bookmark", resp.Data.Description)
				assert.Equal(t, "https://custom-bookmark.com", resp.Data.URL)
				assert.NotEmpty(t, resp.Data.Code)

				var bm model.Bookmark
				err = db.First(&bm, "user_id = ? AND url = ?", userID, "https://custom-bookmark.com").Error
				assert.NoError(t, err)
				assert.Equal(t, "My Custom Bookmark", bm.Description)
				assert.Equal(t, "https://custom-bookmark.com", bm.URL)
				assert.NotEmpty(t, bm.Code)
			},
		},
		{
			name: "invalid input - missing url",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				input := createBookmarkInputBody{
					Description: "Missing URL",
					Url:         "",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
			verifyFunc: func(t *testing.T, db *gorm.DB, recorder *httptest.ResponseRecorder) {
				var count int64
				db.Model(&model.Bookmark{}).Where("description = ?", "Missing URL").Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "invalid input - invalid url format",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				input := createBookmarkInputBody{
					Description: "Invalid URL Format",
					Url:         "invalid-url-string",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
			verifyFunc: func(t *testing.T, db *gorm.DB, recorder *httptest.ResponseRecorder) {
				var count int64
				db.Model(&model.Bookmark{}).Where("description = ?", "Invalid URL Format").Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "unauthorized - no token",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := createBookmarkInputBody{
					Description: "No Token Test",
					Url:         "https://example.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Unauthorized"`,
			verifyFunc: func(t *testing.T, db *gorm.DB, recorder *httptest.ResponseRecorder) {
				var count int64
				db.Model(&model.Bookmark{}).Where("description = ?", "No Token Test").Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)

			testAPI := buildTestAPI(db, nil, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, recorder)
			}
		})
	}
}

func TestGetBookmarksEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "happy path - get all user bookmarks",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var resp dto.SuccessResponse[[]*model.Bookmark]
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 5, len(resp.Data))
				assert.NotNil(t, resp.Pagination)
				assert.Equal(t, int64(5), resp.Pagination.Total)
				assert.Equal(t, 1, resp.Pagination.Page)
				assert.Equal(t, 10, resp.Pagination.Limit)
			},
		},
		{
			name: "happy path - pagination limit",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=2", nil)
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				var resp dto.SuccessResponse[[]*model.Bookmark]
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, 2, len(resp.Data))
				assert.NotNil(t, resp.Pagination)
				assert.Equal(t, int64(5), resp.Pagination.Total)
				assert.Equal(t, 1, resp.Pagination.Page)
				assert.Equal(t, 2, resp.Pagination.Limit)
			},
		},
		{
			name: "invalid query parameter - page 0",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=0&limit=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "unauthorized - no token",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Unauthorized"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)

			testAPI := buildTestAPI(db, nil, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, recorder)
			}
		})
	}
}

type updateBookmarkInputBody struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

func TestUpdateBookmarkEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"
	bookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd993"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				input := updateBookmarkInputBody{
					Description: "Google Updated",
					Url:         "https://www.google.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Success"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var bm model.Bookmark
				err := db.First(&bm, "id = ?", bookmarkID).Error
				assert.NoError(t, err)
				assert.Equal(t, "Google Updated", bm.Description)
				assert.Equal(t, "https://www.google.com", bm.URL)
			},
		},
		{
			name: "bookmark not found - wrong id",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				nonExistentID := "deb745af-1a62-4efa-99a0-f06b274bd999"
				input := updateBookmarkInputBody{
					Description: "Google",
					Url:         "https://www.google.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+nonExistentID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `"message":"Bookmark not found"`,
		},
		{
			name: "invalid input - missing url",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				input := updateBookmarkInputBody{
					Description: "Missing URL",
					Url:         "",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
		{
			name: "unauthorized - no token",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := updateBookmarkInputBody{
					Description: "Google",
					Url:         "https://www.google.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+bookmarkID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Unauthorized"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)

			testAPI := buildTestAPI(db, nil, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db)
			}
		})
	}
}

func TestDeleteBookmarkEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"
	bookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd993"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+bookmarkID, nil)
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Success"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&model.Bookmark{}).Where("id = ?", bookmarkID).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "bookmark not found - wrong id",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				nonExistentID := "deb745af-1a62-4efa-99a0-f06b274bd999"
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+nonExistentID, nil)
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse:   `"message":"Bookmark not found"`,
		},
		{
			name: "unauthorized - no token",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+bookmarkID, nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Unauthorized"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)

			testAPI := buildTestAPI(db, nil, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db)
			}
		})
	}
}

func TestImportBookmarksEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"

	testCases := []struct {
		name               string
		setupRedis         func(t *testing.T) *redis.Client
		setupTestHTTP      func(api api.Engine, redisClient *redis.Client) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, redisClient *redis.Client, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "happy path - import valid csv bookmarks",

			setupRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(api api.Engine, redisClient *redis.Client) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				csvContent := "description,url\nGoogle,https://google.com\nFacebook,https://facebook.com"
				writer, body := csv.CreateTestMultipartRequest(t, csvContent)

				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Imported bookmarks successfully"`,
			verifyFunc: func(t *testing.T, redisClient *redis.Client, recorder *httptest.ResponseRecorder) {
				poppedMsg, err := redisClient.RPop(context.Background(), "bookmark-import").Bytes()
				require.NoError(t, err)
				assert.Contains(t, string(poppedMsg), "https://google.com")
				assert.Contains(t, string(poppedMsg), userID)
			},
		},
		{
			name: "invalid input - invalid csv content",
			setupRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(api api.Engine, redisClient *redis.Client) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, userID, email)
				csvContent := "description,url\nGoogle,invalid-url-string"
				writer, body := csv.CreateTestMultipartRequest(t, csvContent)

				req := httptest.NewRequest(http.MethodPost, "/v1/bookmarks/import", body)
				req.Header.Set("Content-Type", writer.FormDataContentType())
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			redisClient := tc.setupRedis(t)

			testAPI := buildTestAPI(nil, redisClient, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI, redisClient)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			if tc.expectedResponse != "" {
				assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			}
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, redisClient, recorder)
			}
		})
	}
}
