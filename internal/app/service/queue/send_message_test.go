package queue

import (
	"context"
	"errors"
	"fmt"
	"testing"

	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/queue/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errTest = errors.New("test error")

func TestService_SendImportBookmarkJob(t *testing.T) {
	t.Parallel()

	userID := "user-uuid-123"

	// Helper tạo dữ liệu Bookmark mẫu
	makeBookmarks := func(count int) []*ImportBookmarkInput {
		list := make([]*ImportBookmarkInput, count)
		for i := 0; i < count; i++ {
			list[i] = &ImportBookmarkInput{
				Description: fmt.Sprintf("Bookmark %d", i+1),
				URL:         fmt.Sprintf("https://example%d.com", i+1),
			}
		}
		return list
	}

	testCases := []struct {
		name           string
		bookmarkInputs []*ImportBookmarkInput
		setupMockRepo  func(ctx context.Context) *repoMocks.Repository
		expectedErr    error
	}{
		{
			name:           "happy path - single batch (<= 20 items)",
			bookmarkInputs: makeBookmarks(2),
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				// Đẩy 2 items -> Đủ 1 batch -> PushMessage được gọi đúng 1 lần (.Once())
				mockRepo.On("PushMessage", ctx, mock.Anything).Return(nil).Once()
				return mockRepo
			},
			expectedErr: nil,
		},
		{
			name:           "happy path - multiple batches (> 20 items)",
			bookmarkInputs: makeBookmarks(25),
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				// Đẩy 25 items với BatchSize = 20 -> Tách làm 2 batches -> PushMessage được gọi 2 lần (.Times(2))
				mockRepo.On("PushMessage", ctx, mock.Anything).Return(nil).Times(2)
				return mockRepo
			},
			expectedErr: nil,
		},
		{
			name:           "error case - repo PushMessage fails",
			bookmarkInputs: makeBookmarks(2),
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				// Giả lập Repo PushMessage bị lỗi
				mockRepo.On("PushMessage", ctx, mock.Anything).Return(errTest).Once()
				return mockRepo
			},
			expectedErr: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)
			svc := NewService(mockRepo)

			err := svc.SendImportBookmarkJob(ctx, userID, tc.bookmarkInputs)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
