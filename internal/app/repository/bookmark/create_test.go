package bookmark

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_CreateBookmark(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupDB     func(t *testing.T) *gorm.DB
		input       *model.Bookmark
		expectedErr error
		verifyFunc  func(db *gorm.DB, expect *model.Bookmark)
	}{
		{
			name:    "happy path",
			setupDB: func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			input: &model.Bookmark{
				Base:        fixtures.GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd996"),
				Description: "bookmark6",
				URL:         "https://google.com",
				Code:        "666666",
				UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
			},
			expectedErr: nil,
			verifyFunc: func(db *gorm.DB, expect *model.Bookmark) {
				bookmark := &model.Bookmark{}
				err := db.Where("code = ?", "666666").First(bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, expect.ID, bookmark.ID)
				assert.Equal(t, expect.Description, bookmark.Description)
				assert.Equal(t, expect.Code, bookmark.Code)
				assert.Equal(t, expect.URL, bookmark.URL)
			},
		},
		{
			name:    "duplicate code",
			setupDB: func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			input: &model.Bookmark{
				Base:        fixtures.GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd997"),
				Description: "duplicate code test",
				URL:         "https://example.com",
				Code:        "123456",
				UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
			},
			expectedErr: dbutils.ErrUniqueConstraint,
			verifyFunc: func(db *gorm.DB, expect *model.Bookmark) {
				assert.Nil(t, expect)
			},
		},
		{
			name:    "create bookmark with empty description",
			setupDB: func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			input: &model.Bookmark{
				Base:        fixtures.GetTestBase("deb745af-1a62-4efa-99a0-f06b274bd999"),
				Description: "",
				URL:         "https://github.com",
				Code:        "777777",
				UserID:      "deb745af-1a62-4efa-99a0-f06b274bd990",
			},
			expectedErr: nil,
			verifyFunc: func(db *gorm.DB, expect *model.Bookmark) {
				bookmark := &model.Bookmark{}
				err := db.Where("code = ?", "777777").First(bookmark).Error
				assert.NoError(t, err)
				assert.Equal(t, "", bookmark.Description)
				assert.Equal(t, "777777", bookmark.Code)
				assert.Equal(t, "https://github.com", bookmark.URL)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)
			res, err := repo.CreateBookmark(ctx, tc.input)
			assert.ErrorIs(t, err, tc.expectedErr)
			tc.verifyFunc(db, res)
		})
	}
}
