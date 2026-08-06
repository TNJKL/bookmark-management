package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/stretchr/testify/assert"
)

var getErrTest = errors.New("test error")

func TestService_GetBookmarks(t *testing.T) {
	t.Parallel()
	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"

	mockBookmarks := []*model.Bookmark{
		{
			Base: model.Base{
				ID: "bookmark-1",
			},
			Description: "Google",
			URL:         "https://google.com",
			Code:        "12345678",
			UserID:      userID,
		},
		{
			Base: model.Base{
				ID: "bookmark-2",
			},
			Description: "Github",
			URL:         "https://github.com",
			Code:        "87654321",
			UserID:      userID,
		},
	}

	testCases := []struct {
		name           string
		inputPage      int
		inputLimit     int
		setupMockRepo  func(ctx context.Context) *repoMocks.Repository
		expectedResult *GetBookmarksResult
		expectedErr    error
	}{
		{
			name:       "happy path",
			inputPage:  1,
			inputLimit: 10,
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetBookmarks", ctx, userID, 10, 0).Return(mockBookmarks, int64(2), nil)
				return mockRepo
			},
			expectedResult: &GetBookmarksResult{
				Bookmarks: mockBookmarks,
				Total:     2,
			},
			expectedErr: nil,
		},
		{
			name:       "happy path - pagination page 2",
			inputPage:  2,
			inputLimit: 2,
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetBookmarks", ctx, userID, 2, 2).Return(mockBookmarks, int64(5), nil)
				return mockRepo
			},
			expectedResult: &GetBookmarksResult{
				Bookmarks: mockBookmarks,
				Total:     5,
			},
			expectedErr: nil,
		},
		{
			name:       "repo get bookmarks fails",
			inputPage:  1,
			inputLimit: 10,
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("GetBookmarks", ctx, userID, 10, 0).Return(nil, int64(0), getErrTest)
				return mockRepo
			},
			expectedResult: nil,
			expectedErr:    getErrTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockRepo := tc.setupMockRepo(ctx)

			svc := NewService(mockRepo, nil, nil)

			res, err := svc.GetBookmarks(ctx, userID, tc.inputPage, tc.inputLimit)

			assert.Equal(t, tc.expectedResult, res)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
