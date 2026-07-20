package service

import (
	"context"
	"errors"
	"time"

	"github.com/TNJKL/bookmark-management/internal/repository"
	"github.com/TNJKL/bookmark-management/pkg/utils"
)

// ShortenURL defines the service operations for creating and resolving
// shortened URLs.
//
//go:generate mockery --name ShortenURL --filename shortenurl.go
type ShortenURL interface {
	// CreateShortenLink generates a unique short code and maps it to the original URL in storage.
	CreateShortenLink(ctx context.Context, url string, exp int64) (string, error)
	// GetLinkFromCode retrieves the original URL associated with the given short code.
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

// mock GenPass & Storage de test
// shortenURL is the default implementation of the ShortenURL interface.
type shortenURL struct {
	storage repository.URLStorage
	keyGen  utils.KeyGenerator
}

// NewShortenURL creates a new ShortenURL service  instance.
func NewShortenUrl(storage repository.URLStorage, keyGen utils.KeyGenerator) ShortenURL {
	return &shortenURL{
		storage: storage,
		keyGen:  keyGen,
	}
}

const linkKeyLength = 7

// Chỗ này em làm theo cách của em trước,
// sau buổi 3 anh chữa bài --> Em thấy cách của em chưa clean lắm,
// sau khi anh review xong em sẽ sửa lại theo cách của anh như trong video bài giảng
// CreateShortenLink generates a unique short code and maps it to the original URL in storage
func (s *shortenURL) CreateShortenLink(ctx context.Context, url string, exp int64) (string, error) {
	//gen code

	key := s.keyGen.GenerateKey(linkKeyLength)

	res, err := s.storage.GetURL(ctx, key)

	if err != nil && !errors.Is(err, repository.ErrorCodeNotFound) {
		return "", err
	}

	if res != "" {
		return s.CreateShortenLink(ctx, url, exp)
	}

	//goi repo de store URL
	err = s.storage.StoreURL(ctx, key, url, time.Duration(exp)*time.Second)
	if err != nil {
		return "", err
	}
	//return code
	return key, nil
}

// GetLinkFromCode retrieves the original URL associated with the given short code.
func (s *shortenURL) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return s.storage.GetURL(ctx, code)
}
