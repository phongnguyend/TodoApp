package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/service"
)

type TodoItemAttachmentHandler struct {
	svc service.TodoItemAttachmentService
}

func NewTodoItemAttachmentHandler(svc service.TodoItemAttachmentService) *TodoItemAttachmentHandler {
	return &TodoItemAttachmentHandler{svc: svc}
}

func attachmentID(c *gin.Context, name string) (uint, bool) {
	v, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || v == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": name + " must be a positive integer"})
		return 0, false
	}
	return uint(v), true
}
func handleAttachmentError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrNotFound) || errors.Is(err, service.ErrFileNotFound) || errors.Is(err, service.ErrAttachmentNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}
func bindAttachment(c *gin.Context) (dto.SaveTodoItemAttachmentRequest, bool) {
	var request dto.SaveTodoItemAttachmentRequest
	if err := c.ShouldBindJSON(&request); err != nil || request.FileID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "fileId is required and must be positive"})
		return request, false
	}
	return request, true
}
func (h *TodoItemAttachmentHandler) GetAll(c *gin.Context) {
	todoID, ok := attachmentID(c, "id")
	if !ok {
		return
	}
	result, err := h.svc.GetAll(todoID)
	if err != nil {
		handleAttachmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
func (h *TodoItemAttachmentHandler) GetByID(c *gin.Context) {
	todoID, ok := attachmentID(c, "id")
	if !ok {
		return
	}
	id, ok := attachmentID(c, "attachmentId")
	if !ok {
		return
	}
	result, err := h.svc.GetByID(todoID, id)
	if err != nil {
		handleAttachmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
func (h *TodoItemAttachmentHandler) Create(c *gin.Context) {
	todoID, ok := attachmentID(c, "id")
	if !ok {
		return
	}
	request, ok := bindAttachment(c)
	if !ok {
		return
	}
	result, err := h.svc.Create(todoID, request.FileID)
	if err != nil {
		handleAttachmentError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}
func (h *TodoItemAttachmentHandler) Update(c *gin.Context) {
	todoID, ok := attachmentID(c, "id")
	if !ok {
		return
	}
	id, ok := attachmentID(c, "attachmentId")
	if !ok {
		return
	}
	request, ok := bindAttachment(c)
	if !ok {
		return
	}
	result, err := h.svc.Update(todoID, id, request.FileID)
	if err != nil {
		handleAttachmentError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}
func (h *TodoItemAttachmentHandler) Delete(c *gin.Context) {
	todoID, ok := attachmentID(c, "id")
	if !ok {
		return
	}
	id, ok := attachmentID(c, "attachmentId")
	if !ok {
		return
	}
	if err := h.svc.Delete(todoID, id); err != nil {
		handleAttachmentError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
