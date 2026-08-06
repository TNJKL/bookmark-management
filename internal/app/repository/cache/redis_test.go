package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/repository/cache"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisDB_SetCacheData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMockRedis func(ctx context.Context) *redis.Client
		cacheGroupKey  string
		cacheKey       string
		value          []byte
		exp            time.Duration
		expectedErr    error
		verifyFunc     func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "happy path",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			cacheKey:      "1_10",
			value:         []byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y97","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"},{"id":"806fd7e2-0f30-439d-bb5b-071d8377fe9d","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y98","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`),
			exp:           24 * time.Hour,
			expectedErr:   nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				val, err := r.HGet(ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6", "1_10").Bytes()
				assert.NoError(t, err)
				assert.Equal(t, []byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y97","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"},{"id":"806fd7e2-0f30-439d-bb5b-071d8377fe9d","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y98","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`), val)
			},
		},
		{
			name: "set empty value",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			cacheKey:      "2_20",
			value:         []byte(""),
			exp:           24 * time.Hour,
			expectedErr:   nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				val, err := r.HGet(ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6", "2_20").Bytes()
				assert.NoError(t, err)
				assert.Equal(t, []byte(""), val)
			},
		},

		{
			name: "connection error",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				_ = mockR.Close()
				return mockR
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			cacheKey:      "3_30",
			value:         []byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y97","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"},{"id":"806fd7e2-0f30-439d-bb5b-071d8377fe9d","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y98","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`),
			exp:           24 * time.Hour,
			expectedErr:   redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRedis := tc.setupMockRedis(ctx)
			repo := cache.NewRedisDB(mockRedis)
			err := repo.SetCacheData(ctx, tc.cacheGroupKey, tc.cacheKey, tc.value, tc.exp)
			assert.Equal(t, tc.expectedErr, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mockRedis)
			}
		})
	}

}

func TestRedisDB_GetCacheData(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		setupMockRedis func(ctx context.Context) *redis.Client
		cacheGroupKey  string
		cacheKey       string
		expectedValue  []byte
		expectedErr    error
	}{
		{
			name: "happy path - cache hit",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				err := mockR.HSet(ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
					"1_10",
					[]byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169",
					"created_at":"2018-02-18T01:02:03.000000004Z",
					"updated_at":"2018-02-18T01:02:03.000000004Z",
					"description":"A bookmark description",
					"url":"http://google.com","code":"wfaA1y97",
					"user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`)).Err()
				assert.NoError(t, err)
				return mockR
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			cacheKey:      "1_10",
			expectedValue: []byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169",
					"created_at":"2018-02-18T01:02:03.000000004Z",
					"updated_at":"2018-02-18T01:02:03.000000004Z",
					"description":"A bookmark description",
					"url":"http://google.com","code":"wfaA1y97",
					"user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`),
			expectedErr: nil,
		},
		{
			name: "cache miss",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			cacheKey:      "non-existent-key",
			expectedValue: nil,
			expectedErr:   redis.Nil,
		},
		{
			name: "connection error",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				_ = mockR.Close()
				return mockR
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			cacheKey:      "3_30",
			expectedValue: nil,
			expectedErr:   redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRedis := tc.setupMockRedis(ctx)
			repo := cache.NewRedisDB(mockRedis)
			val, err := repo.GetCacheData(ctx, tc.cacheGroupKey, tc.cacheKey)
			assert.Equal(t, tc.expectedValue, val)
			assert.Equal(t, tc.expectedErr, err)

		})

	}

}

func TestRedisDB_DeleteCacheData(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		setupMockRedis func(ctx context.Context) *redis.Client
		cacheGroupKey  string
		expectedErr    error
		verifyFunc     func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "happy path - delete existing group key",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				err := mockR.HSet(ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
					"1_10",
					[]byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169",
					"created_at":"2018-02-18T01:02:03.000000004Z",
					"updated_at":"2018-02-18T01:02:03.000000004Z",
					"description":"A bookmark description",
					"url":"http://google.com","code":"wfaA1y97",
					"user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`)).Err()
				assert.NoError(t, err)
				return mockR
			},
			cacheGroupKey: "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
			expectedErr:   nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				exists, err := r.Exists(ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), exists)
			},
		},
		{
			name: "delete non-existent group key",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			cacheGroupKey: "non_existent_group",
			expectedErr:   nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				exists, err := r.Exists(ctx, "non_existent_group").Result()
				assert.NoError(t, err)
				assert.Equal(t, int64(0), exists)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRedis := tc.setupMockRedis(ctx)
			repo := cache.NewRedisDB(mockRedis)

			err := repo.DeleteCache(ctx, tc.cacheGroupKey)
			assert.Equal(t, tc.expectedErr, err)
			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mockRedis)
			}
		})
	}
}
