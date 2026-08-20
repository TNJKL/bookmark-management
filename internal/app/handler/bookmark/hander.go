package bookmark

import (
	"github.com/TNJKL/bookmark-management/internal/app/service/bookmark"
	"github.com/TNJKL/bookmark-management/internal/app/service/queue"
	"github.com/gin-gonic/gin"
)

// Handler defines HTTP handler methods for bookmark operations
type Handler interface {
	CreateBookmark(ctx *gin.Context)
	GetBookmarks(ctx *gin.Context)
	UpdateBookmark(ctx *gin.Context)
	DeleteBookmark(ctx *gin.Context)
	ImportBookmarks(ctx *gin.Context)
}

// bookmarkHandler implements the Handler interface
type bookmarkHandler struct {
	svc   bookmark.Service
	queue queue.Service
}

// NewHandler creates a new instance of bookmarkHandler with the provided service
func NewHandler(svc bookmark.Service, queue queue.Service) Handler {
	return &bookmarkHandler{
		svc:   svc,
		queue: queue,
	}
}
