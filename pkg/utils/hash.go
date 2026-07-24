package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// Hasher defines the contract for secure password hashing and comparison.
//
//go:generate mockery --name Hasher --filename hasher.go
type Hasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) bool
}

// hasher implements the Hasher interface using the bcrypt algorithm.
type hasher struct{}

// NewHasher creates a new instance of Hasher.
func NewHasher() Hasher {
	return &hasher{}
}

// ErrHashFailed is returned when the password hashing process fails
var ErrHashFailed = errors.New("failed to hash password")

// Hash encrypts a plain password using bcrypt and returns the hashed string
func (h *hasher) Hash(password string) (string, error) {
	hashBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", ErrHashFailed
	}

	return string(hashBytes), nil
}

// Compare checks if the bcrypt hash matches the plain password
func (h *hasher) Compare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
