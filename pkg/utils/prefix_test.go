package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRedisPrefix(t *testing.T) {
	t.Parallel()
	for i := 0; i < 50; i++ {
		prefix := GetRedisPrefix()
		assert.Len(t, prefix, 1)
		assert.True(t, strings.Contains(redisPrefixSet, prefix))
		assert.True(t, IsRedisCode(prefix))
		assert.False(t, IsSQLCode(prefix))
	}
}

func TestGetSQLPrefix(t *testing.T) {
	t.Parallel()
	for i := 0; i < 50; i++ {
		prefix := GetSQLPrefix()
		assert.Len(t, prefix, 1)
		assert.True(t, strings.Contains(sqlPrefixSet, prefix))
		assert.True(t, IsSQLCode(prefix))
		assert.False(t, IsRedisCode(prefix))
	}
}

func TestIsRedisCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid redis prefix 'a'",
			input:    "aSongoku",
			expected: true,
		},
		{
			name:     "valid redis prefix 'h'",
			input:    "h123456",
			expected: true,
		},
		{
			name:     "sql prefix 'i' (boundary outside redis range)",
			input:    "i123456",
			expected: false,
		},
		{
			name:     "sql prefix 'z'",
			input:    "z123456",
			expected: false,
		},
		{
			name:     "numeric prefix '1'",
			input:    "123456",
			expected: false,
		},
		{
			name:     "uppercase prefix 'A'",
			input:    "A123456",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := IsRedisCode(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestIsSQLCode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "valid sql prefix 'i'",
			input:    "i123456",
			expected: true,
		},
		{
			name:     "valid sql prefix 'z'",
			input:    "z123456",
			expected: true,
		},
		{
			name:     "redis prefix 'a'",
			input:    "aSongoku",
			expected: false,
		},
		{
			name:     "redis prefix 'h' (boundary outside sql range)",
			input:    "h123456",
			expected: false,
		},
		{
			name:     "numeric prefix '1'",
			input:    "123456",
			expected: false,
		},
		{
			name:     "uppercase prefix 'I'",
			input:    "I123456",
			expected: false,
		},
		{
			name:     "empty string",
			input:    "",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := IsSQLCode(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
