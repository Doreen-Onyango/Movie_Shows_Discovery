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
