// Package config loads and validates the application configuration.
//
// Configuration is loaded from environment variables.
// During development, variables are read from the .env file.
// In production, values are expected to be provided by the operating system.
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config contains the runtime configuration shared across the application.
type Config struct {
	AppEnv string
	Port   string

	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string

	JWTSecret string
}

// Load reads the application configuration.
//
// If a .env file exists, it is loaded automatically.
// Missing variables are replaced with sensible defaults when possible.
func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using system environment variables")
	}

	return Config{
		AppEnv: getEnv("APP_ENV", "development"),
		Port:   getEnv("APP_PORT", "8080"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "booking"),
		DBUser:     getEnv("DB_USER", "booking"),
		DBPassword: getEnv("DB_PASSWORD", "booking"),

		JWTSecret: getEnv("JWT_SECRET", "change-me"),
	}
}

// getEnv returns the environment variable value if present.
// Otherwise, it falls back to the provided default value.
func getEnv(key, fallback string) string {
	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}
