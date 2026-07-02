package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/database"
	"github.com/todo/backend/go/internal/handler"
	"github.com/todo/backend/go/internal/repository"
	"github.com/todo/backend/go/internal/router"
	"github.com/todo/backend/go/internal/service"
)

// @title           Todo API
// @version         1.0
// @description     RESTful API for managing todo items
// @host            localhost:8080
// @BasePath        /

func main() {
	// ── Configuration (appsettings.json / IConfiguration) ────────────────────
	cfg := config.Load()

	// ── Database (DbContext) ──────────────────────────────────────────────────
	db, err := database.New(cfg.DatabaseDSN)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// ── Dependency injection (composition root - mirrors Program.cs AddScoped/AddSingleton) ──
	repo := repository.NewTodoItemRepository(db)
	svc := service.NewTodoItemService(repo)
	h := handler.NewTodoItemHandler(svc)

	fileRepo := repository.NewFileRepository(db)
	fileSvc := service.NewFileService(fileRepo, cfg)
	fh := handler.NewFileHandler(fileSvc)

	// ── Gin engine (analogous to app.Build() + middleware pipeline) ───────────
	r := gin.Default()

	// Swagger UI - mirrors app.UseSwagger() + app.UseSwaggerUI()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ── Routes ────────────────────────────────────────────────────────────────
	router.Setup(r, h, fh)

	log.Printf("Starting %s v%s on :%s", cfg.AppName, cfg.AppVersion, cfg.Port)
	log.Printf("Swagger UI: http://localhost:%s/swagger/index.html", cfg.Port)

	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
