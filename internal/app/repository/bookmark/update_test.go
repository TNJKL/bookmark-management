package bookmark

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_UpdateBookmark(t *testing.T) {
	t.Parallel()

	userID := "deb745af-1a62-4efa-99a0-f06b274bd990"     //ID of user in fixtures
	bookmarkID := "deb745af-1a62-4efa-99a0-f06b274bd993" // ID of bookmark1 in fixtures

	testCases := []struct {
		name        string
		setupDB     func(t *testing.T) *gorm.DB
		id          string
		userID      string
		description string
		url         string
		expectedErr error
		verifyFunc  func(t *testing.T, db *gorm.DB)
	}{
		{
			name:        "happy path",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			id:          bookmarkID,
			userID:      userID,
			description: "Updated Google Description",
			url:         "https://updated-google.com",
			expectedErr: nil,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var bm model.Bookmark
				err := db.First(&bm, "id = ?", bookmarkID).Error
				assert.NoError(t, err)
				assert.Equal(t, "Updated Google Description", bm.Description)
				assert.Equal(t, "https://updated-google.com", bm.URL)
			},
		},
		{
			name:        "record not found - wrong id",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			id:          "non-existent-id",
			userID:      userID,
			description: "Updated Google Description",
			url:         "https://updated-google.com",
			expectedErr: dbutils.ErrRecordNotFound,
		},
		{
			name:        "record not found - wrong user id",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			id:          bookmarkID,
			userID:      "other-user-id",
			description: "Updated Google Description",
			url:         "https://updated-google.com",
			expectedErr: dbutils.ErrRecordNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)
			err := repo.UpdateBookmark(ctx, tc.id, tc.userID, tc.description, tc.url)
			assert.ErrorIs(t, err, tc.expectedErr)
			if tc.verifyFunc != nil {
				tc.verifyFunc(t, db)
			}
		})
	}
}
