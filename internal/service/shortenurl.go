package service

import (
	"context"
	"time"

	"github.com/TNJKL/bookmark-management/internal/repository"
)

const codeLength = 7

// ShortenURL defines the service operations for creating and resolving
// shortened URLs.
type ShortenURL interface {
	CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error)
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

// mock GenPass & Storage de test
// shortenURL struct implements URLStorage repo & genPass service
type shortenURL struct {
	storage repository.URLStorage
	codeGen GenPass
}

// NewShortenURL creates a new ShortenURL service
func NewShortenUrl(storage repository.URLStorage, codeGen GenPass) ShortenURL {
	return &shortenURL{
		storage: storage,
		codeGen: codeGen,
	}
}

func (s *shortenURL) CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error) {
	//gen code
	code, err := s.codeGen.GeneratePassword(codeLength)
	if err != nil {
		return "", err
	}
	//goi repo de store URL
	err = s.storage.StoreURL(ctx, code, url, exp)
	if err != nil {
		return "", err
	}
	//return code
	return code, nil
}

func (s *shortenURL) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return "", nil
}
