package utils

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBase62(t *testing.T) {
	t.Parallel()
	b := NewBase62(123)
	assert.NotNil(t, b)
}

func TestBase62_EncodeAndDecode(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		inputCode int
	}{
		{
			name:      "zero code",
			inputCode: 0,
		},
		{
			name:      "code equals xorSecret",
			inputCode: 0x3E9A7F2D,
		},
		{
			name:      "code 1",
			inputCode: 1,
		},
		{
			name:      "code 123456",
			inputCode: 123456,
		},
		{
			name:      "large code 999999999",
			inputCode: 999999999,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := NewBase62(123)

			encoded := b.Encode(tc.inputCode)
			assert.NotEmpty(t, encoded)

			decoded := b.Decode(encoded)
			assert.Equal(t, strconv.Itoa(tc.inputCode), decoded)
		})
	}
}

func TestBase62_Decode_InvalidInput(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name     string
		inputStr string
	}{
		{
			name:     "contains special character @",
			inputStr: "abc@123",
		},
		{
			name:     "contains hyphen -",
			inputStr: "-123",
		},
		{
			name:     "contains space",
			inputStr: "abc 123",
		},
		{
			name:     "invalid symbols",
			inputStr: "!#$%^&*",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			b := NewBase62(123)

			decoded := b.Decode(tc.inputStr)
			assert.Equal(t, "", decoded)
		})
	}
}
