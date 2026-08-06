package bookmark

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/TNJKL/bookmark-management/internal/app/model"
	"github.com/TNJKL/bookmark-management/internal/app/repository/cache"
	"github.com/rs/zerolog/log"
)

const (
	getBookmarksCacheGroupKeyFormat = "get_bookmarks_%s"
	getBookmarksCacheKeyFormat      = "%d_%d"
	getBookmarksCacheExp            = 24 * time.Hour
)

// bookmarkServiceWithCache wraps Service to provide transparent Redis caching for bookmark operations.
type bookmarkServiceWithCache struct {
	s Service
	c cache.DB
}

// NewServiceWithCache creates a new bookmark Service decorator with caching support.
func NewServiceWithCache(s Service, c cache.DB) Service {
	return &bookmarkServiceWithCache{
		s,
		c,
	}
}

// GetBookmarks retrieves bookmarks using the cache-aside pattern, checking Redis first before falling back to the database.
func (s *bookmarkServiceWithCache) GetBookmarks(ctx context.Context, userID string, page, limit int) (*GetBookmarksResult, error) {
	//create cache key
	cacheGroupKey := s.getCacheGroupKey(userID)
	cacheKey := fmt.Sprintf(getBookmarksCacheKeyFormat, page, limit)

	//get cache data
	cacheData, err := s.c.GetCacheData(ctx, cacheGroupKey, cacheKey)
	if err == nil && len(cacheData) > 0 {
		result := &GetBookmarksResult{}
		err := json.Unmarshal(cacheData, result)
		if err != nil {
			err = s.c.DeleteCache(ctx, cacheGroupKey)
			if err != nil {
				log.Err(err).Str("key", cacheGroupKey).Msg("failed to delete cache")
			}
		} else {
			return result, nil
		}
	}
	//if not , call service
	result, err := s.s.GetBookmarks(ctx, userID, page, limit)
	if err != nil {
		return nil, err
	}

	//save cache
	resultBytes, err := json.Marshal(result)
	if err == nil {
		cacheErr := s.c.SetCacheData(ctx, cacheGroupKey, cacheKey, resultBytes, getBookmarksCacheExp)
		if cacheErr != nil {
			log.Err(err).Str("key", cacheGroupKey).Msg("failed to set cache")
		}
	}
	return result, nil

}

// CreateBookmark creates a new bookmark and invalidates the user's cached bookmark list.
func (s *bookmarkServiceWithCache) CreateBookmark(ctx context.Context, description, url, userID string) (*model.Bookmark, error) {
	err := s.c.DeleteCache(ctx, s.getCacheGroupKey(userID))
	if err != nil {
		return nil, err
	}

	return s.s.CreateBookmark(ctx, description, url, userID)
}

// UpdateBookmark updates an existing bookmark and invalidates the user's cached bookmark list.
func (s *bookmarkServiceWithCache) UpdateBookmark(ctx context.Context, id, userID, description, url string) error {
	err := s.c.DeleteCache(ctx, s.getCacheGroupKey(userID))
	if err != nil {
		return err
	}
	return s.s.UpdateBookmark(ctx, id, userID, description, url)

}

// DeleteBookmark removes a bookmark and invalidates the user's cached bookmark list.
func (s *bookmarkServiceWithCache) DeleteBookmark(ctx context.Context, id, userID string) error {
	err := s.c.DeleteCache(ctx, s.getCacheGroupKey(userID))
	if err != nil {
		return err
	}

	return s.s.DeleteBookmark(ctx, id, userID)
}

// getCacheGroupKey formats the Redis cache group key for a given user ID.
func (s *bookmarkServiceWithCache) getCacheGroupKey(userID string) string {
	return fmt.Sprintf(getBookmarksCacheGroupKeyFormat, userID)
}
