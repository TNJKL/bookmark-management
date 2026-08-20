package queue

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/repository/queue"
)

//go:generate mockery --name Service --filename service.go
type Service interface {
	SendImportBookmarkJob(ctx context.Context, uid string, bookmarkInputs []*ImportBookmarkInput) error
}

type service struct {
	q queue.Repository
}

func NewService(q queue.Repository) Service {
	return &service{
		q: q,
	}
}
