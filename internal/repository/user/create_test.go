package user

import (
	"testing"

	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/dbutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSqlRepository_CreateUser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		setupDB       func(t *testing.T) *gorm.DB
		inputUserName string
		inputEmail    string
		expectedErr   error
		verifyFunc    func(db *gorm.DB, username, email string)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserName: "user1",
			inputEmail:    "user1@gmail.com",
			expectedErr:   nil,
			verifyFunc: func(db *gorm.DB, username, email string) {
				user := model.User{}
				err := db.First(&user, "username = ?", username).Error
				assert.NoError(t, err)
				assert.Equal(t, username, user.Username)
				assert.Equal(t, email, user.Email)
				assert.NotEmpty(t, user.ID)
			},
		},
		{
			name: "duplicate username",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserName: "johndoe",
			inputEmail:    "test123@gmail.com",
			expectedErr:   dbutils.ErrDuplicationUsername,
			verifyFunc: func(db *gorm.DB, username, email string) {
				user := model.User{}
				err := db.First(&user, "email = ?", email).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
		{
			name: "duplicate email",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			inputUserName: "test123",
			inputEmail:    "johndoe@example.com",
			expectedErr:   dbutils.ErrDuplicationEmail,
			verifyFunc: func(db *gorm.DB, username, email string) {
				user := model.User{}
				err := db.First(&user, "username = ?", username).Error
				assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			ctx := t.Context()
			db := tc.setupDB(t)
			repo := NewSQLRepository(db)
			user, err := repo.CreateUser(ctx, &model.User{
				Username: tc.inputUserName,
				Email:    tc.inputEmail,
			})
			if err != nil {
				assert.Nil(t, user)
			}
			assert.ErrorIs(t, err, tc.expectedErr)
			tc.verifyFunc(db, tc.inputUserName, tc.inputEmail)
		})
	}
}
