package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application settings — analogous to appsettings.json bound via IConfiguration.
type Config struct {
	AppName     string
	AppVersion  string
	Port        string
	DatabaseDSN string
	DefaultPageSize int
	MaxPageSize     int
}

// Load reads settings from the environment (and optionally a .env file).
// Mirrors the host builder configuration in Program.cs.
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment")
	}

	return &Config{
		AppName:         getEnv("APP_NAME", "Todo API"),
		AppVersion:      getEnv("APP_VERSION", "1.0.0"),
		Port:            getEnv("PORT", "8080"),
		DatabaseDSN:     getEnv("DATABASE_DSN", "todo.db"),
		DefaultPageSize: getEnvInt("DEFAULT_PAGE_SIZE", 20),
		MaxPageSize:     getEnvInt("MAX_PAGE_SIZE", 100),
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return fallback
}
