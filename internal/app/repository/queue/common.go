package queue

import "context"

//go:generate mockery --name Repository --filename queueRepo.go
type Repository interface {
	PushMessage(ctx context.Context, message []byte) error
}
