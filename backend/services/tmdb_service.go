package services

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"time"

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

// SearchMedia searches for movies and/or TV shows using TMDB API
func (s *TMDBService) SearchMedia(ctx context.Context, query string, page int, perPage int, mediaType string, includeAdult bool) (*models.MovieSearchResult, error) {
	if perPage <= 0 {
		perPage = 10
	}
	var allResults []models.Movie

	if mediaType == "movie" || mediaType == "all" {
		// Search movies
		baseURL := s.config.BaseURL + "/search/movie"
		params := url.Values{}
		params.Set("api_key", s.config.APIKey)
		params.Set("query", query)
		params.Set("page", strconv.Itoa(1))
		params.Set("include_adult", strconv.FormatBool(includeAdult))
		params.Set("language", "en-US")

		resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			var tmdbResp TMDBSearchResponse
			if err := json.Unmarshal(body, &tmdbResp); err == nil {
				for _, tmdbMovie := range tmdbResp.Results {
					movie := *s.convertTMDBMovie(tmdbMovie)
					movie.MediaType = "movie"
					allResults = append(allResults, movie)
				}
			}
			resp.Body.Close()
		}
	}

	if mediaType == "tv" || mediaType == "all" {
		// Search TV shows
		baseURL := s.config.BaseURL + "/search/tv"
		params := url.Values{}
		params.Set("api_key", s.config.APIKey)
		params.Set("query", query)
		params.Set("page", strconv.Itoa(1))
		params.Set("language", "en-US")

		resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			var tvResp struct {
				Page    int `json:"page"`
				Results []struct {
					ID            int      `json:"id"`
					Name          string   `json:"name"`
					OriginalName  string   `json:"original_name"`
					Overview      string   `json:"overview"`
					PosterPath    string   `json:"poster_path"`
					BackdropPath  string   `json:"backdrop_path"`
					FirstAirDate  string   `json:"first_air_date"`
					VoteAverage   float64  `json:"vote_average"`
					VoteCount     int      `json:"vote_count"`
					Popularity    float64  `json:"popularity"`
					GenreIDs      []int    `json:"genre_ids"`
					OriginCountry []string `json:"origin_country"`
					// Add more fields as needed
				} `json:"results"`
				TotalPages   int `json:"total_pages"`
				TotalResults int `json:"total_results"`
			}
			if err := json.Unmarshal(body, &tvResp); err == nil {
				for _, tv := range tvResp.Results {
					movie := models.Movie{
						ID:            tv.ID,
						Title:         tv.Name,
						OriginalTitle: tv.OriginalName,
						Overview:      tv.Overview,
						PosterPath:    tv.PosterPath,
						BackdropPath:  tv.BackdropPath,
						ReleaseDate:   tv.FirstAirDate,
						VoteAverage:   tv.VoteAverage,
						VoteCount:     tv.VoteCount,
						Popularity:    tv.Popularity,
						GenreIDs:      tv.GenreIDs,
						MediaType:     "tv",
					}
					allResults = append(allResults, movie)
				}
			}
			resp.Body.Close()
		}
	}

	// Manual pagination
	totalResults := len(allResults)
	totalPages := (totalResults + perPage - 1) / perPage
	start := (page - 1) * perPage
	end := start + perPage
	if start > totalResults {
		start = totalResults
	}
	if end > totalResults {
		end = totalResults
	}
	paged := []models.Movie{}
	if start < end {
		paged = allResults[start:end]
	}

	result := &models.MovieSearchResult{
		Page:         page,
		Results:      paged,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	return result, nil
}

// GetMovieDetails retrieves detailed movie information
func (s *TMDBService) GetMovieDetails(ctx context.Context, movieID int) (*models.Movie, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("tmdb_movie", movieID)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if movie, ok := cached.(*models.Movie); ok {
			return movie, nil
		}
	}

	// Build URL
	baseURL := fmt.Sprintf("%s/movie/%d", s.config.BaseURL, movieID)
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("append_to_response", "credits")
	params.Set("language", "en-US")

	// Make request
	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie details: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var tmdbMovie TMDBMovieResponse
	if err := json.Unmarshal(body, &tmdbMovie); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our model
	movie := s.convertTMDBMovie(tmdbMovie)

	// Get credits separately if not included
	if len(movie.Credits.Cast) == 0 {
		credits, err := s.GetMovieCredits(ctx, movieID)
		if err == nil {
			movie.Credits = *credits
		}
	}

	// Fetch trailer video from TMDB
	trailerKey, err := s.getTrailerKey(ctx, movieID)
	if err == nil && trailerKey != "" {
		movie.TrailerKey = trailerKey
	}

	// Validate and set ratings
	movie.Ratings.TMDB = movie.VoteAverage

	// Cache the result
	s.cache.Set(cacheKey, movie, config.AppConfig.Cache.TTL)

	return movie, nil
}

