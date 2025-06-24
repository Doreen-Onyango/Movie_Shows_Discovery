package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
)

// Cache represents a simple in-memory cache
type Cache struct {
	data map[string]CacheItem
}

// CacheItem represents a cached item with expiration
type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{
		data: make(map[string]CacheItem),
	}
}

// Set adds an item to the cache with expiration
func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.data[key] = CacheItem{
		Data:      value,
		ExpiresAt: time.Now().Add(ttl),
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	item, exists := c.data[key]
	if !exists {
		return nil, false
	}

	if time.Now().After(item.ExpiresAt) {
		delete(c.data, key)
		return nil, false
	}

	return item.Data, true
}

// Delete removes an item from the cache
func (c *Cache) Delete(key string) {
	delete(c.data, key)
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.data = make(map[string]CacheItem)
}

// Cleanup removes expired items from the cache
func (c *Cache) Cleanup() {
	now := time.Now()
	for key, item := range c.data {
		if now.After(item.ExpiresAt) {
			delete(c.data, key)
		}
	}
}

// GenerateCacheKey generates a cache key from parameters
func GenerateCacheKey(prefix string, params ...interface{}) string {
	key := prefix
	for _, param := range params {
		key += ":" + fmt.Sprintf("%v", param)
	}

	// Create MD5 hash for consistent key length
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// ValidateMovieData validates movie data and provides fallbacks
func ValidateMovieData(movie *models.Movie) {
	// Ensure required fields have fallback values
	if movie.Title == "" {
		movie.Title = "Unknown Title"
	}
	if movie.Overview == "" {
		movie.Overview = "No overview available"
	}
	if movie.PosterPath == "" {
		movie.PosterPath = "/placeholder-poster.jpg"
	}
	if movie.BackdropPath == "" {
		movie.BackdropPath = "/placeholder-backdrop.jpg"
	}
	if movie.ReleaseDate == "" {
		movie.ReleaseDate = "Unknown"
	}
	if movie.Status == "" {
		movie.Status = "Unknown"
	}
	if movie.Tagline == "" {
		movie.Tagline = "No tagline available"
	}

	// Ensure ratings are within valid range
	if movie.VoteAverage < 0 || movie.VoteAverage > 10 {
		movie.VoteAverage = 0
	}
	if movie.VoteCount < 0 {
		movie.VoteCount = 0
	}
	if movie.Popularity < 0 {
		movie.Popularity = 0
	}
	if movie.Runtime < 0 {
		movie.Runtime = 0
	}
}
