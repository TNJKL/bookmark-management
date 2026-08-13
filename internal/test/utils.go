package integration

import (
	"testing"
	"time"

	"github.com/TNJKL/bookmark-management/internal/api"
	"github.com/TNJKL/bookmark-management/pkg/jwtutils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

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
func buildTestAPI(db *gorm.DB, redisClient *redis.Client, jwtGen jwtutils.JWTGenerator, jwtVal jwtutils.JWTValidator) api.Engine {
	return api.NewEngine(&api.EngineOpts{
		App:         gin.New(),
		Cfg:         &api.Config{},
		RedisClient: redisClient,
		Db:          db,
		JWTGen:      jwtGen,
		JWTVal:      jwtVal,
	})
}