// getTrailerKey fetches the first YouTube trailer key for a movie
func (s *TMDBService) getTrailerKey(ctx context.Context, movieID int) (string, error) {
	baseURL := fmt.Sprintf("%s/movie/%d/videos", s.config.BaseURL, movieID)
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("language", "en-US")

	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var videoResp struct {
		Results []struct {
			Key      string `json:"key"`
			Site     string `json:"site"`
			Type     string `json:"type"`
			Official bool   `json:"official"`
			Name     string `json:"name"`
		} `json:"results"`
	}
	if err := json.Unmarshal(body, &videoResp); err != nil {
		return "", err
	}

	for _, v := range videoResp.Results {
		if v.Site == "YouTube" && v.Type == "Trailer" {
			return v.Key, nil
		}
	}
	return "", nil
}

// GetMovieCredits retrieves cast and crew information
func (s *TMDBService) GetMovieCredits(ctx context.Context, movieID int) (*models.Credits, error) {
	// Build URL
	baseURL := fmt.Sprintf("%s/movie/%d/credits", s.config.BaseURL, movieID)
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("language", "en-US")

	// Make request
	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get movie credits: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var tmdbCredits TMDBCreditsResponse
	if err := json.Unmarshal(body, &tmdbCredits); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our model
	credits := &models.Credits{
		Cast: make([]models.CastMember, len(tmdbCredits.Cast)),
		Crew: make([]models.CrewMember, len(tmdbCredits.Crew)),
	}

	for i, cast := range tmdbCredits.Cast {
		credits.Cast[i] = models.CastMember{
			ID:          cast.ID,
			Name:        cast.Name,
			Character:   cast.Character,
			ProfilePath: cast.ProfilePath,
			Order:       cast.Order,
		}
	}

	for i, crew := range tmdbCredits.Crew {
		credits.Crew[i] = models.CrewMember{
			ID:          crew.ID,
			Name:        crew.Name,
			Job:         crew.Job,
			Department:  crew.Department,
			ProfilePath: crew.ProfilePath,
		}
	}

	return credits, nil
}

// GetTrendingMedia retrieves trending movies and/or TV shows
func (s *TMDBService) GetTrendingMedia(ctx context.Context, timeframe string, page int, mediaType string) (*models.MovieSearchResult, error) {
	if timeframe != "day" && timeframe != "week" {
		timeframe = "day"
	}

	var allResults []models.Movie

	if mediaType == "movie" || mediaType == "all" {
		baseURL := fmt.Sprintf("%s/trending/movie/%s", s.config.BaseURL, timeframe)
		params := url.Values{}
		params.Set("api_key", s.config.APIKey)
		params.Set("page", strconv.Itoa(page))
		params.Set("language", "en-US")

		resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			var tmdbResp TMDBTrendingResponse
			if err := json.Unmarshal(body, &tmdbResp); err == nil {
				for _, tmdbMovie := range tmdbResp.Results {
					movie := *s.convertTMDBMovie(tmdbMovie)
					movie.MediaType = "movie"
					allResults = append(allResults, movie)
				}
			}
			resp.Body.Close()
		}
	}

	if mediaType == "tv" || mediaType == "all" {
		baseURL := fmt.Sprintf("%s/trending/tv/%s", s.config.BaseURL, timeframe)
		params := url.Values{}
		params.Set("api_key", s.config.APIKey)
		params.Set("page", strconv.Itoa(page))
		params.Set("language", "en-US")

		resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
		if err == nil {
			body, _ := io.ReadAll(resp.Body)
			var tvResp struct {
				Page    int `json:"page"`
				Results []struct {
					ID            int      `json:"id"`
					Name          string   `json:"name"`
					OriginalName  string   `json:"original_name"`
					Overview      string   `json:"overview"`
					PosterPath    string   `json:"poster_path"`
					BackdropPath  string   `json:"backdrop_path"`
					FirstAirDate  string   `json:"first_air_date"`
					VoteAverage   float64  `json:"vote_average"`
					VoteCount     int      `json:"vote_count"`
					Popularity    float64  `json:"popularity"`
					GenreIDs      []int    `json:"genre_ids"`
					OriginCountry []string `json:"origin_country"`
				} `json:"results"`
				TotalPages   int `json:"total_pages"`
				TotalResults int `json:"total_results"`
			}
			if err := json.Unmarshal(body, &tvResp); err == nil {
				for _, tv := range tvResp.Results {
					movie := models.Movie{
						ID:            tv.ID,
						Title:         tv.Name,
						OriginalTitle: tv.OriginalName,
						Overview:      tv.Overview,
						PosterPath:    tv.PosterPath,
						BackdropPath:  tv.BackdropPath,
						ReleaseDate:   tv.FirstAirDate,
						VoteAverage:   tv.VoteAverage,
						VoteCount:     tv.VoteCount,
						Popularity:    tv.Popularity,
						GenreIDs:      tv.GenreIDs,
						MediaType:     "tv",
					}
					allResults = append(allResults, movie)
				}
			}
			resp.Body.Close()
		}
	}

	totalResults := len(allResults)
	totalPages := (totalResults + 19) / 20 // 20 per page
	start := (page - 1) * 20
	end := start + 20
	if start > totalResults {
		start = totalResults
	}
	if end > totalResults {
		end = totalResults
	}
	paged := []models.Movie{}
	if start < end {
		paged = allResults[start:end]
	}

	result := &models.MovieSearchResult{
		Page:         page,
		Results:      paged,
		TotalPages:   totalPages,
		TotalResults: totalResults,
	}

	return result, nil
}

