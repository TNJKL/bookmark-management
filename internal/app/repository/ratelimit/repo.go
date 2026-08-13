// Package ratelimit provides the data repository layer for rate limiting operations backed by Redis.
package ratelimit

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

// Repo defines the data access contract for managing rate limit counters in Redis.
//
//go:generate mockery --name Repo --filename repo.go
type Repo interface {
	// IncreaseRateLimit increments the request counter for the specified key and sets an expiration duration.
	IncreaseRateLimit(ctx context.Context, key string, exp time.Duration)

	// GetCurrentRateLimit retrieves the current request count for the specified key.
	GetCurrentRateLimit(ctx context.Context, key string) (int, error)
}

// redisRepo implements the Repo interface using a Redis client connection.
type redisRepo struct {
	client *redis.Client
}

// NewRedisRepo creates a new instance of Repo configured with the provided Redis client.
func NewRedisRepo(client *redis.Client) Repo {
	return &redisRepo{client: client}
}

// IncreaseRateLimit atomically increments the request count for the given key in Redis
// and sets its TTL expiration duration using a Redis transaction pipeline.
func (r *redisRepo) IncreaseRateLimit(ctx context.Context, key string, exp time.Duration) {
	if r == nil || r.client == nil {
		return
	}
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.Incr(ctx, key)
		pipe.Expire(ctx, key, exp)
		return nil
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to increase rate limit")
	}
}

// GetCurrentRateLimit fetches the current request counter integer value stored under the specified key in Redis.
// Returns 0 with no error if the repository or Redis client is uninitialized.
func (r *redisRepo) GetCurrentRateLimit(ctx context.Context, key string) (int, error) {
	if r == nil || r.client == nil {
		return 0, nil
	}
	return r.client.Get(ctx, key).Int()
}
