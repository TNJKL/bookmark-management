package bookmark

import (
	"context"
	"errors"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

func TestService_CreateBookmark(t *testing.T) {
	t.Parallel()
	inputDescription := "Test Bookmark"
	inputURL := "https://google.com"
	inputUserID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	mockCode := "12345678"

	expectedBookmark := &model.Bookmark{
		Description: inputDescription,
		URL:         inputURL,
		UserID:      inputUserID,
		Code:        mockCode,
	}

	testCases := []struct {
		name             string
		setupMockKeyGen  func() *mocks.KeyGenerator
		setupMockRepo    func(ctx context.Context, input *model.Bookmark) *repoMocks.Repository
		expectedBookmark *model.Bookmark
		expectedErr      error
	}{
		{
			name: "happy path",
			setupMockKeyGen: func() *mocks.KeyGenerator {
				mockKeyGen := mocks.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", codeLength).Return(mockCode)
				return mockKeyGen
			},
			setupMockRepo: func(ctx context.Context, input *model.Bookmark) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateBookmark", ctx, input).Return(expectedBookmark, nil)
				return mockRepo
			},
			expectedBookmark: expectedBookmark,
			expectedErr:      nil,
		},
		{
			name: "repo create bookmark fails - general db error",
			setupMockKeyGen: func() *mocks.KeyGenerator {
				mockKeyGen := mocks.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", codeLength).Return(mockCode)
				return mockKeyGen
			},
			setupMockRepo: func(ctx context.Context, input *model.Bookmark) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateBookmark", ctx, input).Return(nil, errTest)
				return mockRepo
			},
			expectedBookmark: nil,
			expectedErr:      errTest,
		},
		{
			name: "repo returns unique constraint error (duplicate code)",
			setupMockKeyGen: func() *mocks.KeyGenerator {
				mockKeyGen := mocks.NewKeyGenerator(t)
				mockKeyGen.On("GenerateKey", codeLength).Return(mockCode)
				return mockKeyGen
			},
			setupMockRepo: func(ctx context.Context, input *model.Bookmark) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateBookmark", ctx, input).Return(nil, dbutils.ErrUniqueConstraint)
				return mockRepo
			},
			expectedBookmark: nil,
			expectedErr:      dbutils.ErrUniqueConstraint,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockKeyGen := tc.setupMockKeyGen()

			inputBookmark := &model.Bookmark{
				Description: inputDescription,
				URL:         inputURL,
				UserID:      inputUserID,
				Code:        mockCode,
			}
			mockRepo := tc.setupMockRepo(ctx, inputBookmark)

			svc := NewService(mockRepo, mockKeyGen)

			bookmark, err := svc.CreateBookmark(ctx, inputDescription, inputURL, inputUserID)
			assert.Equal(t, tc.expectedBookmark, bookmark)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