// GetMoviesByGenre retrieves movies by genre
func (s *TMDBService) GetMoviesByGenre(ctx context.Context, genreID int, page int, sortBy string) (*models.MovieSearchResult, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("tmdb_genre", genreID, page, sortBy)

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if result, ok := cached.(*models.MovieSearchResult); ok {
			return result, nil
		}
	}

	// Build URL
	baseURL := fmt.Sprintf("%s/discover/movie", s.config.BaseURL)
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("with_genres", strconv.Itoa(genreID))
	params.Set("page", strconv.Itoa(page))
	params.Set("sort_by", sortBy)
	params.Set("language", "en-US")
	params.Set("include_adult", "false")

	// Make request
	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get movies by genre: %w", err)
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

	result := &models.MovieSearchResult{
		Page:         tmdbResp.Page,
		Results:      movies,
		TotalPages:   tmdbResp.TotalPages,
		TotalResults: tmdbResp.TotalResults,
	}

	// Cache the result
	s.cache.Set(cacheKey, result, config.AppConfig.Cache.TTL)

	return result, nil
}

// GetGenres retrieves all available genres
func (s *TMDBService) GetGenres(ctx context.Context) ([]models.Genre, error) {
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("tmdb_genres")

	// Check cache first
	if cached, exists := s.cache.Get(cacheKey); exists {
		if genres, ok := cached.([]models.Genre); ok {
			return genres, nil
		}
	}

	// Build URL
	baseURL := s.config.BaseURL + "/genre/movie/list"
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("language", "en-US")

	// Make request
	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get genres: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse response
	var tmdbResp TMDBGenreResponse
	if err := json.Unmarshal(body, &tmdbResp); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Convert to our model
	genres := make([]models.Genre, len(tmdbResp.Genres))
	for i, genre := range tmdbResp.Genres {
		genres[i] = models.Genre{
			ID:   genre.ID,
			Name: genre.Name,
		}
	}

	// Cache the result (genres don't change often)
	s.cache.Set(cacheKey, genres, 24*time.Hour)

	return genres, nil
}

