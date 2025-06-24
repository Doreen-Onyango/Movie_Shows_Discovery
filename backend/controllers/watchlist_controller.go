package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/middleware"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/services"
)

// WatchlistController handles watchlist-related HTTP requests
type WatchlistController struct {
	watchlistService *services.WatchlistService
	logger           *middleware.Logger
}

// NewWatchlistController creates a new watchlist controller
func NewWatchlistController(watchlistService *services.WatchlistService, logger *middleware.Logger) *WatchlistController {
	return &WatchlistController{
		watchlistService: watchlistService,
		logger:           logger,
	}
}

// CreateWatchlist handles watchlist creation requests
func (c *WatchlistController) CreateWatchlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request (in a real app, this would come from authentication)
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request models.WatchlistCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Create watchlist
	watchlist, err := c.watchlistService.CreateWatchlist(r.Context(), userID, request)
	if err != nil {
		c.logger.LogError(err, "CreateWatchlist", r)
		http.Error(w, "Failed to create watchlist", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(watchlist, "Watchlist created successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetWatchlist handles watchlist retrieval requests
func (c *WatchlistController) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get watchlist
	watchlist, err := c.watchlistService.GetWatchlist(r.Context(), userID)
	if err != nil {
		c.logger.LogError(err, "GetWatchlist", r)
		http.Error(w, "Failed to get watchlist", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(watchlist, "Watchlist retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddToWatchlist handles adding movies to watchlist
func (c *WatchlistController) AddToWatchlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request models.WatchlistItemAddRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Add to watchlist
	item, err := c.watchlistService.AddToWatchlist(r.Context(), userID, request)
	if err != nil {
		c.logger.LogError(err, "AddToWatchlist", r)
		http.Error(w, "Failed to add to watchlist", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(item, "Movie added to watchlist successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateWatchlistItem handles updating watchlist items
func (c *WatchlistController) UpdateWatchlistItem(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get item ID from URL
	itemIDStr := r.URL.Query().Get("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request models.WatchlistItemUpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update watchlist item
	item, err := c.watchlistService.UpdateWatchlistItem(r.Context(), userID, itemID, request)
	if err != nil {
		c.logger.LogError(err, "UpdateWatchlistItem", r)
		http.Error(w, "Failed to update watchlist item", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(item, "Watchlist item updated successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// RemoveFromWatchlist handles removing movies from watchlist
func (c *WatchlistController) RemoveFromWatchlist(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get item ID from URL
	itemIDStr := r.URL.Query().Get("item_id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		http.Error(w, "Invalid item ID", http.StatusBadRequest)
		return
	}

	// Remove from watchlist
	err = c.watchlistService.RemoveFromWatchlist(r.Context(), userID, itemID)
	if err != nil {
		c.logger.LogError(err, "RemoveFromWatchlist", r)
		http.Error(w, "Failed to remove from watchlist", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(nil, "Movie removed from watchlist successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetWatchlistStats handles watchlist statistics requests
func (c *WatchlistController) GetWatchlistStats(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get watchlist stats
	stats, err := c.watchlistService.GetWatchlistStats(r.Context(), userID)
	if err != nil {
		c.logger.LogError(err, "GetWatchlistStats", r)
		http.Error(w, "Failed to get watchlist stats", http.StatusInternalServerError)
		return
	}

	// Create response
	response := models.NewSuccessResponse(stats, "Watchlist stats retrieved successfully")

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetRecommendations handles movie recommendation requests
func (c *WatchlistController) GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// Get user ID from request
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Get limit parameter
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10
	}

	// Get recommendations
	recommendations, err := c.watchlistService.GetRecommendations(r.Context(), userID, limit)
	if err != nil {
		c.logger.LogError(err, "GetRecommendations", r)
		http.Error(w, "Failed to get recommendations", http.StatusInternalServerError)
		return
	}

	// Create response
	meta := models.Meta{
		Page:         1,
		PerPage:      limit,
		TotalPages:   1,
		TotalResults: len(recommendations),
		HasNext:      false,
		HasPrev:      false,
	}

	response := models.NewRecommendationResponse(recommendations, userID, meta)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
