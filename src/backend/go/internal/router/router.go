package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/handler"
)

// Setup registers all routes on the Gin engine.
// Mirrors app.MapControllers() / minimal-API route group configuration in ASP.NET Core.
func Setup(r *gin.Engine, h *handler.TodoItemHandler, fh *handler.FileHandler) {
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	api := r.Group("/api/todo-items")
	{
		api.GET("", h.GetAll)
		api.GET("/incomplete", h.GetIncomplete)
		api.GET("/:id", h.GetByID)
		api.POST("", h.Create)
		api.PUT("/:id", h.Update)
		api.PATCH("/:id/complete", h.MarkComplete)
		api.DELETE("/:id", h.Delete)
		api.POST("/import/csv", h.ImportCSV)
		api.GET("/export/csv", h.ExportCSV)
	}

	files := r.Group("/api/files")
	{
		files.GET("", fh.GetAll)
		files.GET("/:id", fh.GetByID)
		files.GET("/:id/download", fh.Download)
		files.POST("", fh.Upload)
		files.DELETE("/:id", fh.Delete)
	}
}
