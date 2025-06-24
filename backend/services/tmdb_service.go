package services

import (
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/config"
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
