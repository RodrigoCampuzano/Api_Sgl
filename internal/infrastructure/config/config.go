package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// Database
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	// Security
	JWTSecretKey          string
	JWTExpirationHours    int
	SessionTimeoutMinutes int

	// Server
	Port    string
	GinMode string

	// Storage
	StorageProvider string
	StoragePath     string

	// Logging
	LogLevel string
}

func Load() *Config {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := &Config{
		// Database
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "sgl_user"),
		DBPassword: getEnv("DB_PASSWORD", "secure_password"),
		DBName:     getEnv("DB_NAME", "sgl_disasur"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		// Security
		JWTSecretKey:          getEnv("JWT_SECRET_KEY", "default-secret-key-CHANGE-IN-PRODUCTION"),
		JWTExpirationHours:    getEnvAsInt("JWT_EXPIRATION_HOURS", 8),
		SessionTimeoutMinutes: getEnvAsInt("SESSION_TIMEOUT_MINUTES", 30),

		// Server
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),

		// Storage
		StorageProvider: getEnv("STORAGE_PROVIDER", "local"),
		StoragePath:     getEnv("STORAGE_PATH", "./uploads"),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	return cfg
}

func (c *Config) GetDatabaseURL() string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
		c.DBSSLMode,
	)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}

	return value
}
