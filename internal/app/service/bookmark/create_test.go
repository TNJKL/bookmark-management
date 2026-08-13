package bookmark

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	utilsMocks "github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

var errTest = errors.New("test error")

func TestService_CreateBookmark(t *testing.T) {
	t.Parallel()
	inputDescription := "Test Bookmark"
	inputURL := "https://google.com"
	inputUserID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	mockCodeInt := 67
	mockEncodedCode := "k1b8X3k"

	inputBM := &model.Bookmark{
		Description: inputDescription,
		URL:         inputURL,
		UserID:      inputUserID,
		Code:        "pending",
	}

	createdInDB := &model.Bookmark{
		Base:        fixtures.GetTestBase("bookmark-id-123"),
		Description: inputDescription,
		URL:         inputURL,
		UserID:      inputUserID,
		Code:        "pending",
		CodeInt:     mockCodeInt,
	}

	//kiểm tra code được update phải chứa SQL Prefix (i-z) ở đầu và kết thúc bằng mockEncodedCode
	matchCodeWithPrefix := mock.MatchedBy(func(code string) bool {
		return strings.HasSuffix(code, mockEncodedCode) && utils.IsSQLCode(code)
	})

	testCases := []struct {
		name             string
		setupMockBase62  func() *utilsMocks.Base62
		setupMockRepo    func(ctx context.Context) *repoMocks.Repository
		setupDB          func(t *testing.T) *gorm.DB
		expectedBookmark *model.Bookmark
		expectedErr      error
	}{
		{
			name: "happy path - creates bookmark and encodes code_int via transaction",
			setupMockBase62: func() *utilsMocks.Base62 {
				mockB62 := utilsMocks.NewBase62(t)
				mockB62.On("Encode", mockCodeInt).Return(mockEncodedCode)
				return mockB62
			},
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateBookmarkTx", ctx, mock.AnythingOfType("*gorm.DB"), inputBM).Return(createdInDB, nil)
				mockRepo.On("UpdateCodeTx", ctx, mock.AnythingOfType("*gorm.DB"), createdInDB.ID, matchCodeWithPrefix).Return(nil)
				return mockRepo
			},
			setupDB: func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			expectedBookmark: &model.Bookmark{
				Base:        fixtures.GetTestBase("bookmark-id-123"),
				Description: inputDescription,
				URL:         inputURL,
				UserID:      inputUserID,
				Code:        mockEncodedCode,
				CodeInt:     mockCodeInt,
			},
			expectedErr: nil,
		},
		{
			name: "fail - repo CreateBookmarkTx error",
			setupMockBase62: func() *utilsMocks.Base62 {
				return utilsMocks.NewBase62(t)
			},
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateBookmarkTx", ctx, mock.AnythingOfType("*gorm.DB"), inputBM).Return(nil, errTest)
				return mockRepo
			},
			setupDB:          func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			expectedBookmark: nil,
			expectedErr:      errTest,
		},
		{
			name: "fail - repo UpdateCodeTx error",
			setupMockBase62: func() *utilsMocks.Base62 {
				mockB62 := utilsMocks.NewBase62(t)
				mockB62.On("Encode", mockCodeInt).Return(mockEncodedCode)
				return mockB62
			},
			setupMockRepo: func(ctx context.Context) *repoMocks.Repository {
				mockRepo := repoMocks.NewRepository(t)
				mockRepo.On("CreateBookmarkTx", ctx, mock.AnythingOfType("*gorm.DB"), inputBM).Return(createdInDB, nil)
				mockRepo.On("UpdateCodeTx", ctx, mock.AnythingOfType("*gorm.DB"), createdInDB.ID, matchCodeWithPrefix).Return(errTest)
				return mockRepo
			},
			setupDB:          func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			expectedBookmark: nil,
			expectedErr:      errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockB62 := tc.setupMockBase62()
			mockRepo := tc.setupMockRepo(ctx)
			db := tc.setupDB(t)
			svc := NewService(mockRepo, mockB62, db)

			bookmark, err := svc.CreateBookmark(ctx, inputDescription, inputURL, inputUserID)
			assert.ErrorIs(t, err, tc.expectedErr)

			if tc.expectedErr == nil {
				assert.NotNil(t, bookmark)
				assert.Equal(t, inputDescription, bookmark.Description)
				assert.Equal(t, inputURL, bookmark.URL)
				assert.Equal(t, inputUserID, bookmark.UserID)
				assert.True(t, strings.HasSuffix(bookmark.Code, mockEncodedCode))
				assert.True(t, utils.IsSQLCode(bookmark.Code))
			}
		})
	}
}
