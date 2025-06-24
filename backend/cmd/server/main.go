package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/config"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/middleware"
	"github.com/Doreen-Onyango/Movie_Shows_Discovery/backend/services"
	"github.com/joho/godotenv"
)

func main() {
	// Load configuration
	godotenv.Load()
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize logger
	logger := middleware.NewLogger()

	// Initialize services
	tmdbService := services.NewTMDBService()
	omdbService := services.NewOMDBService()

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", config.AppConfig.Server.Port),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		logger.InfoLogger.Printf("Starting server on port %s", config.AppConfig.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.ErrorLogger.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.InfoLogger.Println("Shutting down server...")

	// Create a deadline for server shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Attempt graceful shutdown
	if err := server.Shutdown(ctx); err != nil {
		logger.ErrorLogger.Printf("Server forced to shutdown: %v", err)
	}

	// Clean up services
	tmdbService.Close()
	omdbService.Close()

	logger.InfoLogger.Println("Server exited")
}

// init function to set up logging
func init() {
	// Set up basic logging
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("Movie Shows Discovery Backend starting...")
}
