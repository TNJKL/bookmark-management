package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/TNJKL/bookmark-management/internal/repository"
	"github.com/redis/go-redis/v9"
)

const codeLength = 7
const maxRetries = 5

// ShortenURL defines the service operations for creating and resolving
// shortened URLs.
//
//go:generate mockery --name ShortenURL --filename shortenurl.go
type ShortenURL interface {
	// CreateShortenLink generates a unique short code and maps it to the original URL in storage.
	CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error)
	// GetLinkFromCode retrieves the original URL associated with the given short code.
	GetLinkFromCode(ctx context.Context, code string) (string, error)
}

// mock GenPass & Storage de test
// shortenURL is the default implementation of the ShortenURL interface.
type shortenURL struct {
	storage repository.URLStorage
	codeGen GenPass
}

// NewShortenURL creates a new ShortenURL service  instance.
func NewShortenUrl(storage repository.URLStorage, codeGen GenPass) ShortenURL {
	return &shortenURL{
		storage: storage,
		codeGen: codeGen,
	}
}

// Chỗ này em làm theo cách của em trước,
// sau buổi 3 anh chữa bài --> Em thấy cách của em chưa clean lắm,
// sau khi anh review xong em sẽ sửa lại theo cách của anh như trong video bài giảng
// CreateShortenLink generates a unique short code and maps it to the original URL in storage
func (s *shortenURL) CreateShortenLink(ctx context.Context, url string, exp time.Duration) (string, error) {
	//gen code

	var code string
	for i := 0; i < maxRetries; i++ {
		generatedCode, err := s.codeGen.GeneratePassword(codeLength)
		if err != nil {
			return "", err
		}
		_, err = s.storage.GetURL(ctx, generatedCode)
		//truong hop 1 : code chua ton tai trong Redis (redis.nil)
		if errors.Is(err, redis.Nil) {
			code = generatedCode
			break
		}
		//truong hop 2 : code da ton tai trong Redis --> retry generate code
		if err == nil {
			continue //nhay sang vong lap tiep theo --> sinh code moi'
		}
		return "", err
	}
	//sau vòng lặp , nếu code vẫn là "" nghĩa là đã retry hết maxRetries mà lần nào cũng trùng
	if code == "" {
		return "", fmt.Errorf("failed to generate unique code after %d retries", maxRetries)
	}

	//goi repo de store URL
	err := s.storage.StoreURL(ctx, code, url, exp)
	if err != nil {
		return "", err
	}
	//return code
	return code, nil
}

// GetLinkFromCode retrieves the original URL associated with the given short code.
func (s *shortenURL) GetLinkFromCode(ctx context.Context, code string) (string, error) {
	return "", nil
}
