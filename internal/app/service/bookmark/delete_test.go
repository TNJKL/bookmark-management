package bookmark

import (
	"context"
	"testing"

	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	keyGenMocks "github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

func TestService_DeleteBookmark(t *testing.T) {
	t.Parallel()
	id := "bookmark-123"
	userID := "user-123"

	testCases := []struct {
		name          string
		setupMockRepo func(ctx context.Context) *repoMocks.Repository
		expectedErr   error
	}{
		{
			name: "happy path",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("DeleteBookmark", ctx, id, userID).Return(nil)
				return mockRepo
			},
			expectedErr: nil,
		},
		{
			name: "record not found",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("DeleteBookmark", ctx, id, userID).Return(dbutils.ErrRecordNotFound)
				return mockRepo
			},
			expectedErr: dbutils.ErrRecordNotFound,
		},
		{
			name: "internal db error",
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("DeleteBookmark", ctx, id, userID).Return(errTest)
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
			mockKeyGen := keyGenMocks.NewKeyGenerator(t)

			svc := NewService(mockRepo, mockKeyGen)
			err := svc.DeleteBookmark(ctx, id, userID)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
