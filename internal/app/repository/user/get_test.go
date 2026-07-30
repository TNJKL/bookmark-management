package user

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_GetUserByUsername(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserName string
		expectedErr   error
		verifyFunc    func(t *testing.T, user *model.User)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserName: "johndoe",
			expectedErr:   nil,
			verifyFunc: func(t *testing.T, user *model.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "johndoe", user.Username)
				assert.Equal(t, "johndoe@example.com", user.Email)
				assert.Equal(t, "John Doe", user.DisplayName)
			},
		},

		{
			name: "user not found	",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserName: "songoku",
			expectedErr:   dbutils.ErrRecordNotFound,
			verifyFunc: func(t *testing.T, user *model.User) {
				assert.Nil(t, user)
			},
		},

		{
			name: "empty username",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserName: "",
			expectedErr:   dbutils.ErrRecordNotFound,
			verifyFunc: func(t *testing.T, user *model.User) {
				assert.Nil(t, user)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewSQLRepository(db)

			user, err := repo.GetUserByUsername(ctx, tc.inputUserName)

			assert.ErrorIs(t, tc.expectedErr, err)
			tc.verifyFunc(t, user)
		})
	}
}

func TestSqlRepository_GetUserByID(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name        string
		setupDB     func(t *testing.T) *gorm.DB
		inputUserID string
		expectedErr error
		verifyFunc  func(t *testing.T, user *model.User)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserID: "deb745af-1a62-4efa-99a0-f06b274bd990",
			expectedErr: nil,
			verifyFunc: func(t *testing.T, user *model.User) {
				assert.NotNil(t, user)
				assert.Equal(t, "johndoe", user.Username)
				assert.Equal(t, "johndoe@example.com", user.Email)
				assert.Equal(t, "John Doe", user.DisplayName)
			},
		},

		{
			name: "user id not found",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserID: "deb745af-1a62-4efa-99a0-f06b274bd123",
			expectedErr: dbutils.ErrRecordNotFound,
			verifyFunc: func(t *testing.T, user *model.User) {
				assert.Nil(t, user)
			},
		},

		{
			name: "empty userID",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserID: "",
			expectedErr: dbutils.ErrRecordNotFound,
			verifyFunc: func(t *testing.T, user *model.User) {
				assert.Nil(t, user)
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewSQLRepository(db)
			user, err := repo.GetUserByID(ctx, tc.inputUserID)
			assert.ErrorIs(t, tc.expectedErr, err)
			tc.verifyFunc(t, user)
		})
	}

}
