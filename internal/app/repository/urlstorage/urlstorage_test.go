package urlstorage

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

		setupMock func(ctx context.Context) *redis.Client

		expectedErr error

		verifyFunc func(ctx context.Context, r *redis.Client)
	}{
		{
			name: "normal case",

			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedErr: nil,
			verifyFunc: func(ctx context.Context, r *redis.Client) {
				res, err := r.Get(ctx, "test").Result()
				// Đọc lại key "test" từ Redis giả
				// Tương đương lệnh Redis: GET "test"
				// res = giá trị lấy được (string)
				// err = lỗi nếu có (ví dụ key không tồn tại)
				assert.NoError(t, err)

				//Xác nhận: giá trị đọc ra phải là "https://google.com"
				//(đúng với URL mà StoreURL đã lưu)
				assert.Equal(t, "https://google.com", res)

				//Nói cách khác: sau khi  StoreURL  chạy xong,
				//verifyFunc  đóng vai "thanh tra" —
				//mở kho Redis ra kiểm tra hàng có đúng không.
			},
		},

		{
			name: "connection error",

			setupMock: func(ctx context.Context) *redis.Client {
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

			mock := tc.setupMock(ctx)
			storage := NewURLStorage(mock)
			err := storage.StoreURL(ctx, "test", "https://google.com", time.Hour)
			assert.Equal(t, tc.expectedErr, err)

			if tc.verifyFunc != nil {
				tc.verifyFunc(ctx, mock)
			}

		})
	}

}

func TestUrlStorage_GetURL(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupMock   func(ctx context.Context) *redis.Client
		expectedVal string
		expectedErr error
	}{
		{
			name: "normal case",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)

				err := mock.Set(ctx, "Songoku", "https://google.com", 500000).Err()
				assert.NoError(t, err)
				return mock
			},
			expectedVal: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "code not found",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				return mock
			},
			expectedVal: "",
			expectedErr: ErrorCodeNotFound,
		},
		{
			name: "connection error",
			setupMock: func(ctx context.Context) *redis.Client {
				mock := redisPkg.InitMockRedis(t)
				_ = mock.Close()
				return mock
			},
			expectedVal: "",
			expectedErr: redis.ErrClosed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mock := tc.setupMock(ctx)
			storage := NewURLStorage(mock)
			val, err := storage.GetURL(ctx, "Songoku")

			assert.Equal(t, tc.expectedVal, val)
			assert.Equal(t, tc.expectedErr, err)
		})

	}
}
