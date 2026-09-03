package config

import (
	"os"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	AppPort    string
	DBUser     string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	JWTSecret  string
}

func LoadEnv() *EnvConfig {
	// Attempt to load .env file. We ignore the error because in production,
	// environment variables might be injected directly by the host (like Docker).
	_ = godotenv.Load()

	return &EnvConfig{
		AppPort:    getEnvOrDefault("APP_PORT", "8080"),
		DBUser:     getEnvOrPanic("DB_USER"),
		DBPassword: getEnvOrPanic("DB_PASSWORD"),
		DBHost:     getEnvOrPanic("DB_HOST"),
		DBPort:     getEnvOrPanic("DB_PORT"),
		DBName:     getEnvOrPanic("DB_NAME"),
		JWTSecret:  getEnvOrPanic("JWT_SECRET"),
	}
}

// getEnvOrPanic ensures critical configuration is present (Fail-Fast)
func getEnvOrPanic(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists {
		panic("CRITICAL: Missing required environment variable: " + key)
	}
	return val
}

// getEnvOrDefault returns a fallback if the env var is not set
func getEnvOrDefault(key, fallback string) string {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}
	return val
}
