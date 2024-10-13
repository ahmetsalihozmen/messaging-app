package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds the configuration variables
type Config struct {
	WebhookURL      string
	RedisAddr       string
	Port            string
	MessagingPeriod int
	DatabaseURL     string
}

func LoadConfig() Config {
	godotenv.Load()
	cfg := Config{
		WebhookURL:      getEnv("WEBHOOK_URL", ""),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		Port:            getEnv("PORT", "8080"),
		MessagingPeriod: getEnvInt("MESSAGING_PERIOD", 120), // Default is 2 minutes
		DatabaseURL:     getEnv("DATABASE_URL", "messages.db"),
	}
	return cfg
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
