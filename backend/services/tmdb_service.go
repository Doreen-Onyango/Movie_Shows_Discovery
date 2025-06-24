package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/config"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/utils"
)

// TMDBService handles all TMDB API interactions
type TMDBService struct {
	client *utils.HTTPClient
	cache  *utils.Cache
	config *config.TMDBConfig
}

// TMDBMovieResponse represents TMDB movie response
type TMDBMovieResponse struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Overview      string  `json:"overview"`
	PosterPath    string  `json:"poster_path"`
	BackdropPath  string  `json:"backdrop_path"`
	ReleaseDate   string  `json:"release_date"`
	Runtime       int     `json:"runtime"`
	Status        string  `json:"status"`
	Tagline       string  `json:"tagline"`
	VoteAverage   float64 `json:"vote_average"`
	VoteCount     int     `json:"vote_count"`
	Popularity    float64 `json:"popularity"`
	Adult         bool    `json:"adult"`
	Video         bool    `json:"video"`
	GenreIDs      []int   `json:"genre_ids"`
	Genres        []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
	ProductionCompanies []struct {
		ID            int    `json:"id"`
		Name          string `json:"name"`
		LogoPath      string `json:"logo_path"`
		OriginCountry string `json:"origin_country"`
	} `json:"production_companies"`
	SpokenLanguages []struct {
		ISO6391 string `json:"iso_639_1"`
		Name    string `json:"name"`
	} `json:"spoken_languages"`
}

// TMDBCreditsResponse represents TMDB credits response
type TMDBCreditsResponse struct {
	ID   int `json:"id"`
	Cast []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Character   string `json:"character"`
		ProfilePath string `json:"profile_path"`
		Order       int    `json:"order"`
	} `json:"cast"`
	Crew []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Job         string `json:"job"`
		Department  string `json:"department"`
		ProfilePath string `json:"profile_path"`
	} `json:"crew"`
}

// TMDBSearchResponse represents TMDB search response
type TMDBSearchResponse struct {
	Page         int                 `json:"page"`
	Results      []TMDBMovieResponse `json:"results"`
	TotalPages   int                 `json:"total_pages"`
	TotalResults int                 `json:"total_results"`
}

// TMDBTrendingResponse represents TMDB trending response
type TMDBTrendingResponse struct {
	Page         int                 `json:"page"`
	Results      []TMDBMovieResponse `json:"results"`
	TotalPages   int                 `json:"total_pages"`
	TotalResults int                 `json:"total_results"`
}

// TMDBGenreResponse represents TMDB genre response
type TMDBGenreResponse struct {
	Genres []struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"genres"`
}

// NewTMDBService creates a new TMDB service instance
func NewTMDBService() *TMDBService {
	return &TMDBService{
		client: utils.CreateTMDBClient(),
		cache:  utils.NewCache(),
		config: &config.AppConfig.TMDB,
	}
}

// SearchMovies searches for movies using TMDB API
func (s *TMDBService) SearchMovies(ctx context.Context, query string, page int, perPage int, includeAdult bool) (*models.MovieSearchResult, error) {
	if perPage <= 0 {
		perPage = 10
	}
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("tmdb_search", query, page, perPage, includeAdult)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if result, ok := cached.(*models.MovieSearchResult); ok {
			return result, nil
		}
	}

	// Build URL
	baseURL := s.config.BaseURL + "/search/movie"
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("query", query)
	params.Set("page", strconv.Itoa(1)) // Always fetch first page from TMDB
	params.Set("include_adult", strconv.FormatBool(includeAdult))
	params.Set("language", "en-US")

	// Make request
	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search movies: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var tmdbResp TMDBSearchResponse
	if err := json.Unmarshal(body, &tmdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our model
	movies := make([]models.Movie, len(tmdbResp.Results))
	for i, tmdbMovie := range tmdbResp.Results {
		movies[i] = *s.convertTMDBMovie(tmdbMovie)
	}

	// Manual pagination
	totalResults := len(movies)
	totalPages := (totalResults + perPage - 1) / perPage
	start := (page - 1) * perPage
	end := start + perPage
	if start > totalResults {
		start = totalResults
	}
	if end > totalResults {
		end = totalResults
	}
	pagedMovies := []models.Movie{}
	if start < end {
		pagedMovies = movies[start:end]
	}

	result := &models.MovieSearchResult{
		Page:         page,
		Results:      pagedMovies,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	// Cache the result
	s.cache.Set(cacheKey, result, config.AppConfig.Cache.SearchTTL)

	return result, nil
}
