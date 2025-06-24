package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
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

// ValidateRatings validates and normalizes ratings
func ValidateRatings(ratings *models.Ratings) {
	// Ensure ratings are within valid ranges
	if ratings.TMDB < 0 || ratings.TMDB > 10 {
		ratings.TMDB = 0
	}
	if ratings.OMDB < 0 || ratings.OMDB > 10 {
		ratings.OMDB = 0
	}
	if ratings.RottenTomatoes < 0 || ratings.RottenTomatoes > 100 {
		ratings.RottenTomatoes = 0
	}
	if ratings.IMDB < 0 || ratings.IMDB > 10 {
		ratings.IMDB = 0
	}
	if ratings.Metacritic < 0 || ratings.Metacritic > 100 {
		ratings.Metacritic = 0
	}
}

// CalculateAverageRating calculates the average rating from multiple sources
func CalculateAverageRating(ratings models.Ratings) float64 {
	var sum float64
	var count int

	if ratings.TMDB > 0 {
		sum += ratings.TMDB
		count++
	}
	if ratings.OMDB > 0 {
		sum += ratings.OMDB
		count++
	}
	if ratings.IMDB > 0 {
		sum += ratings.IMDB
		count++
	}
	if ratings.RottenTomatoes > 0 {
		// Convert Rotten Tomatoes to 10-point scale
		sum += (ratings.RottenTomatoes / 10.0)
		count++
	}
	if ratings.Metacritic > 0 {
		// Convert Metacritic to 10-point scale
		sum += (ratings.Metacritic / 10.0)
		count++
	}

	if count == 0 {
		return 0
	}

	return math.Round((sum/float64(count))*10) / 10
}

// ParseYear extracts year from date string
func ParseYear(dateStr string) int {
	if dateStr == "" {
		return 0
	}

	// Try to parse as YYYY-MM-DD
	parts := strings.Split(dateStr, "-")
	if len(parts) > 0 {
		if year, err := strconv.Atoi(parts[0]); err == nil {
			return year
		}
	}

	// Try to parse as YYYY
	if year, err := strconv.Atoi(dateStr); err == nil && year > 1900 && year <= time.Now().Year()+10 {
		return year
	}

	return 0
}

// SanitizeString removes special characters and normalizes string
func SanitizeString(s string) string {
	// Remove special characters but keep spaces and basic punctuation
	reg := regexp.MustCompile(`[^\w\s\-.,!?()]`)
	sanitized := reg.ReplaceAllString(s, "")

	// Normalize whitespace
	sanitized = regexp.MustCompile(`\s+`).ReplaceAllString(sanitized, " ")

	return strings.TrimSpace(sanitized)
}

// TruncateString truncates a string to specified length
func TruncateString(s string, maxLength int) string {
	if len(s) <= maxLength {
		return s
	}
	return s[:maxLength-3] + "..."
}

// FormatDuration formats duration in minutes to human readable format
func FormatDuration(minutes int) string {
	if minutes <= 0 {
		return "Unknown"
	}

	hours := minutes / 60
	mins := minutes % 60

	if hours > 0 {
		if mins > 0 {
			return fmt.Sprintf("%dh %dm", hours, mins)
		}
		return fmt.Sprintf("%dh", hours)
	}
	return fmt.Sprintf("%dm", mins)
}

// CalculateSimilarity calculates similarity between two movies based on genres
func CalculateSimilarity(movie1, movie2 models.Movie) float64 {
	if len(movie1.GenreIDs) == 0 || len(movie2.GenreIDs) == 0 {
		return 0
	}

	// Create sets of genre IDs
	genres1 := make(map[int]bool)
	for _, id := range movie1.GenreIDs {
		genres1[id] = true
	}

	genres2 := make(map[int]bool)
	for _, id := range movie2.GenreIDs {
		genres2[id] = true
	}

	// Calculate intersection and union
	intersection := 0
	for id := range genres1 {
		if genres2[id] {
			intersection++
		}
	}

	union := len(genres1) + len(genres2) - intersection

	if union == 0 {
		return 0
	}

	return float64(intersection) / float64(union)
}

// PaginateResults handles pagination for search results
func PaginateResults(results []models.Movie, page, perPage int) ([]models.Movie, models.Meta) {
	totalResults := len(results)
	totalPages := int(math.Ceil(float64(totalResults) / float64(perPage)))

	// Validate page number
	if page < 1 {
		page = 1
	}
	if page > totalPages && totalPages > 0 {
		page = totalPages
	}

	// Calculate start and end indices
	start := (page - 1) * perPage
	end := start + perPage

	if start >= totalResults {
		return []models.Movie{}, models.Meta{
			Page:         page,
			PerPage:      perPage,
			TotalPages:   totalPages,
			TotalResults: totalResults,
			HasNext:      false,
			HasPrev:      page > 1,
		}
	}

	if end > totalResults {
		end = totalResults
	}

	paginatedResults := results[start:end]

	meta := models.Meta{
		Page:         page,
		PerPage:      perPage,
		TotalPages:   totalPages,
		TotalResults: totalResults,
		HasNext:      page < totalPages,
		HasPrev:      page > 1,
	}

	return paginatedResults, meta
}

// MarshalJSON safely marshals data to JSON
func MarshalJSON(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

// UnmarshalJSON safely unmarshals JSON data
func UnmarshalJSON(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
