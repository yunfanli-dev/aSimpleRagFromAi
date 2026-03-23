package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	AppName         string
	AppEnv          string
	HTTPAddr        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	ShutdownTimeout time.Duration
	PostgresDSN     string
	RedisAddr       string
	RedisPassword   string
	RedisDB         string
	EmbeddingModel  string
	EmbeddingDims   int
	LLMProvider     string
	LLMModel        string
	MiniMaxAPIKey   string
	MiniMaxBaseURL  string
	LLMTimeout      time.Duration
}

// Load reads process configuration from environment variables with defaults.
func Load() Config {
	return Config{
		AppName:         getEnv("APP_NAME", "simplerag-go"),
		AppEnv:          getEnv("APP_ENV", "dev"),
		HTTPAddr:        getEnv("HTTP_ADDR", ":8080"),
		ReadTimeout:     getDurationEnv("READ_TIMEOUT", 5*time.Second),
		WriteTimeout:    getDurationEnv("WRITE_TIMEOUT", 10*time.Second),
		ShutdownTimeout: getDurationEnv("SHUTDOWN_TIMEOUT", 10*time.Second),
		PostgresDSN:     getEnv("POSTGRES_DSN", "postgres://postgres:postgres@localhost:5432/simplerag?sslmode=disable"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   getEnv("REDIS_PASSWORD", ""),
		RedisDB:         getEnv("REDIS_DB", "0"),
		EmbeddingModel:  getEnv("EMBEDDING_MODEL", "local-hash-v1"),
		EmbeddingDims:   getIntEnv("EMBEDDING_DIMS", 1024),
		LLMProvider:     getEnv("LLM_PROVIDER", "local"),
		LLMModel:        getEnv("LLM_MODEL", "local-extractive-v1"),
		MiniMaxAPIKey:   getEnv("MINIMAX_API_KEY", ""),
		MiniMaxBaseURL:  getEnv("MINIMAX_BASE_URL", "https://api.minimax.io/v1"),
		LLMTimeout:      getDurationEnv("LLM_TIMEOUT", 30*time.Second),
	}
}

// getEnv reads a string environment variable with fallback.
func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// getDurationEnv reads a duration environment variable with fallback on parse failure.
func getDurationEnv(key string, fallback time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}

// getIntEnv reads an integer environment variable with fallback on parse failure.
func getIntEnv(key string, fallback int) int {
	raw := os.Getenv(key)
	if raw == "" {
		return fallback
	}

	value, err := strconv.Atoi(raw)
	if err != nil {
		return fallback
	}
	return value
}
