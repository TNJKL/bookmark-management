package repository

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// HealthRepository defines the operations for checking database and cache health.
//
//go:generate mockery --name HealthRepository --filename healthcheck.go
type HealthRepository interface {
	// Ping checks the connection health of the underlying storage system.
	Ping(ctx context.Context) error
}

// healthRepository is the Redis-backed implementation of HealthRepository.
type healthRepository struct {
	rclient *redis.Client
}

// NewHealthRepository creates a new HealthRepository using the given Redis client.
func NewHealthRepository(rclient *redis.Client) HealthRepository {
	return &healthRepository{
		rclient: rclient,
	}
}

// Ping checks the connection health of the underlying storage system (e.g. Redis).
func (h *healthRepository) Ping(ctx context.Context) error {
	return h.rclient.Ping(ctx).Err()
}
