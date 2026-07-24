package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasher_Hash(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		password    string
		expectedErr error
	}{
		{
			name:        "happy path  normal password",
			password:    "password123",
			expectedErr: nil,
		},
		{
			name:        "password too long - over 72 bytes",
			password:    strings.Repeat("a", 73),
			expectedErr: ErrHashFailed,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			haser := NewHasher()
			hash, err := haser.Hash(tc.password)
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				assert.NotEqual(t, tc.password, hash)
			}
		})
	}
}

func TestHasher_Compare(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name           string
		password       string
		comparePass    string
		isValidHash    bool
		expectedResult bool
	}{
		{
			name:           "correct password",
			password:       "my_secret_password",
			comparePass:    "my_secret_password",
			isValidHash:    true,
			expectedResult: true,
		},
		{
			name:           "incorrect password",
			password:       "my_secret_password",
			comparePass:    "wrong_password",
			isValidHash:    true,
			expectedResult: false,
		},
		{
			name:           "invalid hash format",
			password:       "my_secret_password",
			comparePass:    "my_secret_password",
			isValidHash:    false,
			expectedResult: false,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			hasher := NewHasher()

			var hash string
			if tc.isValidHash {
				var err error
				hash, err = hasher.Hash(tc.password)
				assert.NoError(t, err)
			} else {
				//gia lap chuoi hash rac ko dung dinh dang
				hash = "invalid_hash_format_123"
			}

			result := hasher.Compare(hash, tc.comparePass)
			assert.Equal(t, tc.expectedResult, result)
		})
	}
}
