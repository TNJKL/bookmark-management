package bookmark

import (
	"context"
)

// DeleteBookmark deletes a bookmark belonging to a specific user by id
func (s *bookmarkService) DeleteBookmark(ctx context.Context, id, userID string) error {
	return s.repo.DeleteBookmark(ctx, id, userID)
}
