package bookmark

import (
	"context"

	"github.com/TNJKL/bookmark-management/internal/app/model"
)

const codeLength = 8

// CreateBookmark generates a unique short code and creates a new bookmark for a user
func (b *bookmarkService) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	//create bookmark model
	code := b.keyGen.GenerateKey(codeLength)

	bookmark := &model.Bookmark{
		Description: description,
		URL:         url,
		UserID:      userID,
		Code:        code,
	}
	//call repo
	res, err := b.repo.CreateBookmark(ctx, bookmark)
	if err != nil {
		return nil, err
	}

	//return bookmark
	return res, nil
}
