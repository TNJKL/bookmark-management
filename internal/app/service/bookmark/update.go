package bookmark

import (
	"context"
)

// UpdateBookmark updates description and url of an existing bookmark for a user
func (s *bookmarkService) UpdateBookmark(ctx context.Context, id, userID, description, url string) error {
	return s.repo.UpdateBookmark(ctx, id, userID, description, url)
}
