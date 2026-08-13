package integration

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestRateLimit_Integration(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)
	userA := "deb745af-1a62-4efa-99a0-f06b274bd990"
	userB := "deb745af-1a62-4efa-99a0-f06b274bd991"
	tokenA := generateTestToken(t, jwtGen, userA, "usera@example.com")
	tokenB := generateTestToken(t, jwtGen, userB, "userb@example.com")

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupRedis         func(t *testing.T) *redis.Client
		setupTestHTTP      func(testAPI api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		verifyFunc         func(t *testing.T, r *redis.Client, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "happy path - 20 consecutive requests pass within threshold",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(testAPI api.Engine) *httptest.ResponseRecorder {
				var rec *httptest.ResponseRecorder
				for i := 1; i <= 20; i++ {
					req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
					req.Header.Set("Authorization", "Bearer "+tokenA)
					rec = httptest.NewRecorder()
					testAPI.ServerHTTP(rec, req)
				}
				return rec
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, r *redis.Client, recorder *httptest.ResponseRecorder) {
				key := fmt.Sprintf("rate_limit:%s", userA)
				val, err := r.Get(t.Context(), key).Int()
				assert.NoError(t, err)
				assert.Equal(t, 20, val)
			},
		},

		{
			name: "error case - 21th request exceeds rate limit and is blocked",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(testAPI api.Engine) *httptest.ResponseRecorder {
				var rec *httptest.ResponseRecorder
				for i := 1; i <= 21; i++ {
					req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
					req.Header.Set("Authorization", "Bearer "+tokenA)
					rec = httptest.NewRecorder()
					testAPI.ServerHTTP(rec, req)
				}
				return rec
			},
			expectedStatusCode: http.StatusTooManyRequests,
			verifyFunc: func(t *testing.T, r *redis.Client, recorder *httptest.ResponseRecorder) {
				assert.Contains(t, recorder.Body.String(), "rate limit exceeded")
			},
		},
		{
			name: "isolation - user B is not blocked when user A exceeds limit",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{})
			},
			setupRedis: func(t *testing.T) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			setupTestHTTP: func(testAPI api.Engine) *httptest.ResponseRecorder {
				for i := 1; i <= 20; i++ {
					req := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
					req.Header.Set("Authorization", "Bearer "+tokenA)
					rec := httptest.NewRecorder()
					testAPI.ServerHTTP(rec, req)
				}
				reqUserB := httptest.NewRequest(http.MethodGet, "/v1/bookmarks?page=1&limit=10", nil)
				reqUserB.Header.Set("Authorization", "Bearer "+tokenB)
				recUserB := httptest.NewRecorder()
				testAPI.ServerHTTP(recUserB, reqUserB)
				return recUserB
			},
			expectedStatusCode: http.StatusOK,
			verifyFunc: func(t *testing.T, r *redis.Client, recorder *httptest.ResponseRecorder) {
				keyB := fmt.Sprintf("rate_limit:%s", userB)
				valB, err := r.Get(t.Context(), keyB).Int()
				assert.NoError(t, err)
				assert.Equal(t, 1, valB)
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
				tc.verifyFunc(t, redisClient, rec)
			}
		})
	}
}
