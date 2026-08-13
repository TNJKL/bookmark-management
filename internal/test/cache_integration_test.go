package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/app/handler/dto"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/app/service/bookmark"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/TNJKL/bookmark-management/pkg/response"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type createBookmarkInput struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

type updateBookmarkInput struct {
	Description string `json:"description"`
	Url         string `json:"url"`
}

func TestBookmarkCache_GetBookmarks(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)
	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"
	token := generateTestToken(t, jwtGen, userID, email)

	cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
	cacheKey := "1_10"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupRedis         func(t *testing.T) *redis.Client
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		verifyFunc         func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "happy path - cache hit returns cached result directly",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				cachedResult := bookmark.GetBookmarksResult{
					Bookmarks: []*model.Bookmark{
						{
							Description: "Cached Bookmark Item 1",
							URL:         "https://cached-item1.com",
							Code:        "CACHE01",
						},
						{
							Description: "Cached Bookmark Item 2",
							URL:         "https://cached-item2.com",
							Code:        "CACHE02",
						},
					},
					Total: 2,
				}
				cacheBytes, err := json.Marshal(cachedResult)
				assert.NoError(t, err)

				err = mockR.HSet(t.Context(), cacheGroupKey, cacheKey, cacheBytes).Err()
				assert.NoError(t, err)
				return mockR
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder) {
				var resp dto.SuccessResponse[[]*model.Bookmark]
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Data, 2)
				assert.Equal(t, "Cached Bookmark Item 1", resp.Data[0].Description)
				assert.Equal(t, "https://cached-item1.com", resp.Data[0].URL)
				assert.Equal(t, "Cached Bookmark Item 2", resp.Data[1].Description)
				assert.Equal(t, "https://cached-item2.com", resp.Data[1].URL)
				assert.NotNil(t, resp.Pagination)
				assert.Equal(t, int64(2), resp.Pagination.Total)

				var dbCount int64
				db.Model(&model.Bookmark{}).Count(&dbCount)
				assert.Equal(t, int64(0), dbCount)
			},
		},
		{
			name: "happy path - cache miss queries DB and populates Redis cache",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder) {
				ctx := t.Context()

				var resp dto.SuccessResponse[[]*model.Bookmark]
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Len(t, resp.Data, 5)
				assert.NotNil(t, resp.Pagination)
				assert.Equal(t, int64(5), resp.Pagination.Total)

				cachedData, err := r.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
				assert.NoError(t, err)
				assert.NotEmpty(t, cachedData)

				var cachedResult bookmark.GetBookmarksResult
				err = json.Unmarshal(cachedData, &cachedResult)
				assert.NoError(t, err)
				assert.Equal(t, int64(5), cachedResult.Total)
				assert.Len(t, cachedResult.Bookmarks, 5)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			redisClient := tc.setupRedis(t)
			testAPI := buildTestAPI(db, redisClient, jwtGen, jwtVal)
			rec := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, redisClient, rec)
			}
			_ = ctx
		})
	}
}

func TestBookmarkCache_CreateBookmark(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)
	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"
	token := generateTestToken(t, jwtGen, userID, email)

	cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupRedis         func(t *testing.T) *redis.Client
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		verifyFunc         func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "create bookmark invalidates user cache group",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)

				err := mockR.HSet(t.Context(), cacheGroupKey, "1_10", []byte(`{"total":5}`)).Err()
				assert.NoError(t, err)
				return mockR
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := createBookmarkInput{
					Description: "Cache Invalidation Test Item",
					Url:         "https://cache-invalidation.com",
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
			verifyFunc: func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder) {
				ctx := t.Context()

				var resp dto.SuccessResponse[*model.Bookmark]
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Create a bookmark successfully!", resp.Message)
				assert.NotNil(t, resp.Data)
				assert.Equal(t, "Cache Invalidation Test Item", resp.Data.Description)
				assert.Equal(t, "https://cache-invalidation.com", resp.Data.URL)

				exists, err := r.Exists(ctx, cacheGroupKey).Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), exists)

				var count int64
				db.Model(&model.Bookmark{}).Where("url = ?", "https://cache-invalidation.com").Count(&count)
				assert.Equal(t, int64(1), count)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)
			redisClient := tc.setupRedis(t)
			testAPI := buildTestAPI(db, redisClient, jwtGen, jwtVal)
			rec := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, redisClient, rec)
			}
		})
	}
}

func TestBookmarkCache_UpdateBookmark(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)
	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"
	token := generateTestToken(t, jwtGen, userID, email)

	cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
	targetBookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd993"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupRedis         func(t *testing.T) *redis.Client
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		verifyFunc         func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "update bookmark invalidates user cache group",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				err := mockR.HSet(t.Context(), cacheGroupKey, "1_10", []byte(`{"total":5}`)).Err()
				assert.NoError(t, err)
				return mockR
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := updateBookmarkInput{
					Description: "Updated Bookmark Description",
					Url:         "https://updated-url.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPut, "/v1/bookmarks/"+targetBookmarkID, bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder) {
				ctx := t.Context()

				var resp response.Message
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Success", resp.Message)

				exists, err := r.Exists(ctx, cacheGroupKey).Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), exists)

				var updatedBm model.Bookmark
				err = db.First(&updatedBm, "id = ?", targetBookmarkID).Error
				assert.NoError(t, err)
				assert.Equal(t, "Updated Bookmark Description", updatedBm.Description)
				assert.Equal(t, "https://updated-url.com", updatedBm.URL)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)
			redisClient := tc.setupRedis(t)
			testAPI := buildTestAPI(db, redisClient, jwtGen, jwtVal)
			rec := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, redisClient, rec)
			}
		})
	}
}

func TestBookmarkCache_DeleteBookmark(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)
	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	email := "johndoe@example.com"
	token := generateTestToken(t, jwtGen, userID, email)

	cacheGroupKey := fmt.Sprintf("get_bookmarks_%s", userID)
	targetBookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd993"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupRedis         func(t *testing.T) *redis.Client
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		verifyFunc         func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "delete bookmark invalidates user cache group",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				err := mockR.HSet(t.Context(), cacheGroupKey, "1_10", []byte(`{"total":5}`)).Err()
				assert.NoError(t, err)
				return mockR
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodDelete, "/v1/bookmarks/"+targetBookmarkID, nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, db *gorm.DB, r *redis.Client, recorder *httptest.ResponseRecorder) {
				ctx := t.Context()

				// Parse and verify exact response
				var resp response.Message
				err := json.Unmarshal(recorder.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, "Success", resp.Message)

				// Verify Redis cache key is DELETED
				exists, err := r.Exists(ctx, cacheGroupKey).Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), exists)

				// Verify record is deleted from DB
				var count int64
				db.Model(&model.Bookmark{}).Where("id = ?", targetBookmarkID).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)
			redisClient := tc.setupRedis(t)
			testAPI := buildTestAPI(db, redisClient, jwtGen, jwtVal)
			rec := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, rec.Code)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db, redisClient, rec)
			}
		})
	}
}
