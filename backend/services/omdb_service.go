package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/config"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/utils"
)

// OMDBService handles all OMDB API interactions
type OMDBService struct {
	client *utils.HTTPClient
	cache  *utils.Cache
	config *config.OMDBConfig
}

// OMDBMovieResponse represents OMDB movie response
type OMDBMovieResponse struct {
	Title    string `json:"Title"`
	Year     string `json:"Year"`
	Rated    string `json:"Rated"`
	Released string `json:"Released"`
	Runtime  string `json:"Runtime"`
	Genre    string `json:"Genre"`
	Director string `json:"Director"`
	Writer   string `json:"Writer"`
	Actors   string `json:"Actors"`
	Plot     string `json:"Plot"`
	Poster   string `json:"Poster"`
	Ratings  []struct {
		Source string `json:"Source"`
		Value  string `json:"Value"`
	} `json:"Ratings"`
	Metascore  string `json:"Metascore"`
	IMDBRating string `json:"imdbRating"`
	IMDBVotes  string `json:"imdbVotes"`
	IMDBID     string `json:"imdbID"`
	Type       string `json:"Type"`
	DVD        string `json:"DVD"`
	BoxOffice  string `json:"BoxOffice"`
	Production string `json:"Production"`
	Website    string `json:"Website"`
	Response   string `json:"Response"`
	Error      string `json:"Error"`
}

// NewOMDBService creates a new OMDB service instance
func NewOMDBService() *OMDBService {
	return &OMDBService{
		client: utils.CreateOMDBClient(),
		cache:  utils.NewCache(),
		config: &config.AppConfig.OMDB,
	}
}

// GetMovieByTitle retrieves movie information by title
func (s *OMDBService) GetMovieByTitle(ctx context.Context, title string, year string) (*OMDBMovieResponse, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("omdb_title", title, year)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if movie, ok := cached.(*OMDBMovieResponse); ok {
			return movie, nil
		}
	}

	// Build URL
	params := url.Values{}
	params.Set("apikey", s.config.APIKey)
	params.Set("t", title)
	if year != "" {
		params.Set("y", year)
	}
	params.Set("plot", "full")

	// Make request
	resp, err := s.client.Get(ctx, s.config.BaseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie by title: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var omdbResp OMDBMovieResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDB API error: %s", omdbResp.Error)
	}

	// Cache the result
	s.cache.Set(cacheKey, &omdbResp, config.AppConfig.Cache.TTL)

	return &omdbResp, nil
}

