package user

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_UpdateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		setupDB     func(t *testing.T) *gorm.DB
		inputUser   *model.User
		expectedErr error
		verifyFunc  func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd990",
				DisplayName: "John New Name",
				Email:       "john_new@example.com",
			},
			expectedErr: nil,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				user := &model.User{}
				err := db.First(&user, "id = ?", "deb745af-1a62-4efa-99a0-f06b274bd990").Error
				assert.NoError(t, err)
				assert.Equal(t, "John New Name", user.DisplayName)
				assert.Equal(t, "john_new@example.com", user.Email)
			},
		},
		{
			name: "duplicate email",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUser: &model.User{
				ID:          "deb745af-1a62-4efa-99a0-f06b274bd990",
				DisplayName: "John New Name",
				Email:       "janedoe@example.com",
			},
			expectedErr: dbutils.ErrDuplicationEmail,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				user := model.User{}
				db.First(&user, "id = ?", "deb745af-1a62-4efa-99a0-f06b274bd990")
				assert.Equal(t, "johndoe@example.com", user.Email)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewSQLRepository(db)
			err := repo.UpdateUser(ctx, tc.inputUser)

			assert.ErrorIs(t, tc.expectedErr, err)
			tc.verifyFunc(t, db)
		})
	}
}
