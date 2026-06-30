package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config holds all application settings — analogous to appsettings.json bound via IConfiguration.
type Config struct {
	AppName         string
	AppVersion      string
	Port            string
	DatabaseDSN     string
	DefaultPageSize int
	MaxPageSize     int

	// SMTP / email settings (used by the background worker)
	SMTPHost       string
	SMTPPort       int
	SMTPUsername   string
	SMTPPassword   string
	SMTPUseTLS     bool
	EmailSender    string
	EmailRecipient string

	// Worker settings
	WorkerIntervalMinutes int
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

		SMTPHost:              getEnv("SMTP_HOST", "localhost"),
		SMTPPort:              getEnvInt("SMTP_PORT", 587),
		SMTPUsername:          getEnv("SMTP_USERNAME", ""),
		SMTPPassword:          getEnv("SMTP_PASSWORD", ""),
		SMTPUseTLS:            getEnvBool("SMTP_USE_TLS", true),
		EmailSender:           getEnv("EMAIL_SENDER", "noreply@example.com"),
		EmailRecipient:        getEnv("EMAIL_RECIPIENT", "admin@example.com"),
		WorkerIntervalMinutes: getEnvInt("WORKER_INTERVAL_MINUTES", 60),
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

func getEnvBool(key string, fallback bool) bool {
	if v, ok := os.LookupEnv(key); ok {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return fallback
}
