package integration

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestShortenURLEndpoint(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupRedis    func(t *testing.T) *redis.Client
		setupTestHTTP func(api api.Engine) *httptest.ResponseRecorder

		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name: "happy path",
			setupRedis: func(t *testing.T) *redis.Client {
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
			setupRedis: func(t *testing.T) *redis.Client {
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
			expectedResponseBody: `{"error":"Invalid input"}`,
		},
		{
			name: "wrong endpoint method",
			setupRedis: func(t *testing.T) *redis.Client {
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
			setupRedis: func(t *testing.T) *redis.Client {
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
			expectedResponseBody: `{"error":"Internal Server Error"}`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			redisClient := tc.setupRedis(t)
			testAPI := api.NewEngine(&api.Config{ServiceName: "test_service", InstanceID: "test_instance"}, redisClient)
			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponseBody)
		})
	}

}
