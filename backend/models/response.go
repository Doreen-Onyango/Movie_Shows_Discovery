package models

import "time"

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

// Meta represents metadata for paginated responses
type Meta struct {
	Page         int  `json:"page"`
	PerPage      int  `json:"per_page"`
	TotalPages   int  `json:"total_pages"`
	TotalResults int  `json:"total_results"`
	HasNext      bool `json:"has_next"`
	HasPrev      bool `json:"has_prev"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
	Message   string    `json:"message"`
	Code      int       `json:"code"`
	Timestamp time.Time `json:"timestamp"`
	Path      string    `json:"path,omitempty"`
}

// SuccessResponse represents a success response
type SuccessResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Timestamp time.Time   `json:"timestamp"`
	Meta      *Meta       `json:"meta,omitempty"`
}

// PaginatedResponse represents a paginated response
type PaginatedResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data"`
	Meta      Meta        `json:"meta"`
	Timestamp time.Time   `json:"timestamp"`
}

// HealthResponse represents a health check response
type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
	Version   string            `json:"version"`
}

// SearchResponse represents a search response
type SearchResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Results   interface{} `json:"results"`
	Query     string      `json:"query"`
	Meta      Meta        `json:"meta"`
	Timestamp time.Time   `json:"timestamp"`
}

// TrendingResponse represents a trending content response
type TrendingResponse struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Trending  interface{} `json:"trending"`
	Timeframe string      `json:"timeframe"`
	Meta      Meta        `json:"meta"`
	Timestamp time.Time   `json:"timestamp"`
}

// RecommendationResponse represents a recommendation response
type RecommendationResponse struct {
	Success         bool        `json:"success"`
	Message         string      `json:"message"`
	Recommendations interface{} `json:"recommendations"`
	UserID          string      `json:"user_id"`
	Meta            Meta        `json:"meta"`
	Timestamp       time.Time   `json:"timestamp"`
}

// Helper functions for creating responses
func NewSuccessResponse(data interface{}, message string) SuccessResponse {
	return SuccessResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Timestamp: time.Now(),
	}
}

func NewErrorResponse(err string, message string, code int) ErrorResponse {
	return ErrorResponse{
		Success:   false,
		Error:     err,
		Message:   message,
		Code:      code,
		Timestamp: time.Now(),
	}
}

func NewPaginatedResponse(data interface{}, meta Meta, message string) PaginatedResponse {
	return PaginatedResponse{
		Success:   true,
		Message:   message,
		Data:      data,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}

func NewSearchResponse(results interface{}, query string, meta Meta) SearchResponse {
	return SearchResponse{
		Success:   true,
		Message:   "Search completed successfully",
		Results:   results,
		Query:     query,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}

func NewTrendingResponse(trending interface{}, timeframe string, meta Meta) TrendingResponse {
	return TrendingResponse{
		Success:   true,
		Message:   "Trending content retrieved successfully",
		Trending:  trending,
		Timeframe: timeframe,
		Meta:      meta,
		Timestamp: time.Now(),
	}
}

func NewRecommendationResponse(recommendations interface{}, userID string, meta Meta) RecommendationResponse {
	return RecommendationResponse{
		Success:         true,
		Message:         "Recommendations generated successfully",
		Recommendations: recommendations,
		UserID:          userID,
		Meta:            meta,
		Timestamp:       time.Now(),
	}
}
