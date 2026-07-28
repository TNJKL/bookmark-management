package jwtutils

import (
	"crypto/rsa"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

// JWTGenerator defines the interface for creating JWT tokens
//
//go:generate mockery --name JWTGenerator --filename jwtGenerator.go
type JWTGenerator interface {
	GenerateJWT(jwtContent jwt.MapClaims) (string, error)
}

// jwtGenerator implements the JWTGenerator interface using RSA private key
type jwtGenerator struct {
	privateKey *rsa.PrivateKey
}

// NewJWTGenerator creates a new JWTGenerator instance with the provided private key path
func NewJWTGenerator(privateKeyPath string) (JWTGenerator, error) {
	privateKeyData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, err
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
	if err != nil {
		return nil, err
	}

	return &jwtGenerator{
		privateKey: privateKey,
	}, nil
}

// GenerateJWT signs and returns a new JWT token string using the private key
func (g *jwtGenerator) GenerateJWT(jwtContent jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwtContent)
	tokenString, err := token.SignedString(g.privateKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
