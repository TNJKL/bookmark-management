package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/repository"
	repoMocks "github.com/TNJKL/bookmark-management/internal/repository/mocks"
	ultilMocks "github.com/TNJKL/bookmark-management/pkg/utils/mocks"
	"github.com/stretchr/testify/assert"
)

var errTest = errors.New("test error")

func TestShortenURL_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                  string
		setupMockKeyGenerator func() *ultilMocks.KeyGenerator
		setupMockStorage      func(ctx context.Context) *repoMocks.URLStorage
		inputURL              string
		inputExp              int64
		expectedCode          string
		expectedErr           error
	}{
		{
			name: "get URL error",
			setupMockKeyGenerator: func() *ultilMocks.KeyGenerator {
				mockKeyGenerator := ultilMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", 7).Return("Songoku")
				return mockKeyGenerator
			},
			setupMockStorage: func(ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "Songoku").Return("", errTest)
				return mockStorage
			},
			inputURL:     "https://google.com",
			inputExp:     3600,
			expectedCode: "",
			expectedErr:  errTest,
		},
		{
			name: "store Url error",
			setupMockKeyGenerator: func() *ultilMocks.KeyGenerator {
				mockKeyGenerator := ultilMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", 7).Return("Songoku")
				return mockKeyGenerator
			},
			setupMockStorage: func(ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "Songoku").Return("", repository.ErrorCodeNotFound)
				mockStorage.On("StoreURL", ctx, "Songoku", "https://google.com", time.Hour).Return(errTest)
				return mockStorage
			},
			inputURL:     "https://google.com",
			inputExp:     3600,
			expectedCode: "",
			expectedErr:  errTest,
		},
		{
			name: "happy path",
			setupMockKeyGenerator: func() *ultilMocks.KeyGenerator {
				mockKeyGenerator := ultilMocks.NewKeyGenerator(t)
				mockKeyGenerator.On("GenerateKey", 7).Return("Songoku")
				return mockKeyGenerator
			},

			setupMockStorage: func(ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "Songoku").Return("", repository.ErrorCodeNotFound)
				mockStorage.On("StoreURL", ctx, "Songoku", "https://google.com", time.Hour).Return(nil)
				return mockStorage
			},
			inputURL:     "https://google.com",
			inputExp:     3600,
			expectedCode: "Songoku",
			expectedErr:  nil,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			//tao context cho test (tu huy khi ket thuc)
			ctx := context.Background()

			//goi 2 ham setup Mock
			mockKeyGenerator := tc.setupMockKeyGenerator()
			mockStorage := tc.setupMockStorage(ctx)

			//tao ShortenURL service , truyen mock vao
			shortenURLsvc := NewShortenUrl(mockStorage, mockKeyGenerator)

			//goi ham can test
			code, err := shortenURLsvc.CreateShortenLink(ctx, tc.inputURL, tc.inputExp)
			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedErr, err)

		})
	}

}

func TestShortenURL_GetShortenLink(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupMockStorage func(ctx context.Context) *repoMocks.URLStorage
		inputCode        string
		expectedUrl      string
		expectedErr      error
	}{
		{
			name: "happy path",
			setupMockStorage: func(ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "Songku").Return("https://google.com", nil)
				return mockStorage
			},
			inputCode:   "Songku",
			expectedUrl: "https://google.com",
			expectedErr: nil,
		},
		{
			name: "code not found",
			setupMockStorage: func(ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "invalid-code").Return("", repository.ErrorCodeNotFound)
				return mockStorage
			},
			inputCode:   "invalid-code",
			expectedUrl: "",
			expectedErr: repository.ErrorCodeNotFound,
		},
		{
			name: "storage error",
			setupMockStorage: func(ctx context.Context) *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", ctx, "Songku").Return("", errTest)
				return mockStorage
			},
			inputCode:   "Songku",
			expectedUrl: "",
			expectedErr: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockStorage := tc.setupMockStorage(ctx)

			shortenURLsvc := NewShortenUrl(mockStorage, nil)
			url, err := shortenURLsvc.GetLinkFromCode(ctx, tc.inputCode)
			assert.Equal(t, tc.expectedUrl, url)
			assert.Equal(t, tc.expectedErr, err)
		})
	}
}
