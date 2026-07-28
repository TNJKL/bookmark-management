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
			keyPath:        filepath.FromSlash("./test.public.key"),
			expectedErrStr: "",
		},

		{
			name:           "err case -file not found",
			keyPath:        filepath.FromSlash("./non_existent.key"),
			expectedErrStr: "open",
		},

		{
			name:           "err case - not a public key",
			keyPath:        filepath.FromSlash("./test.private.key"),
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
	testToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0MTIzIn0.LXoNKqv4WoBu2JH0ELzfu_leagME-CWZ_MZ-rNnl-zzka4gnrRYDukXVE3LbElHbvTqtZdxp3TRo3YbFjiQR0QkbmbcL63y1dVwTuNOpUf04mMN3HMe2AZyfRUakfz4-7xzx0swi5GskT3axSuis42dQzdg7isl7OttKAg20_ti-dZ23VX6m5M0uOUgG6MHbG7qNtbztpRZ3Xl0dtQI-zYPGZpXqS079ydjSQzKOdyVcZq4U8yFdp1Imf5PXe3hVOYnA6Xdo3kp9ZlMKhiOKERo9QwAn5aOR0BzRBFzR7jBxNEgE6cQ6NN7QVY_L7KuxqZQRHqmbTybDY2ziEeHhVA"

	val, err := NewJWTValidator(filepath.FromSlash("./test.public.key"))
	if err != nil {
		t.Fatal("should not fail")
	}
	claims, err := val.ValidateJWT(testToken)
	if err != nil {
		t.Fatal("should not fail")
	}
	assert.Equal(t, claims, jwt.MapClaims{
		"sub": "test123",
	})
}
