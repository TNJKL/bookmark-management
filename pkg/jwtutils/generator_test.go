package jwtutils

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
			keyPath:        filepath.FromSlash("./test.private.key"),
			expectedErrStr: "",
		},
		{
			name:           "err case - file not found",
			keyPath:        filepath.FromSlash("./non_existent.key"),
			expectedErrStr: "open",
		},
		{
			name:           "err case - not a private key",
			keyPath:        filepath.FromSlash("./test.public.key"),
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
	expectedToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ0ZXN0MTIzIn0.LXoNKqv4WoBu2JH0ELzfu_leagME-CWZ_MZ-rNnl-zzka4gnrRYDukXVE3LbElHbvTqtZdxp3TRo3YbFjiQR0QkbmbcL63y1dVwTuNOpUf04mMN3HMe2AZyfRUakfz4-7xzx0swi5GskT3axSuis42dQzdg7isl7OttKAg20_ti-dZ23VX6m5M0uOUgG6MHbG7qNtbztpRZ3Xl0dtQI-zYPGZpXqS079ydjSQzKOdyVcZq4U8yFdp1Imf5PXe3hVOYnA6Xdo3kp9ZlMKhiOKERo9QwAn5aOR0BzRBFzR7jBxNEgE6cQ6NN7QVY_L7KuxqZQRHqmbTybDY2ziEeHhVA"

	gen, err := NewJWTGenerator(filepath.FromSlash("./test.private.key"))
	if err != nil {
		t.Fatal("should not fail")
	}
	token, err := gen.GenerateJWT(map[string]any{
		"sub": "test123",
	})
	if err != nil {
		t.Fatal("should not fail")
	}
	assert.Equal(t, token, expectedToken)
}
