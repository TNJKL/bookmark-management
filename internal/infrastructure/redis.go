package infrastructure

import (
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"github.com/redis/go-redis/v9"
)

// CreateRedisConn creates a new redis connection
func CreateRedisConn() *redis.Client {
	// Create redis db connection
	redisClient, err := redisPkg.NewClient("")
	utils.NoErr(err)

	return redisClient
}
