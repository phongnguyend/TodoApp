package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/handler"
)

// Setup registers all routes on the Gin engine.
// Mirrors app.MapControllers() / minimal-API route group configuration in ASP.NET Core.
func Setup(r *gin.Engine, h *handler.TodoItemHandler, fh *handler.FileHandler, attachmentHandlers ...*handler.TodoItemAttachmentHandler) {
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
		api.POST("/import/excel", h.ImportExcel)
		api.GET("/export/excel", h.ExportExcel)
	}
	if len(attachmentHandlers) > 0 && attachmentHandlers[0] != nil {
		ah := attachmentHandlers[0]
		api.GET("/:id/attachments", ah.GetAll)
		api.POST("/:id/attachments", ah.Create)
		api.GET("/:id/attachments/:attachmentId", ah.GetByID)
		api.PUT("/:id/attachments/:attachmentId", ah.Update)
		api.DELETE("/:id/attachments/:attachmentId", ah.Delete)
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

// RegisterUsers registers user-management and authenticated self-service routes.
func RegisterUsers(r *gin.Engine, h *handler.UserHandler) {
	users := r.Group("/api/users")
	users.POST("/signup", h.SignUp)
	users.POST("/password/change", h.ChangePassword)
	users.POST("/password/reset", h.RequestPasswordReset)
	users.POST("/password/confirm", h.ConfirmPasswordReset)
	users.GET("/profile", h.GetProfile)
	users.PUT("/profile", h.UpdateProfile)
	users.GET("", h.GetAll)
	users.POST("", h.Create)
	users.GET("/:id", h.GetByID)
	users.PUT("/:id", h.Update)
	users.PATCH("/:id/activate", h.Activate)
	users.PATCH("/:id/deactivate", h.Deactivate)
}
