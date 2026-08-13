package bookmark

import (
	"context"
	"testing"

	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
)

func TestService_UpdateBookmark(t *testing.T) {
	t.Parallel()
	id := "bookmark-123"
	userID := "user-123"
	description := "Updated Description"
	url := "https://updated-url.com"

	testCases := []struct {
		name          string
		setupMockRepo func(ctx context.Context) *repoMocks.Repository
		expectedErr   error
	}{
		{
			name: "happy path",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("UpdateBookmark", ctx, id, userID, description, url).Return(nil)
				return mockRepo
			},
			expectedErr: nil,
		},
		{
			name: "record not found",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("UpdateBookmark", ctx, id, userID, description, url).Return(dbutils.ErrRecordNotFound)
				return mockRepo
			},
			expectedErr: dbutils.ErrRecordNotFound,
		},
		{
			name: "internal db error",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("UpdateBookmark", ctx, id, userID, description, url).Return(errTest)
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

			svc := NewService(mockRepo, nil, nil)
			err := svc.UpdateBookmark(ctx, id, userID, description, url)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
