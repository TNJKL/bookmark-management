package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

type registerInput struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

func TestRegisterEndpoint(t *testing.T) {
	t.Parallel()
	newUsername := "newuser"
	newPassword := "Kakarot1996@"
	newDisplayName := "New User"
	newEmail := "newuser@example.com"

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
		verifyFunc         func(t *testing.T, db *gorm.DB)
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := registerInput{
					Username:    newUsername,
					Password:    newPassword,
					DisplayName: newDisplayName,
					Email:       newEmail,
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"Register an user successfully!"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				user := model.User{}
				err := db.First(&user, "username = ?", newUsername).Error
				assert.NoError(t, err)
				assert.Equal(t, newUsername, user.Username)
				assert.Equal(t, newEmail, user.Email)
				assert.NotEqual(t, newPassword, user.Password)
				assert.NotEmpty(t, user.Password)
			},
		},
		{
			name: "invalid input - missing password",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := registerInput{
					Username:    newUsername,
					DisplayName: newDisplayName,
					Email:       newEmail,
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"Input error"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&model.User{}).Where("username = ? OR email = ?", newUsername, newEmail).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "duplicate username",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := registerInput{
					Username:    "johndoe",
					Password:    newPassword,
					DisplayName: newDisplayName,
					Email:       newEmail,
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `"message":"Username already taken"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&model.User{}).Where("username = ? OR email = ?", newUsername, newEmail).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
		{
			name: "duplicate email",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := registerInput{
					Username:    newUsername,
					Password:    newPassword,
					DisplayName: newDisplayName,
					Email:       "janedoe@example.com",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/register", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `"message":"Email already taken"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				var count int64
				db.Model(&model.User{}).Where("username = ? OR email = ?", newUsername, newEmail).Count(&count)
				assert.Equal(t, int64(0), count)
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)

			testAPI := api.NewEngine(&api.Config{}, nil, db)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			tc.verifyFunc(t, db)
		})
	}
}
