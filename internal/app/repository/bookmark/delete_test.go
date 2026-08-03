package bookmark

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_DeleteBookmark(t *testing.T) {
	t.Parallel()

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"
	bookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd993"

	testCases := []struct {
		name        string
		setupDB     func(t *testing.T) *gorm.DB
		id          string
		userID      string
		expectedErr error
		verifyFunc  func(t *testing.T, db *gorm.DB)
	}{
		{
			name:        "happy path",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			id:          bookmarkID,
			userID:      userID,
			expectedErr: nil,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&model.Bookmark{}).Where("id = ?", bookmarkID).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name:        "record not found - wrong id",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			id:          "non-existent-id",
			userID:      userID,
			expectedErr: dbutils.ErrRecordNotFound,
		},
		{
			name:        "record not found - wrong user id",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			id:          bookmarkID,
			userID:      "other-user-id",
			expectedErr: dbutils.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)
			err := repo.DeleteBookmark(ctx, tc.id, tc.userID)
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db)
			}
		})
	}
}
