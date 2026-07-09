package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// URLStorage defines storage operations for shortened URLs.
//
//go:generate mockery --name URLStorage --filename urlstorage.go
type URLStorage interface {
	StoreURL(ctx context.Context, code, url string, exp time.Duration) error
	GetURL(ctx context.Context, code string) (string, error)
}

// urlStorage is the Redis-backed implementation of URLStorage.
type urlStorage struct {
	rclient *redis.Client
}

// NewURLStorage creates a new urlStorage using the given Redis client.
func NewURLStorage(rclient *redis.Client) URLStorage {
	return &urlStorage{
		rclient: rclient,
	}
}

// StoreURL sets code as a key mapping to url, expiring after exp.
func (u *urlStorage) StoreURL(ctx context.Context, code, url string, exp time.Duration) error {
	return u.rclient.Set(ctx, code, url, exp).Err()
}

// GetURL fetches the url stored under code, or an error if not found.
func (u *urlStorage) GetURL(ctx context.Context, code string) (string, error) {
	return u.rclient.Get(ctx, code).Result()
}
