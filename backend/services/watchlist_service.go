package services

import (
	"context"
	"fmt"
	"math"
	"sort"
	"time"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/utils"
)

// WatchlistService handles watchlist operations and recommendations
type WatchlistService struct {
	tmdbService *TMDBService
	cache       *utils.Cache
	// In a real application, you would have a database here
	watchlists map[string]*models.Watchlist // userID -> watchlist
}

// NewWatchlistService creates a new watchlist service instance
func NewWatchlistService(tmdbService *TMDBService) *WatchlistService {
	return &WatchlistService{
		tmdbService: tmdbService,
		cache:       utils.NewCache(),
		watchlists:  make(map[string]*models.Watchlist),
	}
}

// CreateWatchlist creates a new watchlist for a user
func (s *WatchlistService) CreateWatchlist(ctx context.Context, userID string, request models.WatchlistCreateRequest) (*models.Watchlist, error) {
	// Check if user already has a watchlist
	if _, exists := s.watchlists[userID]; exists {
		return nil, fmt.Errorf("user already has a watchlist")
	}

	watchlist := &models.Watchlist{
		ID:          len(s.watchlists) + 1, // Simple ID generation
		UserID:      userID,
		Name:        request.Name,
		Description: request.Description,
		IsPublic:    request.IsPublic,
		Items:       []models.WatchlistItem{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	s.watchlists[userID] = watchlist

	return watchlist, nil
}

// GetWatchlist retrieves a user's watchlist
func (s *WatchlistService) GetWatchlist(ctx context.Context, userID string) (*models.Watchlist, error) {
	watchlist, exists := s.watchlists[userID]
	if !exists {
		return nil, fmt.Errorf("watchlist not found")
	}

	// Populate movie details for each item
	for i := range watchlist.Items {
		movie, err := s.tmdbService.GetMovieDetails(ctx, watchlist.Items[i].MovieID)
		if err == nil {
			watchlist.Items[i].Movie = *movie
		}
	}

	return watchlist, nil
}

// AddToWatchlist adds a movie to a user's watchlist
func (s *WatchlistService) AddToWatchlist(ctx context.Context, userID string, request models.WatchlistItemAddRequest) (*models.WatchlistItem, error) {
	watchlist, exists := s.watchlists[userID]
	if !exists {
		return nil, fmt.Errorf("watchlist not found")
	}

	// Check if movie already exists in watchlist
	for _, item := range watchlist.Items {
		if item.MovieID == request.MovieID {
			return nil, fmt.Errorf("movie already in watchlist")
		}
	}

	// Create new watchlist item
	item := models.WatchlistItem{
		ID:          len(watchlist.Items) + 1,
		WatchlistID: watchlist.ID,
		MovieID:     request.MovieID,
		Status:      request.Status,
		Rating:      request.Rating,
		Notes:       request.Notes,
		AddedAt:     time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Get movie details
	movie, err := s.tmdbService.GetMovieDetails(ctx, request.MovieID)
	if err == nil {
		item.Movie = *movie
	}

	watchlist.Items = append(watchlist.Items, item)
	watchlist.UpdatedAt = time.Now()

	return &item, nil
}

// UpdateWatchlistItem updates a watchlist item
func (s *WatchlistService) UpdateWatchlistItem(ctx context.Context, userID string, itemID int, request models.WatchlistItemUpdateRequest) (*models.WatchlistItem, error) {
	watchlist, exists := s.watchlists[userID]
	if !exists {
		return nil, fmt.Errorf("watchlist not found")
	}

	// Find and update the item
	for i := range watchlist.Items {
		if watchlist.Items[i].ID == itemID {
			watchlist.Items[i].Status = request.Status
			watchlist.Items[i].Rating = request.Rating
			watchlist.Items[i].Notes = request.Notes
			watchlist.Items[i].UpdatedAt = time.Now()
			watchlist.UpdatedAt = time.Now()

			return &watchlist.Items[i], nil
		}
	}

	return nil, fmt.Errorf("watchlist item not found")
}

// RemoveFromWatchlist removes a movie from a user's watchlist
func (s *WatchlistService) RemoveFromWatchlist(ctx context.Context, userID string, itemID int) error {
	watchlist, exists := s.watchlists[userID]
	if !exists {
		return fmt.Errorf("watchlist not found")
	}

	// Find and remove the item
	for i, item := range watchlist.Items {
		if item.ID == itemID {
			watchlist.Items = append(watchlist.Items[:i], watchlist.Items[i+1:]...)
			watchlist.UpdatedAt = time.Now()
			return nil
		}
	}

	return fmt.Errorf("watchlist item not found")
}

// GetWatchlistStats returns statistics for a user's watchlist
func (s *WatchlistService) GetWatchlistStats(ctx context.Context, userID string) (*models.WatchlistStats, error) {
	watchlist, exists := s.watchlists[userID]
	if !exists {
		return nil, fmt.Errorf("watchlist not found")
	}

	stats := &models.WatchlistStats{
		TotalItems:     len(watchlist.Items),
		CompletedItems: 0,
		WatchingItems:  0,
		ToWatchItems:   0,
		DroppedItems:   0,
		AverageRating:  0,
		TotalHours:     0,
	}

	var totalRating float64
	var ratingCount int

	for _, item := range watchlist.Items {
		switch item.Status {
		case "completed":
			stats.CompletedItems++
		case "watching":
			stats.WatchingItems++
		case "to_watch":
			stats.ToWatchItems++
		case "dropped":
			stats.DroppedItems++
		}

		if item.Rating > 0 {
			totalRating += item.Rating
			ratingCount++
		}

		// Add runtime to total hours
		if item.Movie.Runtime > 0 {
			stats.TotalHours += item.Movie.Runtime
		}
	}

	if ratingCount > 0 {
		stats.AverageRating = totalRating / float64(ratingCount)
	}

	// Convert minutes to hours
	stats.TotalHours = stats.TotalHours / 60

	return stats, nil
}

// GetRecommendations generates movie recommendations based on user's watchlist
func (s *WatchlistService) GetRecommendations(ctx context.Context, userID string, limit int) ([]models.WatchlistRecommendation, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("recommendations", userID, limit)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if recommendations, ok := cached.([]models.WatchlistRecommendation); ok {
			return recommendations, nil
		}
	}

	watchlist, exists := s.watchlists[userID]
	if !exists {
		return nil, fmt.Errorf("watchlist not found")
	}

	if len(watchlist.Items) == 0 {
		return []models.WatchlistRecommendation{}, nil
	}

	// Analyze user preferences
	preferences := s.analyzeUserPreferences(watchlist)

	// Get trending movies as candidate recommendations
	trending, err := s.tmdbService.GetTrendingMedia(ctx, "week", 1, "all")
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}

	// Score and rank recommendations
	recommendations := s.scoreRecommendations(trending.Results, preferences, limit)

	// Cache the result
	s.cache.Set(cacheKey, recommendations, 30*time.Minute)

	return recommendations, nil
}

// analyzeUserPreferences analyzes user's watchlist to determine preferences
func (s *WatchlistService) analyzeUserPreferences(watchlist *models.Watchlist) map[string]float64 {
	preferences := make(map[string]float64)
	genreCounts := make(map[int]int)
	ratingSum := 0.0
	ratingCount := 0

	// Count genres and calculate average rating
	for _, item := range watchlist.Items {
		for _, genreID := range item.Movie.GenreIDs {
			genreCounts[genreID]++
		}

		if item.Rating > 0 {
			ratingSum += item.Rating
			ratingCount++
		}
	}

	// Calculate genre preferences
	totalItems := len(watchlist.Items)
	if totalItems > 0 {
		for genreID, count := range genreCounts {
			preferences[fmt.Sprintf("genre_%d", genreID)] = float64(count) / float64(totalItems)
		}
	}

	// Calculate average rating preference
	if ratingCount > 0 {
		preferences["avg_rating"] = ratingSum / float64(ratingCount)
	}

	// Calculate year preferences
	yearCounts := make(map[int]int)
	for _, item := range watchlist.Items {
		year := utils.ParseYear(item.Movie.ReleaseDate)
		if year > 0 {
			yearCounts[year]++
		}
	}

	for year, count := range yearCounts {
		preferences[fmt.Sprintf("year_%d", year)] = float64(count) / float64(totalItems)
	}

	return preferences
}

// scoreRecommendations scores and ranks movie recommendations
func (s *WatchlistService) scoreRecommendations(candidates []models.Movie, preferences map[string]float64, limit int) []models.WatchlistRecommendation {
	var recommendations []models.WatchlistRecommendation

	for _, movie := range candidates {
		score := 0.0
		genreMatch := 0.0
		ratingMatch := 0.0

		// Calculate genre match
		for _, genreID := range movie.GenreIDs {
			if pref, exists := preferences[fmt.Sprintf("genre_%d", genreID)]; exists {
				genreMatch += pref
			}
		}

		// Normalize genre match
		if len(movie.GenreIDs) > 0 {
			genreMatch = genreMatch / float64(len(movie.GenreIDs))
		}

		// Calculate rating match
		if avgRating, exists := preferences["avg_rating"]; exists && movie.VoteAverage > 0 {
			ratingDiff := math.Abs(movie.VoteAverage - avgRating)
			ratingMatch = 1.0 - (ratingDiff / 10.0) // Higher score for closer ratings
		}

		// Calculate year match
		year := utils.ParseYear(movie.ReleaseDate)
		yearMatch := 0.0
		if year > 0 {
			if yearPref, exists := preferences[fmt.Sprintf("year_%d", year)]; exists {
				yearMatch = yearPref
			}
		}

		// Calculate overall score
		score = (genreMatch * 0.5) + (ratingMatch * 0.3) + (yearMatch * 0.2)

		// Add popularity bonus
		if movie.Popularity > 0 {
			popularityBonus := math.Min(movie.Popularity/100.0, 0.1) // Max 10% bonus
			score += popularityBonus
		}

		// Determine reason for recommendation
		reason := "Based on your watchlist preferences"
		if genreMatch > 0.5 {
			reason = "Similar genres to your favorites"
		} else if ratingMatch > 0.7 {
			reason = "Matches your rating preferences"
		} else if yearMatch > 0.3 {
			reason = "From your preferred time period"
		}

		recommendations = append(recommendations, models.WatchlistRecommendation{
			Movie:       movie,
			Score:       score,
			Reason:      reason,
			GenreMatch:  genreMatch,
			RatingMatch: ratingMatch,
		})
	}

	// Sort by score (highest first)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// Return top recommendations
	if len(recommendations) > limit {
		recommendations = recommendations[:limit]
	}

	return recommendations
}

// GetSimilarMovies finds movies similar to a given movie
func (s *WatchlistService) GetSimilarMovies(ctx context.Context, movieID int, limit int) ([]models.Movie, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("similar_movies", movieID, limit)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if movies, ok := cached.([]models.Movie); ok {
			return movies, nil
		}
	}

	// Get the source movie
	sourceMovie, err := s.tmdbService.GetMovieDetails(ctx, movieID)
	if err != nil {
		return nil, fmt.Errorf("failed to get source movie: %w", err)
	}

	// Get trending movies as candidates
	trending, err := s.tmdbService.GetTrendingMedia(ctx, "week", 1, "all")
	if err != nil {
		return nil, fmt.Errorf("failed to get trending movies: %w", err)
	}

	// Calculate similarities
	var similarMovies []models.Movie
	for _, movie := range trending.Results {
		if movie.ID != movieID {
			similarity := utils.CalculateSimilarity(*sourceMovie, movie)
			if similarity > 0.1 { // Only include movies with some similarity
				similarMovies = append(similarMovies, movie)
			}
		}
	}

	// Sort by similarity
	sort.Slice(similarMovies, func(i, j int) bool {
		sim1 := utils.CalculateSimilarity(*sourceMovie, similarMovies[i])
		sim2 := utils.CalculateSimilarity(*sourceMovie, similarMovies[j])
		return sim1 > sim2
	})

	// Limit results
	if len(similarMovies) > limit {
		similarMovies = similarMovies[:limit]
	}

	// Cache the result
	s.cache.Set(cacheKey, similarMovies, 30*time.Minute)

	return similarMovies, nil
}

// Close closes the service and cleans up resources
func (s *WatchlistService) Close() {
	// Clean up cache
	s.cache.Cleanup()
}
