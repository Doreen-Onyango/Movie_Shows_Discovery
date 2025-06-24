package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/middleware"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/services"
)

// TrendingController handles trending content requests
type TrendingController struct {
	tmdbService *services.TMDBService
	logger      *middleware.Logger
}

// NewTrendingController creates a new trending controller
func NewTrendingController(tmdbService *services.TMDBService, logger *middleware.Logger) *TrendingController {
	return &TrendingController{
		tmdbService: tmdbService,
		logger:      logger,
	}
}

// GetTrending handles trending content requests
func (c *TrendingController) GetTrending(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		timeframe = "day"
	}

	// Validate timeframe
	if timeframe != "day" && timeframe != "week" {
		http.Error(w, "Invalid timeframe. Use 'day' or 'week'", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	// Get trending movies
	result, err := c.tmdbService.GetTrendingMovies(r.Context(), timeframe, page)
	if err != nil {
		c.logger.LogError(err, "GetTrending", r)
		http.Error(w, "Failed to get trending content", http.StatusInternalServerError)
		return
	}

	// Create response
	meta := models.Meta{
		Page:         result.Page,
		PerPage:      20,
		TotalPages:   result.TotalPages,
		TotalResults: result.TotalResults,
		HasNext:      result.Page < result.TotalPages,
		HasPrev:      result.Page > 1,
	}

	response := models.NewTrendingResponse(result.Results, timeframe, meta)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTrendingByGenre handles trending content by genre
func (c *TrendingController) GetTrendingByGenre(w http.ResponseWriter, r *http.Request) {
	// Get genre ID from URL
	genreIDStr := r.URL.Query().Get("genre_id")
	genreID, err := strconv.Atoi(genreIDStr)
	if err != nil {
		http.Error(w, "Invalid genre ID", http.StatusBadRequest)
		return
	}

	// Get query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	sortBy := r.URL.Query().Get("sort_by")
	if sortBy == "" {
		sortBy = "popularity.desc"
	}

	// Get trending movies by genre
	result, err := c.tmdbService.GetMoviesByGenre(r.Context(), genreID, page, sortBy)
	if err != nil {
		c.logger.LogError(err, "GetTrendingByGenre", r)
		http.Error(w, "Failed to get trending content by genre", http.StatusInternalServerError)
		return
	}

	// Create response
	meta := models.Meta{
		Page:         result.Page,
		PerPage:      20,
		TotalPages:   result.TotalPages,
		TotalResults: result.TotalResults,
		HasNext:      result.Page < result.TotalPages,
		HasPrev:      result.Page > 1,
	}

	response := models.NewPaginatedResponse(result.Results, meta, "Trending content by genre retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTrendingStats handles trending statistics requests
func (c *TrendingController) GetTrendingStats(w http.ResponseWriter, r *http.Request) {
	// Get trending movies for both day and week
	dayTrending, err := c.tmdbService.GetTrendingMovies(r.Context(), "day", 1)
	if err != nil {
		c.logger.LogError(err, "GetTrendingStats - day", r)
		http.Error(w, "Failed to get daily trending stats", http.StatusInternalServerError)
		return
	}

	weekTrending, err := c.tmdbService.GetTrendingMovies(r.Context(), "week", 1)
	if err != nil {
		c.logger.LogError(err, "GetTrendingStats - week", r)
		http.Error(w, "Failed to get weekly trending stats", http.StatusInternalServerError)
		return
	}

	// Create stats response
	stats := map[string]interface{}{
		"daily_trending_count":  len(dayTrending.Results),
		"weekly_trending_count": len(weekTrending.Results),
		"daily_top_movies":      dayTrending.Results[:min(5, len(dayTrending.Results))],
		"weekly_top_movies":     weekTrending.Results[:min(5, len(weekTrending.Results))],
	}

	// Create response
	response := models.NewSuccessResponse(stats, "Trending stats retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTrendingGenres handles trending genres requests
func (c *TrendingController) GetTrendingGenres(w http.ResponseWriter, r *http.Request) {
	// Get all genres
	genres, err := c.tmdbService.GetGenres(r.Context())
	if err != nil {
		c.logger.LogError(err, "GetTrendingGenres", r)
		http.Error(w, "Failed to get trending genres", http.StatusInternalServerError)
		return
	}

	// Get trending movies to analyze popular genres
	trending, err := c.tmdbService.GetTrendingMovies(r.Context(), "week", 1)
	if err != nil {
		c.logger.LogError(err, "GetTrendingGenres - trending", r)
		http.Error(w, "Failed to get trending genres", http.StatusInternalServerError)
		return
	}

	// Count genre occurrences in trending movies
	genreCounts := make(map[int]int)
	for _, movie := range trending.Results {
		for _, genreID := range movie.GenreIDs {
			genreCounts[genreID]++
		}
	}

	// Create trending genres response
	type TrendingGenre struct {
		Genre   models.Genre `json:"genre"`
		Count   int          `json:"count"`
		Popular bool         `json:"popular"`
	}

	var trendingGenres []TrendingGenre
	for _, genre := range genres {
		count := genreCounts[genre.ID]
		trendingGenres = append(trendingGenres, TrendingGenre{
			Genre:   genre,
			Count:   count,
			Popular: count >= 3, // Consider popular if appears in 3+ trending movies
		})
	}

	// Create response
	response := models.NewSuccessResponse(trendingGenres, "Trending genres retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Helper function to get minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
