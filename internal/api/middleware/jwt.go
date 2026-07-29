package middleware

import (
	"net/http"
	"strings"

	"github.com/TNJKL/bookmark-management/pkg/jwtutils"
	"github.com/gin-gonic/gin"
)

// JWTAuth defines the middleware interface for JWT authentication
type JWTAuth interface {
	JWTAuth() gin.HandlerFunc
}

// jwtAuth implements the JWTAuth interface using JWTValidator
type jwtAuth struct {
	jwtVal jwtutils.JWTValidator
}

// NewJWTAuth creates a new JWTAuth middleware instance
func NewJWTAuth(jwtVal jwtutils.JWTValidator) JWTAuth {
	return &jwtAuth{
		jwtVal: jwtVal,
	}
}

// JWTAuth returns a Gin middleware handler function to validate JWT tokens
func (j *jwtAuth) JWTAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//get token from header
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		tokenStr := parts[1]
		//validate token
		tokenClaims, err := j.jwtVal.ValidateJWT(tokenStr)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}
		//set claims to context
		ctx.Set("claims", tokenClaims)
		ctx.Next()
	}
}
