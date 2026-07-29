package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestNewJWTValidation(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		keyPath        string
		expectedErrStr string
	}{
		{
			name:           "happy path",
			keyPath:        filepath.FromSlash(publicKeyPath),
			expectedErrStr: "",
		},

		{
			name:           "err case -file not found",
			keyPath:        filepath.FromSlash("./non_existent.key"),
			expectedErrStr: "open",
		},

		{
			name:           "err case - not a public key",
			keyPath:        filepath.FromSlash(privateKeyPath),
			expectedErrStr: "structure error",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewJWTValidator(tc.keyPath)
			if err != nil {
				assert.ErrorContains(t, err, tc.expectedErrStr)
			}
		})
	}
}

func TestValidator_ValidateJWT(t *testing.T) {
	val, err := NewJWTValidator(filepath.FromSlash(publicKeyPath))
	if err != nil {
		t.Fatal("should not fail")
	}
	// Reused the sharedTestToken defined in generator_test.go to avoid duplication
	claims, err := val.ValidateJWT(sharedTestToken)
	if err != nil {
		t.Fatal("should not fail")
	}
	assert.Equal(t, claims, jwt.MapClaims{
		"sub": "test123",
	})
}
