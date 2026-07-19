package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/todo/backend/go/internal/config"
)

// unsetEnv temporarily clears env vars for the duration of a test and restores
// their previous values via t.Cleanup.
func unsetEnv(t *testing.T, keys ...string) {
	t.Helper()
	for _, k := range keys {
		k := k
		prev, exists := os.LookupEnv(k)
		os.Unsetenv(k)
		if exists {
			t.Cleanup(func() { os.Setenv(k, prev) })
		} else {
			t.Cleanup(func() { os.Unsetenv(k) })
		}
	}
}

func TestLoad_DefaultValues(t *testing.T) {
	unsetEnv(t, "APP_NAME", "APP_VERSION", "PORT", "DATABASE_DSN", "DEFAULT_PAGE_SIZE", "MAX_PAGE_SIZE", "FILE_STORAGE_PATH", "MAX_UPLOAD_SIZE_BYTES", "JWT_SECRET_KEY", "PASSWORD_HASH_ITERATIONS", "PASSWORD_RESET_SECRET_KEY", "PASSWORD_RESET_TOKEN_LIFETIME_MINUTES", "PASSWORD_RESET_CONFIRMATION_URL")

	cfg := config.Load()

	assert.Equal(t, "Todo API", cfg.AppName)
	assert.Equal(t, "1.0.0", cfg.AppVersion)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "todo.db", cfg.DatabaseDSN)
	assert.Equal(t, 20, cfg.DefaultPageSize)
	assert.Equal(t, 100, cfg.MaxPageSize)
	assert.Equal(t, "./uploads", cfg.FileStoragePath)
	assert.Equal(t, int64(10485760), cfg.MaxUploadSizeBytes)
	assert.Equal(t, "change-me-use-at-least-32-bytes-long", cfg.JWTSecretKey)
	assert.Equal(t, 120000, cfg.PasswordHashIterations)
	assert.Equal(t, "change-me-use-at-least-32-bytes-long", cfg.PasswordResetSecretKey)
	assert.Equal(t, 60, cfg.PasswordResetTokenLifetimeMinutes)
	assert.Equal(t, "/reset-password", cfg.PasswordResetConfirmationURL)
}

func TestLoad_EnvVarsOverrideDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "My App")
	t.Setenv("APP_VERSION", "2.0.0")
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_DSN", "test.db")
	t.Setenv("DEFAULT_PAGE_SIZE", "50")
	t.Setenv("MAX_PAGE_SIZE", "200")
	t.Setenv("FILE_STORAGE_PATH", "/data/uploads")
	t.Setenv("MAX_UPLOAD_SIZE_BYTES", "2048")
	t.Setenv("JWT_SECRET_KEY", "jwt-test")
	t.Setenv("PASSWORD_HASH_ITERATIONS", "1000")
	t.Setenv("PASSWORD_RESET_SECRET_KEY", "reset-test")
	t.Setenv("PASSWORD_RESET_TOKEN_LIFETIME_MINUTES", "30")
	t.Setenv("PASSWORD_RESET_CONFIRMATION_URL", "https://example.com/reset")

	cfg := config.Load()

	assert.Equal(t, "My App", cfg.AppName)
	assert.Equal(t, "2.0.0", cfg.AppVersion)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "test.db", cfg.DatabaseDSN)
	assert.Equal(t, 50, cfg.DefaultPageSize)
	assert.Equal(t, 200, cfg.MaxPageSize)
	assert.Equal(t, "/data/uploads", cfg.FileStoragePath)
	assert.Equal(t, int64(2048), cfg.MaxUploadSizeBytes)
	assert.Equal(t, "jwt-test", cfg.JWTSecretKey)
	assert.Equal(t, 1000, cfg.PasswordHashIterations)
	assert.Equal(t, "reset-test", cfg.PasswordResetSecretKey)
	assert.Equal(t, 30, cfg.PasswordResetTokenLifetimeMinutes)
	assert.Equal(t, "https://example.com/reset", cfg.PasswordResetConfirmationURL)
}

func TestLoad_InvalidPageSizeEnvVar_UsesDefault(t *testing.T) {
	unsetEnv(t, "DEFAULT_PAGE_SIZE", "MAX_PAGE_SIZE")
	t.Setenv("DEFAULT_PAGE_SIZE", "not-a-number")
	t.Setenv("MAX_PAGE_SIZE", "abc")

	cfg := config.Load()

	assert.Equal(t, 20, cfg.DefaultPageSize)
	assert.Equal(t, 100, cfg.MaxPageSize)
}

func TestLoad_InvalidMaxUploadSizeEnvVar_UsesDefault(t *testing.T) {
	unsetEnv(t, "MAX_UPLOAD_SIZE_BYTES")
	t.Setenv("MAX_UPLOAD_SIZE_BYTES", "not-a-number")

	cfg := config.Load()

	assert.Equal(t, int64(10485760), cfg.MaxUploadSizeBytes)
}
