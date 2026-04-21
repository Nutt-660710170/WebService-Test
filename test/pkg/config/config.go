package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config holds all application configuration.
type Config struct {
	DBHost         string
	DBPort         string
	DBUser         string
	DBPassword     string
	DBName         string
	AppPort        string
	BangchakAPIURL string
}

// Load reads configuration from .env file and environment variables.
func Load() *Config {
	// Load .env file (ignore error if file doesn't exist)
	_ = godotenv.Load()

	return &Config{
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "oilprice"),
		DBPassword:     getEnv("DB_PASSWORD", "oilprice123"),
		DBName:         getEnv("DB_NAME", "oilprice_db"),
		AppPort:        getEnv("APP_PORT", "8080"),
		BangchakAPIURL: getEnv("BANGCHAK_API_URL", "https://oil-price.bangchak.co.th/ApiOilPrice2/en"),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