// GetMovieByIMDBID retrieves movie information by IMDB ID
func (s *OMDBService) GetMovieByIMDBID(ctx context.Context, imdbID string) (*OMDBMovieResponse, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("omdb_imdb", imdbID)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if movie, ok := cached.(*OMDBMovieResponse); ok {
			return movie, nil
		}
	}

	// Build URL
	params := url.Values{}
	params.Set("apikey", s.config.APIKey)
	params.Set("i", imdbID)
	params.Set("plot", "full")

	// Make request
	resp, err := s.client.Get(ctx, s.config.BaseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie by IMDB ID: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var omdbResp OMDBMovieResponse
	if err := json.Unmarshal(body, &omdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if omdbResp.Response == "False" {
		return nil, fmt.Errorf("OMDB API error: %s", omdbResp.Error)
	}

	// Cache the result
	s.cache.Set(cacheKey, &omdbResp, config.AppConfig.Cache.TTL)

	return &omdbResp, nil
}

// SearchMovies searches for movies using OMDB API
func (s *OMDBService) SearchMovies(ctx context.Context, query string, page int) ([]OMDBMovieResponse, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("omdb_search", query, page)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if movies, ok := cached.([]OMDBMovieResponse); ok {
			return movies, nil
		}
	}

	// Build URL
	params := url.Values{}
	params.Set("apikey", s.config.APIKey)
	params.Set("s", query)
	params.Set("type", "movie")
	params.Set("page", strconv.Itoa(page))

	// Make request
	resp, err := s.client.Get(ctx, s.config.BaseURL+"?"+params.Encode(), nil)
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
	var searchResp struct {
		Search       []OMDBMovieResponse `json:"Search"`
		TotalResults string              `json:"totalResults"`
		Response     string              `json:"Response"`
		Error        string              `json:"Error"`
	}

	if err := json.Unmarshal(body, &searchResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if searchResp.Response == "False" {
		return nil, fmt.Errorf("OMDB API error: %s", searchResp.Error)
	}

	// Cache the result
	s.cache.Set(cacheKey, searchResp.Search, config.AppConfig.Cache.SearchTTL)

	return searchResp.Search, nil
}

// ExtractRatings extracts ratings from OMDB response
func (s *OMDBService) ExtractRatings(omdbResp *OMDBMovieResponse) models.Ratings {
	ratings := models.Ratings{}

	// Extract IMDB rating
	if omdbResp.IMDBRating != "" && omdbResp.IMDBRating != "N/A" {
		if rating, err := strconv.ParseFloat(omdbResp.IMDBRating, 64); err == nil {
			ratings.IMDB = rating
		}
	}

	// Extract Metascore
	if omdbResp.Metascore != "" && omdbResp.Metascore != "N/A" {
		if rating, err := strconv.ParseFloat(omdbResp.Metascore, 64); err == nil {
			ratings.Metacritic = rating
		}
	}

	// Extract ratings from Ratings array
	for _, rating := range omdbResp.Ratings {
		switch strings.ToLower(rating.Source) {
		case "rotten tomatoes":
			if rating.Value != "" && rating.Value != "N/A" {
				// Remove % sign and convert to float
				value := strings.TrimSuffix(rating.Value, "%")
				if rtRating, err := strconv.ParseFloat(value, 64); err == nil {
					ratings.RottenTomatoes = rtRating
				}
			}
		}
	}

	// Calculate OMDB rating (average of available ratings)
	var sum float64
	var count int

	if ratings.IMDB > 0 {
		sum += ratings.IMDB
		count++
	}
	if ratings.Metacritic > 0 {
		sum += (ratings.Metacritic / 10.0) // Convert to 10-point scale
		count++
	}
	if ratings.RottenTomatoes > 0 {
		sum += (ratings.RottenTomatoes / 10.0) // Convert to 10-point scale
		count++
	}

	if count > 0 {
		ratings.OMDB = sum / float64(count)
	}

	// Validate ratings
	utils.ValidateRatings(&ratings)

	return ratings
}

// EnrichMovieWithOMDBData enriches a movie with OMDB data
func (s *OMDBService) EnrichMovieWithOMDBData(ctx context.Context, movie *models.Movie) error {
	// Try to get OMDB data by title and year
	year := utils.ParseYear(movie.ReleaseDate)
	yearStr := ""
	if year > 0 {
		yearStr = strconv.Itoa(year)
	}

	omdbResp, err := s.GetMovieByTitle(ctx, movie.Title, yearStr)
	if err != nil {
		// If title search fails, try with original title
		if movie.OriginalTitle != "" && movie.OriginalTitle != movie.Title {
			omdbResp, err = s.GetMovieByTitle(ctx, movie.OriginalTitle, yearStr)
			if err != nil {
				return fmt.Errorf("failed to get OMDB data: %w", err)
			}
		} else {
			return fmt.Errorf("failed to get OMDB data: %w", err)
		}
	}

	// Extract and merge ratings
	omdbRatings := s.ExtractRatings(omdbResp)

	// Merge ratings (keep existing TMDB rating)
	if omdbRatings.IMDB > 0 {
		movie.Ratings.IMDB = omdbRatings.IMDB
	}
	if omdbRatings.RottenTomatoes > 0 {
		movie.Ratings.RottenTomatoes = omdbRatings.RottenTomatoes
	}
	if omdbRatings.Metacritic > 0 {
		movie.Ratings.Metacritic = omdbRatings.Metacritic
	}
	if omdbRatings.OMDB > 0 {
		movie.Ratings.OMDB = omdbRatings.OMDB
	}

	// Update movie with additional OMDB data if missing
	if movie.Overview == "" || movie.Overview == "No overview available" {
		if omdbResp.Plot != "" && omdbResp.Plot != "N/A" {
			movie.Overview = omdbResp.Plot
		}
	}

	if movie.Runtime == 0 && omdbResp.Runtime != "" && omdbResp.Runtime != "N/A" {
		// Parse runtime (format: "120 min")
		runtimeStr := strings.TrimSuffix(omdbResp.Runtime, " min")
		if runtime, err := strconv.Atoi(runtimeStr); err == nil {
			movie.Runtime = runtime
		}
	}

	// Validate the updated movie data
	utils.ValidateMovieData(movie)

	return nil
}

// Close closes the service and cleans up resources
func (s *OMDBService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}
