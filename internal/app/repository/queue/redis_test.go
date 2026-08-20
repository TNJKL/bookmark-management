package queue_test

import (
	"context"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/repository/queue"
	redisPkg "github.com/TNJKL/bookmark-management/pkg/redis"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRedisQueue_PushMessage(t *testing.T) {
	t.Parallel()

	const testQueueName = "test_bookmark_queue"

	testCases := []struct {
		name           string
		setupMockRedis func(ctx context.Context) *redis.Client
		message        []byte
		expectedErr    error
		verifyFunc     func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "happy path",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			message:     []byte(`{"uid":"123","url":"https://google.com"}`),
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				poppedMsg, err := r.RPop(ctx, testQueueName).Bytes()
				assert.NoError(t, err)
				assert.Equal(t, []byte(`{"uid":"123","url":"https://google.com"}`), poppedMsg)
			},
		},
		{
			name: "set empty value",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				return redisPkg.InitMockRedis(t)
			},
			message:     []byte(""),
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				poppedMsg, err := r.RPop(ctx, testQueueName).Bytes()
				require.NoError(t, err)
				assert.Equal(t, []byte(""), poppedMsg)
			},
		},
		{
			name: "connection error",
			setupMockRedis: func(ctx context.Context) *redis.Client {
				mockR := redisPkg.InitMockRedis(t)
				_ = mockR.Close()
				return mockR
			},
			message:     []byte(`{"uid":"123"}`),
			expectedErr: redis.ErrClosed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mockRedis := tc.setupMockRedis(ctx)
			repo := queue.NewRedisQueueRepo(mockRedis, testQueueName)

			err := repo.PushMessage(ctx, tc.message)
			assert.Equal(t, tc.expectedErr, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mockRedis)
			}
		})
	}
}
