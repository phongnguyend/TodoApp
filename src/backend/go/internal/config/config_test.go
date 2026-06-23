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
	unsetEnv(t, "APP_NAME", "APP_VERSION", "PORT", "DATABASE_DSN", "DEFAULT_PAGE_SIZE", "MAX_PAGE_SIZE")

	cfg := config.Load()

	assert.Equal(t, "Todo API", cfg.AppName)
	assert.Equal(t, "1.0.0", cfg.AppVersion)
	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "todo.db", cfg.DatabaseDSN)
	assert.Equal(t, 20, cfg.DefaultPageSize)
	assert.Equal(t, 100, cfg.MaxPageSize)
}

func TestLoad_EnvVarsOverrideDefaults(t *testing.T) {
	t.Setenv("APP_NAME", "My App")
	t.Setenv("APP_VERSION", "2.0.0")
	t.Setenv("PORT", "9090")
	t.Setenv("DATABASE_DSN", "test.db")
	t.Setenv("DEFAULT_PAGE_SIZE", "50")
	t.Setenv("MAX_PAGE_SIZE", "200")

	cfg := config.Load()

	assert.Equal(t, "My App", cfg.AppName)
	assert.Equal(t, "2.0.0", cfg.AppVersion)
	assert.Equal(t, "9090", cfg.Port)
	assert.Equal(t, "test.db", cfg.DatabaseDSN)
	assert.Equal(t, 50, cfg.DefaultPageSize)
	assert.Equal(t, 200, cfg.MaxPageSize)
}

func TestLoad_InvalidPageSizeEnvVar_UsesDefault(t *testing.T) {
	unsetEnv(t, "DEFAULT_PAGE_SIZE", "MAX_PAGE_SIZE")
	t.Setenv("DEFAULT_PAGE_SIZE", "not-a-number")
	t.Setenv("MAX_PAGE_SIZE", "abc")

	cfg := config.Load()

	assert.Equal(t, 20, cfg.DefaultPageSize)
	assert.Equal(t, 100, cfg.MaxPageSize)
}
