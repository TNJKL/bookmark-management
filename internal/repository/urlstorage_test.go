package repository

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"

	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
)

// Setup mock nhận context vì hàm StoreURL(ctx context.Context) nhận context
// verifyFunc : cái này dùng để tạo 1 function xem data của mình có chuẩn hay ko
func TestUrlStorage_StoreURL(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupMock func(ctx context.Context, t *testing.T) *redis.Client

		expectedErr error

		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "test").Result()
				assert.NoError(t, err)
				assert.Equal(t, "https://google.com", res)
			},
		},

		{
			name: "connection error",

			setupMock: func(ctx context.Context, t *testing.T) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mock := tc.setupMock(ctx, t)
			storage := NewURLStorage(mock)
			err := storage.StoreURL(ctx, "test", "https://google.com", time.Hour)
			assert.Equal(t, tc.expectedErr, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}

		})
	}

}
