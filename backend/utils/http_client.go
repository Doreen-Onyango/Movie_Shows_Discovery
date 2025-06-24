package utils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/config"
)

// HTTPClient represents a custom HTTP client with retry and rate limiting
type HTTPClient struct {
	client    *http.Client
	rateLimit *RateLimiter
}

// RateLimiter handles rate limiting for API calls
type RateLimiter struct {
	requests chan struct{}
	ticker   *time.Ticker
}

// NewHTTPClient creates a new HTTP client with custom configuration
func NewHTTPClient(timeout time.Duration, requestsPerSecond int) *HTTPClient {
	client := &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 10,
			IdleConnTimeout:     90 * time.Second,
		},
	}

	rateLimiter := NewRateLimiter(requestsPerSecond)

	return &HTTPClient{
		client:    client,
		rateLimit: rateLimiter,
	}
}

// NewRateLimiter creates a new rate limiter
func NewRateLimiter(requestsPerSecond int) *RateLimiter {
	interval := time.Second / time.Duration(requestsPerSecond)
	ticker := time.NewTicker(interval)
	requests := make(chan struct{}, requestsPerSecond)

	go func() {
		for range ticker.C {
			select {
			case requests <- struct{}{}:
			default:
			}
		}
	}()

	return &RateLimiter{
		requests: requests,
		ticker:   ticker,
	}
}

// Get performs a GET request with retry logic
func (h *HTTPClient) Get(ctx context.Context, url string, headers map[string]string) (*http.Response, error) {
	return h.RequestWithRetry(ctx, "GET", url, nil, headers)
}

// Post performs a POST request with retry logic
func (h *HTTPClient) Post(ctx context.Context, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	return h.RequestWithRetry(ctx, "POST", url, body, headers)
}

// RequestWithRetry performs an HTTP request with exponential backoff retry
func (h *HTTPClient) RequestWithRetry(ctx context.Context, method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	maxRetries := 3
	backoff := time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// Wait for rate limiter
		select {
		case <-h.rateLimit.requests:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		// Create request
		req, err := http.NewRequestWithContext(ctx, method, url, body)
		if err != nil {
			return nil, fmt.Errorf("failed to create request: %w", err)
		}

		// Add headers
		for key, value := range headers {
			req.Header.Set(key, value)
		}

		// Perform request
		resp, err := h.client.Do(req)
		if err != nil {
			if attempt == maxRetries {
				return nil, fmt.Errorf("request failed after %d attempts: %w", maxRetries+1, err)
			}
			time.Sleep(backoff)
			backoff *= 2
			continue
		}

		// Check if response is successful
		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return resp, nil
		}

		// Handle specific error codes
		switch resp.StatusCode {
		case 429: // Rate limited
			if attempt == maxRetries {
				return resp, fmt.Errorf("rate limited after %d attempts", maxRetries+1)
			}
			// Wait longer for rate limit
			time.Sleep(backoff * 2)
			backoff *= 2
			resp.Body.Close()
			continue
		case 500, 502, 503, 504: // Server errors
			if attempt == maxRetries {
				return resp, fmt.Errorf("server error after %d attempts", maxRetries+1)
			}
			time.Sleep(backoff)
			backoff *= 2
			resp.Body.Close()
			continue
		default:
			// For other status codes, don't retry
			return resp, nil
		}
	}

	return nil, fmt.Errorf("request failed after %d attempts", maxRetries+1)
}

// Close closes the HTTP client and cleans up resources
func (h *HTTPClient) Close() {
	if h.rateLimit.ticker != nil {
		h.rateLimit.ticker.Stop()
	}
}

// CreateTMDBClient creates an HTTP client configured for TMDB API
func CreateTMDBClient() *HTTPClient {
	cfg := config.AppConfig
	return NewHTTPClient(30*time.Second, cfg.TMDB.RateLimit)
}

// CreateOMDBClient creates an HTTP client configured for OMDB API
func CreateOMDBClient() *HTTPClient {
	cfg := config.AppConfig
	return NewHTTPClient(30*time.Second, cfg.OMDB.RateLimit)
}

// CreateDefaultClient creates a default HTTP client
func CreateDefaultClient() *HTTPClient {
	return NewHTTPClient(30*time.Second, 100)
}
