package jwtutils

import (
	"crypto/rsa"
	"errors"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTValidator defines the interface for validating JWT tokens
//
//go:generate mockery --name JWTValidator --filename jwtValidator.go
type JWTValidator interface {
	ValidateJWT(tokenStr string) (jwt.MapClaims, error)
}

// jwtValidator implements the JWTValidator interface using RSA public key
type jwtValidator struct {
	publicKey *rsa.PublicKey
}

// NewJWTValidator creates a new JWTValidator instance with the provided public key path
func NewJWTValidator(publicKeyPath string) (JWTValidator, error) {
	publicKeyData, err := os.ReadFile(publicKeyPath)
	if err != nil {
		return nil, err
	}

	publicKey, err := jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
	if err != nil {
		return nil, err
	}
	return &jwtValidator{
		publicKey: publicKey,
	}, nil
}

// Common JWT validation errors
var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExtractToken = errors.New("failed to extract token")
)

// ValidateJWT parses and verifies the token string, returning its claims if valid
func (v *jwtValidator) ValidateJWT(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return v.publicKey, nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		return claims, nil
	}

	return nil, ErrExtractToken
}
