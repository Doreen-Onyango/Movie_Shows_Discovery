package routes

import (
	"net/http"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/controllers"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/middleware"
	"github.com/gorilla/mux"
)

// SetupRoutes configures all API routes
func SetupRoutes(
	movieController *controllers.MovieController,
	watchlistController *controllers.WatchlistController,
	trendingController *controllers.TrendingController,
	logger *middleware.Logger,
) *mux.Router {
	router := mux.NewRouter()

	// Apply middleware
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.ErrorLoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.RateLimitMiddleware(100)) // 100 requests per minute

	// Health check endpoint
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"healthy","service":"movie-discovery-api"}`))
	}).Methods("GET")

	// API version prefix
	api := router.PathPrefix("/api/v1").Subrouter()

	// Movie routes
	movieRoutes := api.PathPrefix("/movies").Subrouter()
	movieRoutes.HandleFunc("/search", movieController.SearchMovies).Methods("GET")
	movieRoutes.HandleFunc("/{id:[0-9]+}", movieController.GetMovieDetails).Methods("GET")
	movieRoutes.HandleFunc("/{id:[0-9]+}/similar", movieController.GetSimilarMovies).Methods("GET")
	movieRoutes.HandleFunc("/genres", movieController.GetGenres).Methods("GET")
	movieRoutes.HandleFunc("/genres/{genreId:[0-9]+}", movieController.GetMoviesByGenre).Methods("GET")

	// Trending routes
	trendingRoutes := api.PathPrefix("/trending").Subrouter()
	trendingRoutes.HandleFunc("", trendingController.GetTrending).Methods("GET")
	trendingRoutes.HandleFunc("/by-genre", trendingController.GetTrendingByGenre).Methods("GET")
	trendingRoutes.HandleFunc("/stats", trendingController.GetTrendingStats).Methods("GET")
	trendingRoutes.HandleFunc("/genres", trendingController.GetTrendingGenres).Methods("GET")

	// Media details route (unified for movie and tv)
	api.HandleFunc("/media/{type}/{id:[0-9]+}", movieController.GetMediaDetails).Methods("GET")

	// Watchlist routes
	watchlistRoutes := api.PathPrefix("/watchlist").Subrouter()
	watchlistRoutes.HandleFunc("", watchlistController.CreateWatchlist).Methods("POST")
	watchlistRoutes.HandleFunc("", watchlistController.GetWatchlist).Methods("GET")
	watchlistRoutes.HandleFunc("/items", watchlistController.AddToWatchlist).Methods("POST")
	watchlistRoutes.HandleFunc("/items", watchlistController.UpdateWatchlistItem).Methods("PUT")
	watchlistRoutes.HandleFunc("/items", watchlistController.RemoveFromWatchlist).Methods("DELETE")
	watchlistRoutes.HandleFunc("/stats", watchlistController.GetWatchlistStats).Methods("GET")
	watchlistRoutes.HandleFunc("/recommendations", watchlistController.GetRecommendations).Methods("GET")

	// 404 handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error":"Not Found","message":"The requested resource was not found"}`))
	})

	// Method not allowed handler
	router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(`{"error":"Method Not Allowed","message":"The requested method is not allowed for this resource"}`))
	})

	return router
}

// API Documentation
const APIDocumentation = `
# Movie Shows Discovery API

## Base URL
http://localhost:8080/api/v1

## Authentication
Most endpoints require a user ID in the X-User-ID header.

## Endpoints

### Health Check
GET /health
- Returns service health status

### Movies

#### Search Movies
GET /movies/search?q={query}&page={page}&include_adult={boolean}
- Search for movies by title
- Parameters:
  - q (required): Search query
  - page (optional): Page number (default: 1)
  - include_adult (optional): Include adult content (default: false)

#### Get Movie Details
GET /movies/{id}
- Get detailed information about a specific movie
- Parameters:
  - id (required): TMDB movie ID

#### Get Similar Movies
GET /movies/{id}/similar?limit={limit}
- Get movies similar to the specified movie
- Parameters:
  - id (required): TMDB movie ID
  - limit (optional): Number of results (default: 10)

#### Get Genres
GET /movies/genres
- Get all available movie genres

#### Get Movies by Genre
GET /movies/genres/{genreId}?page={page}&sort_by={sort_by}
- Get movies filtered by genre
- Parameters:
  - genreId (required): Genre ID
  - page (optional): Page number (default: 1)
  - sort_by (optional): Sort order (default: popularity.desc)

### Trending

#### Get Trending Movies
GET /trending?timeframe={timeframe}&page={page}
- Get trending movies
- Parameters:
  - timeframe (optional): "day" or "week" (default: "day")
  - page (optional): Page number (default: 1)

#### Get Trending by Genre
GET /trending/by-genre?genre_id={genreId}&page={page}&sort_by={sort_by}
- Get trending movies filtered by genre
- Parameters:
  - genre_id (required): Genre ID
  - page (optional): Page number (default: 1)
  - sort_by (optional): Sort order (default: popularity.desc)

#### Get Trending Stats
GET /trending/stats
- Get trending statistics

#### Get Trending Genres
GET /trending/genres
- Get trending genres with popularity information

### Watchlist

#### Create Watchlist
POST /watchlist
- Create a new watchlist for a user
- Headers: X-User-ID (required)
- Body: {"name": "string", "description": "string", "is_public": boolean}

#### Get Watchlist
GET /watchlist
- Get user's watchlist
- Headers: X-User-ID (required)

#### Add to Watchlist
POST /watchlist/items
- Add a movie to the watchlist
- Headers: X-User-ID (required)
- Body: {"movie_id": number, "status": "string", "rating": number, "notes": "string"}

#### Update Watchlist Item
PUT /watchlist/items?item_id={itemId}
- Update a watchlist item
- Headers: X-User-ID (required)
- Body: {"status": "string", "rating": number, "notes": "string"}

#### Remove from Watchlist
DELETE /watchlist/items?item_id={itemId}
- Remove a movie from the watchlist
- Headers: X-User-ID (required)

#### Get Watchlist Stats
GET /watchlist/stats
- Get watchlist statistics
- Headers: X-User-ID (required)

#### Get Recommendations
GET /watchlist/recommendations?limit={limit}
- Get movie recommendations based on watchlist
- Headers: X-User-ID (required)
- Parameters:
  - limit (optional): Number of recommendations (default: 10)

## Response Format

All responses follow this format:
{
  "success": boolean,
  "message": "string",
  "data": object,
  "timestamp": "datetime",
  "meta": {
    "page": number,
    "per_page": number,
    "total_pages": number,
    "total_results": number,
    "has_next": boolean,
    "has_prev": boolean
  }
}

## Error Responses

Error responses include:
{
  "success": false,
  "error": "string",
  "message": "string",
  "code": number,
  "timestamp": "datetime"
}

## Rate Limiting

- 100 requests per minute per IP address
- Rate limit headers are included in responses

## Caching

- Search results are cached for 30 minutes
- Trending content is cached for 1 hour
- Movie details are cached for 1 hour
- Genres are cached for 24 hours
`
