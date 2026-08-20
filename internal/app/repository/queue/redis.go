package queue

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisQueue struct {
	client    *redis.Client
	queueName string
}

func NewRedisQueueRepo(client *redis.Client, queueName string) Repository {
	return &redisQueue{
		client:    client,
		queueName: queueName,
	}
}

func (r *redisQueue) PushMessage(ctx context.Context, message []byte) error {
	return r.client.LPush(ctx, r.queueName, message).Err()
}
