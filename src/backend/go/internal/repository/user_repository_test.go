package repository_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestUserRepositoryQueriesAreCaseInsensitive(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil && strings.Contains(err.Error(), "requires cgo") {
		t.Skip("SQLite-backed repository tests require CGO")
	}
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&models.User{}, &models.EmailLog{}))
	repo := repository.NewUserRepository(db)

	created, err := repo.Create(&models.User{Username: "Alice", Email: "Alice@Example.com", PasswordHash: "hash", IsActive: true})
	require.NoError(t, err)
	byEmail, err := repo.FindByEmail("alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, created.ID, byEmail.ID)

	exists, err := repo.UsernameExists("alice", nil)
	require.NoError(t, err)
	assert.True(t, exists)
	exists, err = repo.EmailExists("ALICE@EXAMPLE.COM", &created.ID)
	require.NoError(t, err)
	assert.False(t, exists)
}