// convertTMDBMovie converts TMDB movie response to our model
func (s *TMDBService) convertTMDBMovie(tmdbMovie TMDBMovieResponse) *models.Movie {
	movie := &models.Movie{
		ID:            tmdbMovie.ID,
		Title:         tmdbMovie.Title,
		OriginalTitle: tmdbMovie.OriginalTitle,
		Overview:      tmdbMovie.Overview,
		PosterPath:    tmdbMovie.PosterPath,
		BackdropPath:  tmdbMovie.BackdropPath,
		ReleaseDate:   tmdbMovie.ReleaseDate,
		Runtime:       tmdbMovie.Runtime,
		Status:        tmdbMovie.Status,
		Tagline:       tmdbMovie.Tagline,
		VoteAverage:   tmdbMovie.VoteAverage,
		VoteCount:     tmdbMovie.VoteCount,
		Popularity:    tmdbMovie.Popularity,
		Adult:         tmdbMovie.Adult,
		Video:         tmdbMovie.Video,
		GenreIDs:      tmdbMovie.GenreIDs,
		MediaType:     "movie",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Convert genres
	if len(tmdbMovie.Genres) > 0 {
		movie.Genres = make([]models.Genre, len(tmdbMovie.Genres))
		for i, genre := range tmdbMovie.Genres {
			movie.Genres[i] = models.Genre{
				ID:   genre.ID,
				Name: genre.Name,
			}
		}
	}

	// Convert production companies
	if len(tmdbMovie.ProductionCompanies) > 0 {
		movie.ProductionCompanies = make([]models.ProductionCompany, len(tmdbMovie.ProductionCompanies))
		for i, company := range tmdbMovie.ProductionCompanies {
			movie.ProductionCompanies[i] = models.ProductionCompany{
				ID:            company.ID,
				Name:          company.Name,
				LogoPath:      company.LogoPath,
				OriginCountry: company.OriginCountry,
			}
		}
	}

	// Convert spoken languages
	if len(tmdbMovie.SpokenLanguages) > 0 {
		movie.SpokenLanguages = make([]models.SpokenLanguage, len(tmdbMovie.SpokenLanguages))
		for i, lang := range tmdbMovie.SpokenLanguages {
			movie.SpokenLanguages[i] = models.SpokenLanguage{
				ISO6391: lang.ISO6391,
				Name:    lang.Name,
			}
		}
	}

	// Initialize ratings
	movie.Ratings = models.Ratings{
		TMDB: tmdbMovie.VoteAverage,
	}

	// Validate movie data
	utils.ValidateMovieData(movie)

	return movie
}

// Close closes the service and cleans up resources
func (s *TMDBService) Close() {
	if s.client != nil {
		s.client.Close()
	}
}

// GetTVDetails fetches TV show details from TMDB
func (s *TMDBService) GetTVDetails(ctx context.Context, tvID int) (*models.TV, error) {
	baseURL := s.config.BaseURL + "/tv/" + strconv.Itoa(tvID)
	params := url.Values{}
	params.Set("api_key", s.config.APIKey)
	params.Set("language", "en-US")

	resp, err := s.client.Get(ctx, baseURL+"?"+params.Encode(), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tmdbTV struct {
		ID               int     `json:"id"`
		Name             string  `json:"name"`
		OriginalName     string  `json:"original_name"`
		Overview         string  `json:"overview"`
		PosterPath       string  `json:"poster_path"`
		BackdropPath     string  `json:"backdrop_path"`
		FirstAirDate     string  `json:"first_air_date"`
		LastAirDate      string  `json:"last_air_date"`
		NumberOfSeasons  int     `json:"number_of_seasons"`
		NumberOfEpisodes int     `json:"number_of_episodes"`
		Status           string  `json:"status"`
		Tagline          string  `json:"tagline"`
		VoteAverage      float64 `json:"vote_average"`
		VoteCount        int     `json:"vote_count"`
		Popularity       float64 `json:"popularity"`
		Genres           []struct {
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
	if err := json.Unmarshal(body, &tmdbTV); err != nil {
		return nil, err
	}

	// Convert to models.TV
	genres := make([]models.Genre, len(tmdbTV.Genres))
	for i, g := range tmdbTV.Genres {
		genres[i] = models.Genre{ID: g.ID, Name: g.Name}
	}
	prodCompanies := make([]models.ProductionCompany, len(tmdbTV.ProductionCompanies))
	for i, pc := range tmdbTV.ProductionCompanies {
		prodCompanies[i] = models.ProductionCompany{ID: pc.ID, Name: pc.Name, LogoPath: pc.LogoPath, OriginCountry: pc.OriginCountry}
	}
	languages := make([]models.SpokenLanguage, len(tmdbTV.SpokenLanguages))
	for i, l := range tmdbTV.SpokenLanguages {
		languages[i] = models.SpokenLanguage{ISO6391: l.ISO6391, Name: l.Name}
	}

	tv := &models.TV{
		ID:                  tmdbTV.ID,
		Name:                tmdbTV.Name,
		OriginalName:        tmdbTV.OriginalName,
		Overview:            tmdbTV.Overview,
		PosterPath:          tmdbTV.PosterPath,
		BackdropPath:        tmdbTV.BackdropPath,
		FirstAirDate:        tmdbTV.FirstAirDate,
		LastAirDate:         tmdbTV.LastAirDate,
		NumberOfSeasons:     tmdbTV.NumberOfSeasons,
		NumberOfEpisodes:    tmdbTV.NumberOfEpisodes,
		Status:              tmdbTV.Status,
		Tagline:             tmdbTV.Tagline,
		VoteAverage:         tmdbTV.VoteAverage,
		VoteCount:           tmdbTV.VoteCount,
		Popularity:          tmdbTV.Popularity,
		Genres:              genres,
		ProductionCompanies: prodCompanies,
		SpokenLanguages:     languages,
		MediaType:           "tv",
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}
	return tv, nil
}
