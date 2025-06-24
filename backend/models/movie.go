package models

import "time"

// Movie represents a movie with comprehensive information
type Movie struct {
	ID                  int                 `json:"id"`
	Title               string              `json:"title"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	PosterPath          string              `json:"poster_path"`
	BackdropPath        string              `json:"backdrop_path"`
	ReleaseDate         string              `json:"release_date"`
	Runtime             int                 `json:"runtime"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	VoteAverage         float64             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
	Popularity          float64             `json:"popularity"`
	Adult               bool                `json:"adult"`
	Video               bool                `json:"video"`
	GenreIDs            []int               `json:"genre_ids"`
	Genres              []Genre             `json:"genres"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Credits             Credits             `json:"credits"`
	Ratings             Ratings             `json:"ratings"`
	MediaType           string              `json:"media_type"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	TrailerKey          string              `json:"trailerKey,omitempty"`
}

// Genre represents a movie genre
type Genre struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ProductionCompany represents a production company
type ProductionCompany struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

// SpokenLanguage represents a spoken language
type SpokenLanguage struct {
	ISO6391 string `json:"iso_639_1"`
	Name    string `json:"name"`
}

// Credits represents cast and crew information
type Credits struct {
	Cast []CastMember `json:"cast"`
	Crew []CrewMember `json:"crew"`
}

// CastMember represents a cast member
type CastMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Character   string `json:"character"`
	ProfilePath string `json:"profile_path"`
	Order       int    `json:"order"`
}

// CrewMember represents a crew member
type CrewMember struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Job         string `json:"job"`
	Department  string `json:"department"`
	ProfilePath string `json:"profile_path"`
}

// Ratings represents ratings from different sources
type Ratings struct {
	TMDB           float64 `json:"tmdb"`
	OMDB           float64 `json:"omdb"`
	RottenTomatoes float64 `json:"rotten_tomatoes"`
	IMDB           float64 `json:"imdb"`
	Metacritic     float64 `json:"metacritic"`
}

// MovieSearchResult represents a movie search result
type MovieSearchResult struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

// MovieRecommendation represents a movie recommendation
type MovieRecommendation struct {
	Movie      Movie   `json:"movie"`
	Score      float64 `json:"score"`
	Reason     string  `json:"reason"`
	Similarity float64 `json:"similarity"`
}

// MovieFilter represents filters for movie search
type MovieFilter struct {
	GenreIDs     []int   `json:"genre_ids"`
	Year         int     `json:"year"`
	MinRating    float64 `json:"min_rating"`
	MaxRating    float64 `json:"max_rating"`
	Language     string  `json:"language"`
	SortBy       string  `json:"sort_by"`
	IncludeAdult bool    `json:"include_adult"`
}

// MovieList represents a list of movies with pagination
type MovieList struct {
	Movies       []Movie `json:"movies"`
	Page         int     `json:"page"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
	HasNext      bool    `json:"has_next"`
	HasPrev      bool    `json:"has_prev"`
}
