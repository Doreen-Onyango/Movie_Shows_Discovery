package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/models"
)

// Logger represents a custom logger
type Logger struct {
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
}

// NewLogger creates a new logger instance
func NewLogger() *Logger {
	return &Logger{
		InfoLogger:  log.New(log.Writer(), "INFO: ", log.LstdFlags),
		ErrorLogger: log.New(log.Writer(), "ERROR: ", log.LstdFlags),
	}
}

// LoggingMiddleware logs HTTP requests with timing and status
func LoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Create a custom response writer to capture status code
			responseWriter := &ResponseWriter{
				ResponseWriter: w,
				StatusCode:     http.StatusOK,
			}

			// Process request
			next.ServeHTTP(responseWriter, r)

			// Calculate duration
			duration := time.Since(start)

			// Log request details
			logger.InfoLogger.Printf(
				"%s %s %s %d %v %s",
				r.Method,
				r.URL.Path,
				r.RemoteAddr,
				responseWriter.StatusCode,
				duration,
				r.UserAgent(),
			)
		})
	}
}

// ResponseWriter wraps http.ResponseWriter to capture status code
type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.StatusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// ErrorLoggingMiddleware logs errors with additional context
func ErrorLoggingMiddleware(logger *Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if err := recover(); err != nil {
					logger.ErrorLogger.Printf(
						"PANIC: %s %s %s %v",
						r.Method,
						r.URL.Path,
						r.RemoteAddr,
						err,
					)

					// Return 500 Internal Server Error
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

// CORSMiddleware handles Cross-Origin Resource Sharing
func CORSMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Set CORS headers
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Max-Age", "86400")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware implements basic rate limiting
func RateLimitMiddleware(requestsPerMinute int) func(http.Handler) http.Handler {
	// Simple in-memory rate limiter
	clients := make(map[string][]time.Time)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := r.RemoteAddr
			now := time.Now()

			// Clean old requests
			if times, exists := clients[clientIP]; exists {
				var validTimes []time.Time
				for _, t := range times {
					if now.Sub(t) < time.Minute {
						validTimes = append(validTimes, t)
					}
				}
				clients[clientIP] = validTimes
			}

			// Check rate limit
			if times, exists := clients[clientIP]; exists && len(times) >= requestsPerMinute {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			// Add current request
			clients[clientIP] = append(clients[clientIP], now)

			next.ServeHTTP(w, r)
		})
	}
}

// RequestIDMiddleware adds a unique request ID to each request
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := generateRequestID()
			r.Header.Set("X-Request-ID", requestID)
			w.Header().Set("X-Request-ID", requestID)

			next.ServeHTTP(w, r)
		})
	}
}

// generateRequestID generates a simple request ID
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" + randomString(8)
}

// randomString generates a random string of specified length
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[time.Now().UnixNano()%int64(len(charset))]
	}
	return string(b)
}

// LogError logs an error with context
func (l *Logger) LogError(err error, context string, r *http.Request) {
	l.ErrorLogger.Printf(
		"ERROR: %s - %s %s %s %v",
		context,
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		err,
	)
}

// LogInfo logs an info message
func (l *Logger) LogInfo(message string, r *http.Request) {
	l.InfoLogger.Printf(
		"INFO: %s - %s %s %s",
		message,
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
	)
}

// LogAPIResponse logs API response details
func (l *Logger) LogAPIResponse(response models.APIResponse, r *http.Request) {
	status := "SUCCESS"
	if !response.Success {
		status = "ERROR"
	}

	l.InfoLogger.Printf(
		"API_RESPONSE: %s - %s %s %s - %s",
		status,
		r.Method,
		r.URL.Path,
		r.RemoteAddr,
		response.Message,
	)
}
