package service

import (
	"crypto/rand"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// service.GenPass // service.GenPassService

//name = interface name
//go:generate mockery --name GenPass --filename genpass.go

// Interface: Service có thể làm gì
type GenPass interface {
	GeneratePassword(length int) (string, error)
}

// Struct: Implement thực tế
type genPassService struct {
}

// Constructor
func NewGenPass() GenPass {
	return &genPassService{}
}

// Business logic thuần
func (s *genPassService) GeneratePassword(length int) (string, error) {
	password := make([]byte, length)

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password[i] = charset[n.Int64()]
	}
	return string(password), nil
}
