package config

import (
	"os"
)

type Config struct {
	ServerPort    string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	DBSSLRootCert string
	JWTSecret     string
}

// LoadConfig loads configuration from environment variables or defaults
func LoadConfig() *Config {
	return &Config{
		ServerPort:    getEnv("SERVER_PORT", "3000"),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "postgres"),
		DBPassword:    getEnv("DB_PASSWORD", ""),
		DBName:        getEnv("DB_NAME", "book_exchange"),
		DBPort:        getEnv("DB_PORT", "5432"),
		DBSSLRootCert: getEnv("DB_SSL_ROOT_CERT", "ca.pem"), // Path to SSL certificate
		JWTSecret:     getEnv("JWT_SECRET", ""),
	}
}

// getEnv retrieves environment variables with a fallback default value
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
