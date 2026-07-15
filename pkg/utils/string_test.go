package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKeyGenerator_GenerateKey(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name           string
		expectedLength int
	}{
		{
			name:           "normal case - length 7",
			expectedLength: 7,
		},
		{
			name:           "normal case - length 12",
			expectedLength: 12,
		},
		{
			name:           "normal case - length 100",
			expectedLength: 100,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			generator := NewKeyGenerator()
			key := generator.GenerateKey(tc.expectedLength)
			assert.Equal(t, tc.expectedLength, len(key))
			for _, char := range key {
				assert.Contains(t, charset, string(char))
			}
		})
	}

}
