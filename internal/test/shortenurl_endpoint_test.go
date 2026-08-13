package integration

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		setupRedis           func(ctx context.Context) *redis.Client
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "happy path",
			setupRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := bytes.NewBufferString(`{"url":"https://google.com","exp":1000000}`)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"code":`, //khúc này chỉ cần kiểm tra key "code" có tồn tại ko
		},
		{
			name: "invalid input",
			setupRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := bytes.NewBufferString("invalid test")
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusBadRequest,
			expectedResponseBody: `{"message":"Input error"}`,
		},
		{
			name: "wrong endpoint method",
			setupRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/shorten", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedResponseBody: ``,
		},
		{
			name: "redis connection error",
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				body := bytes.NewBufferString(`{"url":"https://google.com","exp":1000000}`)
				req := httptest.NewRequest(http.MethodPost, "/v1/links/shorten", body)
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusInternalServerError,
			expectedResponseBody: `{"message":"Processing error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			mockRedis := tc.setupRedis(ctx)
			testAPI := api.NewEngine(&api.EngineOpts{
				App:         gin.New(),
				Cfg:         &api.Config{},
				RedisClient: mockRedis,
			})
			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}

}

func TestRedirectEnpoint(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name                 string
		setupRedis           func(ctx context.Context) *redis.Client
		setupDB              func(t *testing.T) *gorm.DB
		setupTestHTTP        func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode   int
		expectedURL          string
		expectedResponseBody string
	}{
		{
			name: "redis happy path (redis code)",
			setupRedis: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				err := mock.Set(ctx, "aSongoku", "https://google.com", 0).Err()
				assert.NoError(t, err)
				return mock
			},
			setupDB: func(t *testing.T) *gorm.DB { return nil },
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/aSongoku", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusFound,
			expectedURL:          "https://google.com",
			expectedResponseBody: "",
		},
		{
			name: "redis code not found",
			setupRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupDB: func(t *testing.T) *gorm.DB { return nil },
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/ablahhhh", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedURL:          "",
			expectedResponseBody: `{"message":"Input error"}`,
		},

		{
			name:       "postgres happy path (bookmark code)",
			setupRedis: func(ctx context.Context) *redis.Client { return redisPkg.InitMockRedis(t) },
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
				db.Model(&model.Bookmark{}).Where("id = ?", "deb745af-1a62-4efa-99a0-f06b274bd993").Update("code", "k123456")
				return db
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/k123456", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusFound,
			expectedURL:          "https://google.com",
			expectedResponseBody: "",
		},
		{
			name:       "postgres bookmark code not found",
			setupRedis: func(ctx context.Context) *redis.Client { return redisPkg.InitMockRedis(t) },
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/links/redirect/kNotFound", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode:   http.StatusNotFound,
			expectedURL:          "",
			expectedResponseBody: `{"message":"Input error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRedis := tc.setupRedis(ctx)
			db := tc.setupDB(t)
			testAPI := api.NewEngine(&api.EngineOpts{
				App:         gin.New(),
				Cfg:         &api.Config{},
				RedisClient: mockRedis,
				Db:          db,
			})
			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Equal(t, tc.expectedURL, recorder.Header().Get("Location"))

		})
	}
}
