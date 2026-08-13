// Package middleware provides HTTP middleware handlers for request authentication, rate limiting, and context processing.
package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/repository/ratelimit"
	"github.com/TNJKL/bookmark-management/pkg/requestutils"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

// RateLimit defines the contract for rate limiting middleware.
type RateLimit interface {
	// RateLimit returns a Gin middleware HandlerFunc that enforces rate limits on HTTP endpoints.
	RateLimit() gin.HandlerFunc
}

// rateLimitMiddleWare implements the RateLimit interface using a Redis-backed rate limit repository.
type rateLimitMiddleWare struct {
	repo ratelimit.Repo
}

// NewRateLimit creates a new instance of RateLimit middleware configured with the provided rate limit repository.
func NewRateLimit(repo ratelimit.Repo) RateLimit {
	return &rateLimitMiddleWare{
		repo: repo,
	}
}

const (
	// rateLimitInterval defines the time window duration for tracking request rates (10 seconds).
	rateLimitInterval = 10 * time.Second

	// rateLimitCount defines the maximum number of allowed requests per user within rateLimitInterval (20 requests).
	rateLimitCount = 20

	// rateLimitKeyFormat defines the Redis cache key pattern used to store rate limit counters per user.
	rateLimitKeyFormat = "rate_limit:%s"
)

// RateLimit returns a Gin middleware handler function that tracks and enforces per-user request rate limits.
// It extracts the user ID from the request context, checks the current request count against Redis,
// aborts with HTTP 429 (Too Many Requests) if the limit is exceeded, or increments the counter and proceeds.
func (r *rateLimitMiddleWare) RateLimit() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get userID from request
		uid, err := requestutils.GetUserIDFromRequest(ctx)
		if err != nil {
			return
		}

		//create rate limit key
		rateLimitKey := fmt.Sprintf(rateLimitKeyFormat, uid)

		//get current rate limit
		currentRate, err := r.repo.GetCurrentRateLimit(ctx, rateLimitKey)
		if err != nil {
			log.Error().Err(err).Msg("failed to get current rate limit")
		}

		//check if rate limit exceeded
		if currentRate >= rateLimitCount {
			ctx.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "rate limit exceeded",
			})
			return
		}

		//increase rate limit
		r.repo.IncreaseRateLimit(ctx, rateLimitKey, rateLimitInterval)

		ctx.Next()
	}
}
