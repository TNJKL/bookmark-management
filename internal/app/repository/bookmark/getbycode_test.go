package bookmark

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/TNJKL/bookmark-management/pkg/sqldb"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBookmarkRepo_GetByCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupDB     func(t *testing.T) *gorm.DB
		inputCode   string
		expectedErr error
		verifyFunc  func(t *testing.T, bm *model.Bookmark)
	}{
		{
			name:        "happy path - bookmark found by code",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			inputCode:   "123456",
			expectedErr: nil,
			verifyFunc: func(t *testing.T, bm *model.Bookmark) {
				assert.NotNil(t, bm)
				assert.Equal(t, "123456", bm.Code)
				assert.Equal(t, "bookmark1", bm.Description)
				assert.Equal(t, "https://google.com", bm.URL)
				assert.Equal(t, "deb745af-1a62-4efa-99a0-f06b274bd990", bm.UserID)
			},
		},
		{
			name:        "record not found",
			setupDB:     func(t *testing.T) *gorm.DB { return fixtures.NewFixture(t, &fixtures.BookmarkCommonTestDB{}) },
			inputCode:   "non_existent_code",
			expectedErr: dbutils.ErrRecordNotFound,
			verifyFunc: func(t *testing.T, bm *model.Bookmark) {
				assert.Nil(t, bm)
			},
		},
		{
			name: "db error - table does not exist",
			setupDB: func(t *testing.T) *gorm.DB {
				return sqldb.InitMockDB(t)
			},
			inputCode:   "123456",
			expectedErr: nil,
			verifyFunc: func(t *testing.T, bm *model.Bookmark) {
				assert.Nil(t, bm)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewRepository(db)

			bm, err := repo.GetByCode(ctx, tc.inputCode)

			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else if tc.name == "db error - table does not exist" {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tc.verifyFunc != nil {
				tc.verifyFunc(t, bm)
			}
		})
	}
}
