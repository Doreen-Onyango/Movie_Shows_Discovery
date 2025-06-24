package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/middleware"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/services"
	"github.com/gorilla/mux"
)

// MovieController handles movie-related HTTP requests
type MovieController struct {
	tmdbService      *services.TMDBService
	omdbService      *services.OMDBService
	logger           *middleware.Logger
	watchlistService *services.WatchlistService
}

// NewMovieController creates a new movie controller
func NewMovieController(tmdbService *services.TMDBService, omdbService *services.OMDBService, logger *middleware.Logger, watchlistService *services.WatchlistService) *MovieController {
	return &MovieController{
		tmdbService:      tmdbService,
		omdbService:      omdbService,
		logger:           logger,
		watchlistService: watchlistService,
	}
}

// SearchMovies handles movie search requests
func (c *MovieController) SearchMovies(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage <= 0 {
		perPage = 10 // default to 10
	}

	includeAdult := r.URL.Query().Get("include_adult") == "true"
	mediaType := r.URL.Query().Get("type")
	if mediaType == "" {
		mediaType = "all"
	}

	// Search movies and/or TV shows
	result, err := c.tmdbService.SearchMedia(r.Context(), query, page, perPage, mediaType, includeAdult)
	if err != nil {
		c.logger.LogError(err, "SearchMovies", r)
		http.Error(w, "Failed to search movies/TV shows", http.StatusInternalServerError)
		return
	}

	// Create response
	meta := models.Meta{
		Page:         result.Page,
		PerPage:      perPage,
		TotalPages:   result.TotalPages,
		TotalResults: result.TotalResults,
		HasNext:      result.Page < result.TotalPages,
		HasPrev:      result.Page > 1,
	}

	response := models.NewSearchResponse(result.Results, query, meta)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMovieDetails handles movie details requests
func (c *MovieController) GetMovieDetails(w http.ResponseWriter, r *http.Request) {
	// Get movie ID from URL
	vars := mux.Vars(r)
	movieIDStr := vars["id"]

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Get movie details from TMDB
	movie, err := c.tmdbService.GetMovieDetails(r.Context(), movieID)
	if err != nil {
		c.logger.LogError(err, "GetMovieDetails", r)
		http.Error(w, "Failed to get movie details", http.StatusInternalServerError)
		return
	}

	// Enrich with OMDB data
	err = c.omdbService.EnrichMovieWithOMDBData(r.Context(), movie)
	if err != nil {
		// Log error but don't fail the request
		c.logger.LogError(err, "EnrichMovieWithOMDBData", r)
	}

	// Create response
	response := models.NewSuccessResponse(movie, "Movie details retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetTrendingMovies handles trending movies requests
func (c *MovieController) GetTrendingMovies(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	timeframe := r.URL.Query().Get("timeframe")
	if timeframe == "" {
		timeframe = "day"
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1
	}

	mediaType := r.URL.Query().Get("type")
	if mediaType == "" {
		mediaType = "all"
	}

	// Get trending movies and/or TV shows
	result, err := c.tmdbService.GetTrendingMedia(r.Context(), timeframe, page, mediaType)
	if err != nil {
		c.logger.LogError(err, "GetTrendingMovies", r)
		http.Error(w, "Failed to get trending movies", http.StatusInternalServerError)
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

// GetMoviesByGenre handles genre-based movie requests
func (c *MovieController) GetMoviesByGenre(w http.ResponseWriter, r *http.Request) {
	// Get genre ID from URL
	vars := mux.Vars(r)
	genreIDStr := vars["genreId"]

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

	// Get movies by genre
	result, err := c.tmdbService.GetMoviesByGenre(r.Context(), genreID, page, sortBy)
	if err != nil {
		c.logger.LogError(err, "GetMoviesByGenre", r)
		http.Error(w, "Failed to get movies by genre", http.StatusInternalServerError)
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

	response := models.NewPaginatedResponse(result.Results, meta, "Movies by genre retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetGenres handles genre list requests
func (c *MovieController) GetGenres(w http.ResponseWriter, r *http.Request) {
	// Get genres
	genres, err := c.tmdbService.GetGenres(r.Context())
	if err != nil {
		c.logger.LogError(err, "GetGenres", r)
		http.Error(w, "Failed to get genres", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(genres, "Genres retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetSimilarMovies handles similar movies requests
func (c *MovieController) GetSimilarMovies(w http.ResponseWriter, r *http.Request) {
	// Get movie ID from URL
	vars := mux.Vars(r)
	movieIDStr := vars["id"]

	movieID, err := strconv.Atoi(movieIDStr)
	if err != nil {
		http.Error(w, "Invalid movie ID", http.StatusBadRequest)
		return
	}

	// Get limit parameter
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	// Get similar movies
	similarMovies, err := c.watchlistService.GetSimilarMovies(r.Context(), movieID, limit)
	if err != nil {
		c.logger.LogError(err, "GetSimilarMovies", r)
		http.Error(w, "Failed to get similar movies", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(similarMovies, "Similar movies retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetMediaDetails handles unified details requests for both movies and TV shows
func (c *MovieController) GetMediaDetails(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaType := vars["type"]
	idStr := vars["id"]

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	if mediaType == "movie" {
		movie, err := c.tmdbService.GetMovieDetails(r.Context(), id)
		if err != nil {
			c.logger.LogError(err, "GetMediaDetails (movie)", r)
			http.Error(w, "Failed to get movie details", http.StatusInternalServerError)
			return
		}
		// Enrich with OMDB data if you want
		_ = c.omdbService.EnrichMovieWithOMDBData(r.Context(), movie)
		response := models.NewSuccessResponse(movie, "Movie details retrieved successfully")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	} else if mediaType == "tv" {
		tv, err := c.tmdbService.GetTVDetails(r.Context(), id)
		if err != nil {
			c.logger.LogError(err, "GetMediaDetails (tv)", r)
			http.Error(w, "Failed to get TV show details", http.StatusInternalServerError)
			return
		}
		response := models.NewSuccessResponse(tv, "TV show details retrieved successfully")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	} else {
		http.Error(w, "Invalid media type", http.StatusBadRequest)
		return
	}
}
