package utils

import "strconv"

const char = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

// Base62 defines operations for encoding integer IDs to obfuscated Base62 strings
// and decoding them back to integer strings.
//
//go:generate mockery --name Base62 --filename base62.go
type Base62 interface {
	// Encode converts an integer ID into a Base62 encoded string.
	Encode(code int) string
	// Decode converts a Base62 encoded string back to its original integer string representation.
	Decode(str string) string
}

// base62 is the default implementation of the Base62 interface.
type base62 struct {
	xorSecret int
}

// NewBase62 creates a new instance of the Base62 service.
func NewBase62(secret int) Base62 {
	return &base62{
		xorSecret: secret,
	}
}

// Encode converts an integer code into an obfuscated Base62 string using XOR obfuscation.
func (b *base62) Encode(code int) string {
	xorInput := code ^ b.xorSecret

	if xorInput == 0 {
		return string(char[0])
	}

	result := make([]byte, 0, 8)

	for xorInput > 0 {
		remainder := xorInput % 62

		result = append(result, char[remainder])

		xorInput /= 62
	}
	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return string(result)
}

// Decode parses a Base62 string back into its original integer string representation.
// Returns an empty string if str contains invalid Base62 characters.
func (b *base62) Decode(str string) string {
	var n int
	for _, ch := range str {
		index := -1
		for i, c := range char {
			if c == ch {
				index = i
				break
			}
		}

		if index == -1 {
			return ""
		}

		n = n*62 + index
	}

	decoded := n ^ b.xorSecret

	return strconv.Itoa(decoded)
}
