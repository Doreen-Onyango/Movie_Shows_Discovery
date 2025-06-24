package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	TMDB     TMDBConfig
	OMDB     OMDBConfig
	Cache    CacheConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type TMDBConfig struct {
	APIKey    string
	BaseURL   string
	RateLimit int
}

type OMDBConfig struct {
	APIKey    string
	BaseURL   string
	RateLimit int
}

type CacheConfig struct {
	TTL         time.Duration
	SearchTTL   time.Duration
	TrendingTTL time.Duration
}

type LoggingConfig struct {
	Level string
}

var AppConfig *Config

func LoadConfig() error {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		// Don't fail if .env doesn't exist, use system env vars
		fmt.Println("Warning: .env file not found, using system environment variables")
	}

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("PORT", "8080"),
			Host: getEnv("HOST", "localhost"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			Name:     getEnv("DB_NAME", "movie_discovery"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		TMDB: TMDBConfig{
			APIKey:    getEnv("TMDB_API_KEY", ""),
			BaseURL:   getEnv("TMDB_BASE_URL", "https://api.themoviedb.org/3"),
			RateLimit: getEnvAsInt("TMDB_RATE_LIMIT", 40),
		},
		OMDB: OMDBConfig{
			APIKey:    getEnv("OMDB_API_KEY", ""),
			BaseURL:   getEnv("OMDB_BASE_URL", "http://www.omdbapi.com"),
			RateLimit: getEnvAsInt("OMDB_RATE_LIMIT", 1000),
		},
		Cache: CacheConfig{
			TTL:         getEnvAsDuration("CACHE_TTL", 3600),
			SearchTTL:   getEnvAsDuration("SEARCH_CACHE_TTL", 1800),
			TrendingTTL: getEnvAsDuration("TRENDING_CACHE_TTL", 3600),
		},
		Logging: LoggingConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}

	// Validate required configuration
	if AppConfig.TMDB.APIKey == "" {
		return fmt.Errorf("TMDB_API_KEY is required")
	}
	if AppConfig.OMDB.APIKey == "" {
		return fmt.Errorf("OMDB_API_KEY is required")
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValueSeconds int) time.Duration {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return time.Duration(intValue) * time.Second
		}
	}
	return time.Duration(defaultValueSeconds) * time.Second
}
