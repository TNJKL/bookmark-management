package link

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	bmMocks "github.com/TNJKL/bookmark-management/internal/app/repository/bookmark/mocks"
	repoMocks "github.com/TNJKL/bookmark-management/internal/app/repository/mocks"
	"github.com/TNJKL/bookmark-management/internal/app/repository/urlstorage"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	keyGenMocks "github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errTest = errors.New("test error")

func TestShortenURL_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		setupMockKeyGenerator func(t *testing.T) *keyGenMocks.KeyGenerator
		setupMockStorage      func(t *testing.T, ctx context.Context) *repoMocks.URLStorage
		inputURL              string
		inputExp              int64
		expectedKeySuffix     string
		expectedErr           error
	}{
		{
			name: "happy path",
			setupMockKeyGenerator: func(t *testing.T) *keyGenMocks.KeyGenerator {
				mockKeyGenerator := keyGenMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", linkKeyLength-1).Return("Songok")
				return mockKeyGenerator
			},
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				matchKey := mock.MatchedBy(func(key string) bool {
					return len(key) == 7 && utils.IsRedisCode(key) && strings.HasSuffix(key, "Songok")
				})
				mockStorage.On("GetURL", ctx, matchKey).Return("", urlstorage.ErrorCodeNotFound)
				mockStorage.On("StoreURL", ctx, matchKey, "https://google.com", 3600*time.Second).Return(nil)
				return mockStorage
			},
			inputURL:          "https://google.com",
			inputExp:          3600,
			expectedKeySuffix: "Songok",
			expectedErr:       nil,
		},
		{
			name: "key collision - retries and generates new key successfully",
			setupMockKeyGenerator: func(t *testing.T) *keyGenMocks.KeyGenerator {
				mockKeyGenerator := keyGenMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", linkKeyLength-1).Return("DupKey").Once()
				mockKeyGenerator.On("GenerateKey", linkKeyLength-1).Return("NewKey").Once()
				return mockKeyGenerator
			},
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				matchDupKey := mock.MatchedBy(func(key string) bool {
					return len(key) == 7 && utils.IsRedisCode(key) && strings.HasSuffix(key, "DupKey")
				})
				matchNewKey := mock.MatchedBy(func(key string) bool {
					return len(key) == 7 && utils.IsRedisCode(key) && strings.HasSuffix(key, "NewKey")
				})
				mockStorage.On("GetURL", ctx, matchDupKey).Return("https://existing.com", nil)
				mockStorage.On("GetURL", ctx, matchNewKey).Return("", urlstorage.ErrorCodeNotFound)
				mockStorage.On("StoreURL", ctx, matchNewKey, "https://google.com", 3600*time.Second).Return(nil)
				return mockStorage
			},
			inputURL:          "https://google.com",
			inputExp:          3600,
			expectedKeySuffix: "NewKey",
			expectedErr:       nil,
		},
		{
			name: "get URL error",
			setupMockKeyGenerator: func(t *testing.T) *keyGenMocks.KeyGenerator {
				mockKeyGenerator := keyGenMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", linkKeyLength-1).Return("Songok")
				return mockKeyGenerator
			},
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				matchKey := mock.MatchedBy(func(key string) bool {
					return len(key) == 7 && utils.IsRedisCode(key) && strings.HasSuffix(key, "Songok")
				})
				mockStorage.On("GetURL", ctx, matchKey).Return("", errTest)
				return mockStorage
			},
			inputURL:          "https://google.com",
			inputExp:          3600,
			expectedKeySuffix: "",
			expectedErr:       errTest,
		},
		{
			name: "store URL error",
			setupMockKeyGenerator: func(t *testing.T) *keyGenMocks.KeyGenerator {
				mockKeyGenerator := keyGenMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", linkKeyLength-1).Return("Songok")
				return mockKeyGenerator
			},
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				matchKey := mock.MatchedBy(func(key string) bool {
					return len(key) == 7 && utils.IsRedisCode(key) && strings.HasSuffix(key, "Songok")
				})
				mockStorage.On("GetURL", ctx, matchKey).Return("", urlstorage.ErrorCodeNotFound)
				mockStorage.On("StoreURL", ctx, matchKey, "https://google.com", 3600*time.Second).Return(errTest)
				return mockStorage
			},
			inputURL:          "https://google.com",
			inputExp:          3600,
			expectedKeySuffix: "",
			expectedErr:       errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockKeyGenerator := tc.setupMockKeyGenerator(t)
			mockStorage := tc.setupMockStorage(t, ctx)
			mockBookmarkRepo := bmMocks.NewRepository(t)

			shortenURLsvc := NewShortenUrl(mockStorage, mockKeyGenerator, mockBookmarkRepo)

			code, err := shortenURLsvc.CreateShortenLink(ctx, tc.inputURL, tc.inputExp)
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.expectedErr == nil {
				assert.Len(t, code, 7)
				assert.True(t, utils.IsRedisCode(code))
				assert.True(t, strings.HasSuffix(code, tc.expectedKeySuffix))
			} else {
				assert.Empty(t, code)
			}
		})
	}
}

