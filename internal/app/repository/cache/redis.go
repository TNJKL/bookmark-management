package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// redisDB is the Redis implementation of the cache DB interface.
type redisDB struct {
	client *redis.Client
}

// NewRedisDB creates a new Redis cache DB instance using the provided redis.Client.
func NewRedisDB(client *redis.Client) DB {
	return &redisDB{
		client: client,
	}
}

// SetCacheData stores data in a Redis hash under cacheGroupKey with field cacheKey and sets expiration using a transaction pipeline.
func (r *redisDB) SetCacheData(ctx context.Context, cacheGroupKey, cacheKey string, value []byte, exp time.Duration) error {
	if r.client == nil {
		return nil
	}
	_, err := r.client.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HSet(ctx, cacheGroupKey, cacheKey, value)
		pipe.Expire(ctx, cacheGroupKey, exp)
		return nil
	})
	return err
}

// GetCacheData retrieves cached byte data from a Redis hash field under cacheGroupKey.
func (r *redisDB) GetCacheData(ctx context.Context, cacheGroupKey, cacheKey string) ([]byte, error) {
	if r.client == nil {
		return nil, redis.Nil
	}
	return r.client.HGet(ctx, cacheGroupKey, cacheKey).Bytes()
}

// DeleteCache removes a cache key or group key from Redis.
func (r *redisDB) DeleteCache(ctx context.Context, key string) error {
	if r.client == nil {
		return nil
	}
	return r.client.Del(ctx, key).Err()
}
