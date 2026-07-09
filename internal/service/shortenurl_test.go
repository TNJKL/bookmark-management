package service

import (
	"context"
	"errors"
	"testing"
	"time"

	repoMocks "github.com/TNJKL/bookmark-management/internal/repository/mocks"
	serviceMocks "github.com/TNJKL/bookmark-management/internal/service/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var errTest = errors.New("test error")

func TestShortenURL_CreateShortenLink(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		setupMockGenPass func() *serviceMocks.GenPass
		setupMockStorage func() *repoMocks.URLStorage
		inputURL         string
		inputExp         time.Duration
		expectedCode     string
		expectedErr      error
	}{
		{
			name: "generate code error",
			setupMockGenPass: func() *serviceMocks.GenPass {
				mockGenPass := serviceMocks.NewGenPass(t)
				mockGenPass.On("GeneratePassword", 7).Return("", errTest)
				return mockGenPass
			},
			setupMockStorage: func() *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				return mockStorage
			},
			inputURL:     "https://google.com",
			inputExp:     time.Hour,
			expectedCode: "",
			expectedErr:  errTest,
		},
		{
			name: "store Url error",
			setupMockGenPass: func() *serviceMocks.GenPass {
				mockGenPass := serviceMocks.NewGenPass(t)
				mockGenPass.On("GeneratePassword", 7).Return("Songoku", nil)
				return mockGenPass
			},
			setupMockStorage: func() *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", mock.Anything, "Songoku").Return("", redis.Nil)
				mockStorage.On("StoreURL", mock.Anything, "Songoku", "https://google.com", time.Hour).Return(errTest)
				return mockStorage
			},
			inputURL:     "https://google.com",
			inputExp:     time.Hour,
			expectedCode: "",
			expectedErr:  errTest,
		},
		{
			name: "happy path",
			setupMockGenPass: func() *serviceMocks.GenPass {
				mockGenPass := serviceMocks.NewGenPass(t)
				mockGenPass.On("GeneratePassword", 7).Return("Songoku", nil)
				return mockGenPass
			},

			setupMockStorage: func() *repoMocks.URLStorage {
				mockStorage := repoMocks.NewURLStorage(t)
				mockStorage.On("GetURL", mock.Anything, "Songoku").Return("", redis.Nil)
				mockStorage.On("StoreURL", mock.Anything, "Songoku", "https://google.com", time.Hour).Return(nil)
				return mockStorage
			},
			inputURL:     "https://google.com",
			inputExp:     time.Hour,
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
			mockGenPass := tc.setupMockGenPass()
			mockStorage := tc.setupMockStorage()

			//tao ShortenURL service , truyen mock vao
			shortenURLsvc := NewShortenUrl(mockStorage, mockGenPass)

			//goi ham can test
			code, err := shortenURLsvc.CreateShortenLink(ctx, tc.inputURL, tc.inputExp)
			assert.Equal(t, tc.expectedCode, code)
			assert.Equal(t, tc.expectedErr, err)

		})
	}

}