func TestShortenURL_GetLinkFromCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		setupMockStorage      func(t *testing.T, ctx context.Context) *repoMocks.URLStorage
		setupMockBookmarkRepo func(t *testing.T, ctx context.Context) *bmMocks.Repository
		inputCode             string
		expectedURL           string
		expectedErr           error
	}{
		{
			name: "happy path - redis code",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "aSongok").Return("https://google.com", nil)
				return mockStorage
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				return bmMocks.NewRepository(t)
			},
			inputCode:   "aSongok",
			expectedURL: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "redis code not found",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "aInvalid").Return("", urlstorage.ErrorCodeNotFound)
				return mockStorage
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				return bmMocks.NewRepository(t)
			},
			inputCode:   "aInvalid",
			expectedURL: "",
			expectedErr: urlstorage.ErrorCodeNotFound,
		},
		{
			name: "redis storage error",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "aSongok").Return("", errTest)
				return mockStorage
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				return bmMocks.NewRepository(t)
			},
			inputCode:   "aSongok",
			expectedURL: "",
			expectedErr: errTest,
		},
		{
			name: "happy path - sql code",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				return repoMocks.NewURLStorage(t)
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				mockRepo := bmMocks.NewRepository(t)
				mockRepo.On("GetByCode", ctx, "i123456").Return(&model.Bookmark{
					URL: "https://sql-google.com",
				}, nil)
				return mockRepo
			},
			inputCode:   "i123456",
			expectedURL: "https://sql-google.com",
			expectedErr: nil,
		},
		{
			name: "sql code record not found",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				return repoMocks.NewURLStorage(t)
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				mockRepo := bmMocks.NewRepository(t)
				mockRepo.On("GetByCode", ctx, "i999999").Return(nil, dbutils.ErrRecordNotFound)
				return mockRepo
			},
			inputCode:   "i999999",
			expectedURL: "",
			expectedErr: urlstorage.ErrorCodeNotFound,
		},
		{
			name: "sql repo error",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				return repoMocks.NewURLStorage(t)
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				mockRepo := bmMocks.NewRepository(t)
				mockRepo.On("GetByCode", ctx, "i123456").Return(nil, errTest)
				return mockRepo
			},
			inputCode:   "i123456",
			expectedURL: "",
			expectedErr: errTest,
		},
		{
			name: "invalid code prefix",
			setupMockStorage: func(t *testing.T, ctx context.Context) *repoMocks.URLStorage {
				return repoMocks.NewURLStorage(t)
			},
			setupMockBookmarkRepo: func(t *testing.T, ctx context.Context) *bmMocks.Repository {
				return bmMocks.NewRepository(t)
			},
			inputCode:   "1234567",
			expectedURL: "",
			expectedErr: urlstorage.ErrorCodeNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()

			mockStorage := tc.setupMockStorage(t, ctx)
			mockBookmarkRepo := tc.setupMockBookmarkRepo(t, ctx)

			shortenURLsvc := NewShortenUrl(mockStorage, nil, mockBookmarkRepo)

			url, err := shortenURLsvc.GetLinkFromCode(ctx, tc.inputCode)
			assert.Equal(t, tc.expectedURL, url)
			assert.ErrorIs(t, err, tc.expectedErr)
		})
	}
}
