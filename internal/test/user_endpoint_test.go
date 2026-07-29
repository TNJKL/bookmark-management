package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/internal/model"
	"github.com/TNJKL/bookmark-management/internal/test/data/fixtures"
	"github.com/TNJKL/bookmark-management/pkg/jwtutils"
	"github.com/TNJKL/bookmark-management/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

			testAPI := buildTestAPI(db, nil, nil)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			tc.verifyFunc(t, db)
		})
	}
}

type loginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestLoginEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
				hasher := utils.NewHasher()
				hashedPassword, _ := hasher.Hash("Kakarot1996@")
				db.Model(&model.User{}).Where("username = ?", "johndoe").Update("password", hashedPassword)
				return db
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := loginInput{
					Username: "johndoe",
					Password: "Kakarot1996@",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Logged in successfully"`,
		},
		{
			name: "invalid credentials - wrong password",
			setupDB: func(t *testing.T) *gorm.DB {
				db := fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
				hasher := utils.NewHasher()
				hashedPassword, _ := hasher.Hash("Kakarot1996@")
				db.Model(&model.User{}).Where("username = ?", "johndoe").Update("password", hashedPassword)
				return db
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				input := loginInput{
					Username: "johndoe",
					Password: "Wrongpassword@123",
				}
				body, _ := json.Marshal(input)
				req := httptest.NewRequest(http.MethodPost, "/v1/users/login", bytes.NewBuffer(body))
				req.Header.Set("Content-Type", "application/json")
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `"message":"invalid credentials"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)

			testAPI := buildTestAPI(db, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
		})
	}
}

func TestGetSelfInfoEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

	testCases := []struct {
		name               string
		setupDB            func(t *testing.T) *gorm.DB
		setupTestHTTP      func(api api.Engine) *httptest.ResponseRecorder
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name: "happy path",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, "deb745af-1a62-4efa-99a0-f06b274bd990", "johndoe@example.com")
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				req.Header.Set("Authorization", "Bearer "+token)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"username":"johndoe"`,
		},
		{
			name: "unauthorized - no token",
			setupDB: func(t *testing.T) *gorm.DB {
				return fixtures.NewFixture(t, &fixtures.UserCommonTestDB{})
			},
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				req := httptest.NewRequest(http.MethodGet, "/v1/self/info", nil)
				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse:   `"error":"Unauthorized"`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			db := tc.setupDB(t)
			testAPI := buildTestAPI(db, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)
			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
		})
	}
}

func TestUpdateSelfInfoEndpoint(t *testing.T) {
	t.Parallel()
	jwtGen, jwtVal := setupTestJWT(t)

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
				token := generateTestToken(t, jwtGen, "deb745af-1a62-4efa-99a0-f06b274bd990", "johndoe@example.com")

				body := bytes.NewBufferString(`{"display_name":"John New Name","email":"john_new@example.com"}`)
				req := httptest.NewRequest(http.MethodPut, "/v1/self/info", body)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusOK,
			expectedResponse:   `"message":"Edit current user successfully!"`,
			verifyFunc: func(t *testing.T, db *gorm.DB) {
				user := model.User{}
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
			setupTestHTTP: func(api api.Engine) *httptest.ResponseRecorder {
				token := generateTestToken(t, jwtGen, "deb745af-1a62-4efa-99a0-f06b274bd990", "johndoe@example.com")

				body := bytes.NewBufferString(`{"display_name":"John New Name","email":"janedoe@example.com"}`)
				req := httptest.NewRequest(http.MethodPut, "/v1/self/info", body)
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("Authorization", "Bearer "+token)

				rec := httptest.NewRecorder()
				api.ServerHTTP(rec, req)
				return rec
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse:   `"message":"Email already taken"`,
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
			db := tc.setupDB(t)

			testAPI := buildTestAPI(db, jwtGen, jwtVal)
			recorder := tc.setupTestHTTP(testAPI)

			assert.Equal(t, tc.expectedStatusCode, recorder.Code)
			assert.Contains(t, recorder.Body.String(), tc.expectedResponse)
			tc.verifyFunc(t, db)
		})
	}
}

// setupTestJWT loads the test RSA keys from pkg/jwtutils
func setupTestJWT(t *testing.T) (jwtutils.JWTGenerator, jwtutils.JWTValidator) {
	jwtGen, err := jwtutils.NewJWTGenerator("../../pkg/jwtutils/test.private.key")
	assert.NoError(t, err)
	jwtVal, err := jwtutils.NewJWTValidator("../../pkg/jwtutils/test.public.key")
	assert.NoError(t, err)
	return jwtGen, jwtVal
}

// generateTestToken creates and signs a valid JWT token for test purposes
func generateTestToken(t *testing.T, jwtGen jwtutils.JWTGenerator, sub, email string) string {
	tokenContent := jwt.MapClaims{
		"sub":   sub,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token, err := jwtGen.GenerateJWT(tokenContent)
	assert.NoError(t, err)
	return token
}

// buildTestAPI instantiates the Gin engine with mocking dependencies
func buildTestAPI(db *gorm.DB, jwtGen jwtutils.JWTGenerator, jwtVal jwtutils.JWTValidator) api.Engine {
	return api.NewEngine(&api.EngineOpts{
		App:    gin.New(),
		Cfg:    &api.Config{},
		Db:     db,
		JWTGen: jwtGen,
		JWTVal: jwtVal,
	})
}
