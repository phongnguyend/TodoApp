package handler

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/service"
)

// FileHandler handles HTTP requests for uploaded files - analogous to a
// [ApiController] mapped to /api/files in ASP.NET Core.
type FileHandler struct {
	svc service.FileService
}

// NewFileHandler creates the handler with its service dependency injected.
func NewFileHandler(svc service.FileService) *FileHandler {
	return &FileHandler{svc: svc}
}

func handleFileServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrFileNotFound), errors.Is(err, service.ErrFileContentMissing):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, service.ErrFileTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

// ── Endpoints ─────────────────────────────────────────────────────────────────

// GetAll godoc
// @Summary      List all uploaded files
// @Tags         Files
// @Produce      json
// @Param        page     query  int  false  "Page number"   default(1)
// @Param        pageSize query  int  false  "Page size"     default(20)
// @Success      200  {object}  dto.PaginatedResponse[dto.FileResponse]
// @Router       /api/files [get]
func (h *FileHandler) GetAll(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)

	result, err := h.svc.GetAll(page, pageSize)
	if err != nil {
		handleFileServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetByID godoc
// @Summary      Get a file's metadata by ID
// @Tags         Files
// @Produce      json
// @Param        id   path  int  true  "File ID"
// @Success      200  {object}  dto.FileResponse
// @Failure      404  {object}  map[string]string
// @Router       /api/files/{id} [get]
func (h *FileHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	result, err := h.svc.GetByID(id)
	if err != nil {
		handleFileServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Download godoc
// @Summary      Download a file's content
// @Tags         Files
// @Produce      application/octet-stream
// @Param        id   path  int  true  "File ID"
// @Success      200  {file}    binary
// @Failure      404  {object}  map[string]string
// @Router       /api/files/{id}/download [get]
func (h *FileHandler) Download(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	target, err := h.svc.GetDownloadTarget(id)
	if err != nil {
		handleFileServiceError(c, err)
		return
	}
	c.Header("Content-Type", target.ContentType)
	c.FileAttachment(target.Path, target.Name)
}

// Upload godoc
// @Summary      Upload a file
// @Tags         Files
// @Accept       multipart/form-data
// @Produce      json
// @Param        file  formData  file  true  "File to upload"
// @Success      201   {object}  dto.FileResponse
// @Failure      400   {object}  map[string]string
// @Failure      413   {object}  map[string]string
// @Router       /api/files [post]
func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file is required"})
		return
	}

	src, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unable to read uploaded file"})
		return
	}
	defer src.Close()

	result, err := h.svc.Upload(service.UploadInput{
		OriginalName: fileHeader.Filename,
		ContentType:  fileHeader.Header.Get("Content-Type"),
		Size:         fileHeader.Size,
		Reader:       src,
	})
	if err != nil {
		handleFileServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// Delete godoc
// @Summary      Delete a file
// @Tags         Files
// @Param        id  path  int  true  "File ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /api/files/{id} [delete]
func (h *FileHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleFileServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
