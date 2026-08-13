package ratelimit

import (
	"context"
	"testing"
	"time"

	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestRedisRepo_IncreaseRateLimit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		setupRepo  func() Repo
		verifyFunc func(ctx context.Context, client *redis.Client)
	}{
		{
			name: "happy path - increase rate limit and set TT",
			setupRepo: func() Repo {
				mockRedis := redisPkg.InitMockRedis(t)
				return NewRedisRepo(mockRedis)
			},
			verifyFunc: func(ctx context.Context, client *redis.Client) {
				val, err := client.Get(ctx, "rate_limit:user123").Int()
				assert.NoError(t, err)
				assert.Equal(t, 1, val)

				ttl, err := client.TTL(ctx, "rate_limit:user123").Result()
				assert.NoError(t, err)
				assert.True(t, ttl > 0)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			var client *redis.Client
			repo := tc.setupRepo()
			if rRepo, ok := repo.(*redisRepo); ok {
				client = rRepo.client
			}
			repo.IncreaseRateLimit(ctx, "rate_limit:user123", 1*time.Minute)
			tc.verifyFunc(ctx, client)
		})
	}
}

func TestRedisRepo_GetCurrentRateLimit(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupRepo   func(ctx context.Context) Repo
		expectedVal int
		expectedErr error
	}{
		{
			name: "happy path - key exists",
			setupRepo: func(ctx context.Context) Repo {
				mockRedis := redisPkg.InitMockRedis(t)
				_ = mockRedis.Set(ctx, "rate_limit:user123", 5, 1*time.Minute).Err()
				return NewRedisRepo(mockRedis)
			},
			expectedVal: 5,
			expectedErr: nil,
		},

		{
			name: "key not found",
			setupRepo: func(ctx context.Context) Repo {
				mockRedis := redisPkg.InitMockRedis(t)
				return NewRedisRepo(mockRedis)
			},
			expectedVal: 0,
			expectedErr: redis.Nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			repo := tc.setupRepo(ctx)
			val, err := repo.GetCurrentRateLimit(ctx, "rate_limit:user123")

			assert.Equal(t, tc.expectedVal, val)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
