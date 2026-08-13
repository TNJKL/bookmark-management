package utils

import (
	"crypto/rand"
	"math/big"
)

const (
	redisPrefixSet = "abcdefgh"
	sqlPrefixSet   = "ijklmnopqrstuvwxyz"
)

// GetRedisPrefix generates a random single-character prefix from the Redis prefix set ("abcdefgh").
func GetRedisPrefix() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(redisPrefixSet))))
	return string(redisPrefixSet[n.Int64()])
}

// GetSQLPrefix generates a random single-character prefix from the SQL prefix set ("ijklmnopqrstuvwxyz").
func GetSQLPrefix() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(sqlPrefixSet))))
	return string(sqlPrefixSet[n.Int64()])
}

// IsRedisCode checks whether the given short code starts with a valid Redis prefix.
func IsRedisCode(code string) bool {
	if len(code) == 0 {
		return false
	}
	prefix := code[0]
	return prefix >= 'a' && prefix <= 'h'
}

// IsSQLCode checks whether the given short code starts with a valid SQL prefix.
func IsSQLCode(code string) bool {
	if len(code) == 0 {
		return false
	}
	prefix := code[0]
	return prefix >= 'i' && prefix <= 'z'
}
