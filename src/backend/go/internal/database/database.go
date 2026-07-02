package database

import (
	"log"

	"github.com/todo/backend/go/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// New opens a GORM database connection and runs AutoMigrate.
// This is analogous to registering DbContext and calling EnsureCreated / Migrate
// in Program.cs. Swap the driver import to use PostgreSQL or MySQL.
func New(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// AutoMigrate creates/updates tables to match the model structs.
	// For production, use a dedicated migration tool (golang-migrate, goose, etc.)
	// instead of AutoMigrate - similar to running `dotnet ef database update`.
	if err := db.AutoMigrate(&models.TodoItem{}, &models.EmailLog{}, &models.File{}); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated")
	return db, nil
}
