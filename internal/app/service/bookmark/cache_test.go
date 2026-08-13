package bookmark_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	mock_cache "github.com/TNJKL/bookmark-management/internal/app/repository/cache/mocks"
	bookmarkSvc "github.com/TNJKL/bookmark-management/internal/app/service/bookmark"
	"github.com/TNJKL/bookmark-management/internal/app/service/bookmark/mocks"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

var (
	userID        = "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"
	cacheGroupKey = "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"
	bookmarkID    = "9f763ee7-6dfc-4ede-802f-f78bdc368169"
	desc          = "A bookmark description"
	url           = "http://google.com"
	errTest       = errors.New("test error")

	createdBookmark = &model.Bookmark{
		Base:        fixtures.GetTestBase(bookmarkID),
		Description: desc,
		URL:         url,
		Code:        "wfaA1y97",
		UserID:      userID,
	}
)

func TestBookmarkServiceWithCache_GetBookmarks(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name string

		setupSvc       func(ctx context.Context) *mocks.Service
		setupCache     func(ctx context.Context) *mock_cache.DB
		expectedResult *bookmarkSvc.GetBookmarksResult
		expectErr      error
	}{
		{
			name: "happy path - success with cache",
			setupSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("GetCacheData", ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6", "2_20").Return(
					[]byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y97","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"},{"id":"806fd7e2-0f30-439d-bb5b-071d8377fe9d","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y98","user_id":"da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6"}],"total":132}`),
					nil,
				)
				return mockC
			},
			expectedResult: &bookmarkSvc.GetBookmarksResult{
				Bookmarks: []*model.Bookmark{
					{
						Base:        fixtures.GetTestBase("9f763ee7-6dfc-4ede-802f-f78bdc368169"),
						Description: "A bookmark description",
						URL:         "http://google.com",
						Code:        "wfaA1y97",
					},
					{
						Base:        fixtures.GetTestBase("806fd7e2-0f30-439d-bb5b-071d8377fe9d"),
						Description: "A bookmark description",
						URL:         "http://google.com",
						Code:        "wfaA1y98",
					},
				},
				Total: 132,
			},
			expectErr: nil,
		},

		{
			name: "error case - success - no cache",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockS := mocks.NewService(t)
				mockS.On("GetBookmarks", ctx, "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6", 2, 20).Return(
					&bookmarkSvc.GetBookmarksResult{
						Bookmarks: []*model.Bookmark{
							{
								Base:        fixtures.GetTestBase("9f763ee7-6dfc-4ede-802f-f78bdc368169"),
								Description: "A bookmark description",
								URL:         "http://google.com",
								Code:        "wfaA1y97",
								UserID:      "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
							},
							{
								Base:        fixtures.GetTestBase("806fd7e2-0f30-439d-bb5b-071d8377fe9d"),
								Description: "A bookmark description",
								URL:         "http://google.com",
								Code:        "wfaA1y98",
								UserID:      "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
							},
						},
						Total: 132,
					},
					nil)
				return mockS
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("GetCacheData", ctx, "get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6", "2_20").Return(nil, redis.Nil)
				mockC.On(
					"SetCacheData",
					ctx,
					"get_bookmarks_da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
					"2_20",
					[]byte(`{"bookmarks":[{"id":"9f763ee7-6dfc-4ede-802f-f78bdc368169","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y97"},{"id":"806fd7e2-0f30-439d-bb5b-071d8377fe9d","created_at":"2018-02-18T01:02:03.000000004Z","updated_at":"2018-02-18T01:02:03.000000004Z","description":"A bookmark description","url":"http://google.com","code":"wfaA1y98"}],"total":132}`),
					24*time.Hour).Return(nil)
				return mockC
			},
			expectedResult: &bookmarkSvc.GetBookmarksResult{
				Bookmarks: []*model.Bookmark{
					{
						Base:        fixtures.GetTestBase("9f763ee7-6dfc-4ede-802f-f78bdc368169"),
						Description: "A bookmark description",
						URL:         "http://google.com",
						Code:        "wfaA1y97",
						UserID:      "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
					},
					{
						Base:        fixtures.GetTestBase("806fd7e2-0f30-439d-bb5b-071d8377fe9d"),
						Description: "A bookmark description",
						URL:         "http://google.com",
						Code:        "wfaA1y98",
						UserID:      "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6",
					},
				},
				Total: 132,
			},
			expectErr: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockSvc := tc.setupSvc(ctx)
			mockC := tc.setupCache(ctx)

			cacheSvc := bookmarkSvc.NewServiceWithCache(mockSvc, mockC)
			res, err := cacheSvc.GetBookmarks(ctx, "da72d2a5-dce8-4f2e-ac4b-e0dcdd0eedc6", 2, 20)
			assert.Equal(t, tc.expectErr, err)
			assert.Equal(t, tc.expectedResult, res)
		})

	}
}

