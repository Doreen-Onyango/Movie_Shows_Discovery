package models

import "time"

// Watchlist represents a user's watchlist
type Watchlist struct {
	ID          int             `json:"id"`
	UserID      string          `json:"user_id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	IsPublic    bool            `json:"is_public"`
	Items       []WatchlistItem `json:"items"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// WatchlistItem represents an item in a watchlist
type WatchlistItem struct {
	ID          int       `json:"id"`
	WatchlistID int       `json:"watchlist_id"`
	MovieID     int       `json:"movie_id"`
	Movie       Movie     `json:"movie"`
	Status      string    `json:"status"` // "to_watch", "watching", "completed", "dropped"
	Rating      float64   `json:"rating"`
	Notes       string    `json:"notes"`
	AddedAt     time.Time `json:"added_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// WatchlistCreateRequest represents a request to create a watchlist
type WatchlistCreateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// WatchlistUpdateRequest represents a request to update a watchlist
type WatchlistUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	IsPublic    bool   `json:"is_public"`
}

// WatchlistItemAddRequest represents a request to add an item to watchlist
type WatchlistItemAddRequest struct {
	MovieID int     `json:"movie_id" validate:"required"`
	Status  string  `json:"status" validate:"oneof=to_watch watching completed dropped"`
	Rating  float64 `json:"rating" validate:"min=0,max=10"`
	Notes   string  `json:"notes"`
}

// WatchlistItemUpdateRequest represents a request to update a watchlist item
type WatchlistItemUpdateRequest struct {
	Status string  `json:"status" validate:"oneof=to_watch watching completed dropped"`
	Rating float64 `json:"rating" validate:"min=0,max=10"`
	Notes  string  `json:"notes"`
}

// WatchlistStats represents statistics for a watchlist
type WatchlistStats struct {
	TotalItems     int     `json:"total_items"`
	CompletedItems int     `json:"completed_items"`
	WatchingItems  int     `json:"watching_items"`
	ToWatchItems   int     `json:"to_watch_items"`
	DroppedItems   int     `json:"dropped_items"`
	AverageRating  float64 `json:"average_rating"`
	TotalHours     int     `json:"total_hours"`
}

// WatchlistRecommendation represents a recommendation based on watchlist
type WatchlistRecommendation struct {
	Movie       Movie   `json:"movie"`
	Score       float64 `json:"score"`
	Reason      string  `json:"reason"`
	GenreMatch  float64 `json:"genre_match"`
	RatingMatch float64 `json:"rating_match"`
}
