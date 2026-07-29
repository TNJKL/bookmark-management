package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// sharedTestToken is used across both generator_test.go and validator_test.go
// to avoid duplicating the long signature string literal.
const sharedTestToken = "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0MTIzIn0.LXoNKqv4WoBu2JH0ELzfu_leagME-CWZ_MZ-rNnl-zzka4gnrRYDukXVE3LbElHbvTqtZdxp3TRo3YbFjiQR0QkbmbcL63y1dVwTuNOpUf04mMN3HMe2AZyfRUakfz4-7xzx0swi5GskT3axSuis42dQzdg7isl7OttKAg20_ti-dZ23VX6m5M0uOUgG6MHbG7qNtbztpRZ3Xl0dtQI-zYPGZpXqS079ydjSQzKOdyVcZq4U8yFdp1Imf5PXe3hVOYnA6Xdo3kp9ZlMKhiOKERo9QwAn5aOR0BzRBFzR7jBxNEgE6cQ6NN7QVY_L7KuxqZQRHqmbTybDY2ziEeHhVA"
const (
	privateKeyPath = "./test.private.key"
	publicKeyPath  = "./test.public.key"
)

func TestNewJWTGenerator(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		keyPath string

		expectedErrStr string
	}{
		{
			name:           "happy path",
			keyPath:        filepath.FromSlash(privateKeyPath),
			expectedErrStr: "",
		},
		{
			name:           "err case - file not found",
			keyPath:        filepath.FromSlash("./non_existent.key"),
			expectedErrStr: "open",
		},
		{
			name:           "err case - not a private key",
			keyPath:        filepath.FromSlash(publicKeyPath),
			expectedErrStr: "structure error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewJWTGenerator(tc.keyPath)
			if err != nil {
				assert.ErrorContains(t, err, tc.expectedErrStr)
			}
		})
	}
}

func TestGenerator_GenerateJWT(t *testing.T) {
	gen, err := NewJWTGenerator(filepath.FromSlash(privateKeyPath))
	if err != nil {
		t.Fatal("should not fail")
	}
	token, err := gen.GenerateJWT(map[string]any{
		"sub": "test123",
	})
	if err != nil {
		t.Fatal("should not fail")
	}
	assert.Equal(t, token, sharedTestToken)
}
