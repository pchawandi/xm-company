package config

import (
	"log"
	"os"
)

// DatabaseConfig holds the database connection settings
type DatabaseConfig struct {
	Hostname string
	Port     string
	DBName   string
	User     string
	Password string
}

// LoadDatabaseConfig loads database-related environment variables and returns a DatabaseConfig
func LoadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Hostname: getEnv("POSTGRES_HOST", "localhost"),
		Port:     getEnv("POSTGRES_PORT", "5432"),
		DBName:   getEnv("POSTGRES_DB", "companies"),
		User:     getEnv("POSTGRES_USER", "postgres"),
		Password: getEnv("POSTGRES_PASSWORD", "postgres"),
	}
}

// Helper function to get environment variables with a fallback default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Printf("Environment variable %s not set. Using default value: %s", key, defaultValue)
		return defaultValue
	}
	return value
}