func TestBookmarkServiceWithCache_CreateBookmark(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name             string
		setupSvc         func(ctx context.Context) *mocks.Service
		setupCache       func(ctx context.Context) *mock_cache.DB
		expectedBookmark *model.Bookmark
		expectedErr      error
	}{
		{
			name: "happy path",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateBookmark", ctx, desc, url, userID).Return(createdBookmark, nil)
				return mockSvc
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(nil)
				return mockC
			},
			expectedBookmark: createdBookmark,
			expectedErr:      nil,
		},
		{
			name: "fail - delete cache error",
			setupSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(errTest)
				return mockC
			},
			expectedBookmark: nil,
			expectedErr:      errTest,
		},

		{
			name: "fail - service create error",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("CreateBookmark", ctx, desc, url, userID).Return(nil, errTest)
				return mockSvc
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(nil)
				return mockC
			},
			expectedBookmark: nil,
			expectedErr:      errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockSvc := tc.setupSvc(ctx)
			mockC := tc.setupCache(ctx)
			cacheSvc := bookmarkSvc.NewServiceWithCache(mockSvc, mockC)
			res, err := cacheSvc.CreateBookmark(ctx, desc, url, userID)
			assert.Equal(t, tc.expectedBookmark, res)
			assert.Equal(t, tc.expectedErr, err)
		})
	}

}

func TestBookmarkServiceWithCache_UpdateBookmark(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupSvc    func(ctx context.Context) *mocks.Service
		setupCache  func(ctx context.Context) *mock_cache.DB
		expectedErr error
	}{
		{
			name: "happy path",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("UpdateBookmark", ctx, bookmarkID, userID, desc, url).Return(nil)
				return mockSvc
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(nil)
				return mockC
			},
			expectedErr: nil,
		},

		{
			name: "fail - delete cache error",
			setupSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(errTest)
				return mockC
			},
			expectedErr: errTest,
		},
		{
			name: "fail - service update error",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("UpdateBookmark", ctx, bookmarkID, userID, desc, url).Return(errTest)
				return mockSvc
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(nil)
				return mockC
			},
			expectedErr: errTest,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockSvc := tc.setupSvc(ctx)
			mockC := tc.setupCache(ctx)

			cacheSvc := bookmarkSvc.NewServiceWithCache(mockSvc, mockC)
			err := cacheSvc.UpdateBookmark(ctx, bookmarkID, userID, desc, url)

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}

func TestBookmarkServiceWithCache_DeleteBookmark(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupSvc    func(ctx context.Context) *mocks.Service
		setupCache  func(ctx context.Context) *mock_cache.DB
		expectedErr error
	}{
		{
			name: "happy path",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("DeleteBookmark", ctx, bookmarkID, userID).Return(nil)
				return mockSvc
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(nil)
				return mockC
			},
			expectedErr: nil,
		},
		{
			name: "fail - delete cache error",
			setupSvc: func(ctx context.Context) *mocks.Service {
				return mocks.NewService(t)
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(errTest)
				return mockC
			},
			expectedErr: errTest,
		},
		{
			name: "fail - service delete error",
			setupSvc: func(ctx context.Context) *mocks.Service {
				mockSvc := mocks.NewService(t)
				mockSvc.On("DeleteBookmark", ctx, bookmarkID, userID).Return(errTest)
				return mockSvc
			},
			setupCache: func(ctx context.Context) *mock_cache.DB {
				mockC := mock_cache.NewDB(t)
				mockC.On("DeleteCache", ctx, cacheGroupKey).Return(nil)
				return mockC
			},
			expectedErr: errTest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			mockSvc := tc.setupSvc(ctx)
			mockC := tc.setupCache(ctx)

			cacheSvc := bookmarkSvc.NewServiceWithCache(mockSvc, mockC)
			err := cacheSvc.DeleteBookmark(ctx, bookmarkID, userID)

			assert.Equal(t, err, tc.expectedErr)
		})
	}
}
