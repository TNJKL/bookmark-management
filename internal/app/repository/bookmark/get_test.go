package bookmark

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_GetBookmarks(t *testing.T) {
	t.Parallel()

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserID   string
		limit         int
		offset        int
		expectedCount int64
		expectedErr   error
		verifyFunc    func(t *testing.T, res []*model.Bookmark)
	}{
		{
			name:          "happy path - get all 5 bookmarks for user",
			setupDB:       func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			inputUserID:   userID,
			limit:         10,
			offset:        0,
			expectedCount: 5,
			expectedErr:   nil,
			verifyFunc: func(t *testing.T, res []*model.Bookmark) {
				assert.Len(t, res, 5)
				assert.Equal(t, "123456", res[0].Code)
				assert.Equal(t, "123457", res[1].Code)
				assert.Equal(t, "123458", res[2].Code)
				assert.Equal(t, "123459", res[3].Code)
				assert.Equal(t, "123460", res[4].Code)
				assert.Equal(t, userID, res[0].UserID)
			},
		},
		{
			name:          "pagination page 1 - limit 2 offset 0",
			setupDB:       func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			inputUserID:   userID,
			limit:         2,
			offset:        0,
			expectedCount: 5,
			expectedErr:   nil,
			verifyFunc: func(t *testing.T, res []*model.Bookmark) {
				assert.Len(t, res, 2)
				assert.Equal(t, "123456", res[0].Code)
				assert.Equal(t, "123457", res[1].Code)
			},
		},
		{
			name:          "pagination page 2 - limit 2 offset 2",
			setupDB:       func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			inputUserID:   userID,
			limit:         2,
			offset:        2,
			expectedCount: 5,
			expectedErr:   nil,
			verifyFunc: func(t *testing.T, res []*model.Bookmark) {
				assert.Len(t, res, 2)
				assert.Equal(t, "123458", res[0].Code)
				assert.Equal(t, "123459", res[1].Code)
			},
		},
		{
			name:          "user has no bookmarks",
			setupDB:       func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			inputUserID:   "non-existent-user-id",
			limit:         10,
			offset:        0,
			expectedCount: 0,
			expectedErr:   nil,
			verifyFunc: func(t *testing.T, res []*model.Bookmark) {
				assert.Empty(t, res)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			res, count, err := repo.GetBookmarks(ctx, tc.inputUserID, tc.limit, tc.offset)

			assert.ErrorIs(t, err, tc.expectedErr)
			assert.Equal(t, tc.expectedCount, count)
			tc.verifyFunc(t, res)
		})
	}
}
