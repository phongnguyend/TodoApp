package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/handler"
)

// Setup registers all routes on the Gin engine.
// Mirrors app.MapControllers() / minimal-API route group configuration in ASP.NET Core.
func Setup(r *gin.Engine, h *handler.TodoItemHandler) {
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
	}
}
