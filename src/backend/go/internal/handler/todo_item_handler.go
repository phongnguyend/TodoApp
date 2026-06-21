package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/service"
)

// TodoItemHandler handles HTTP requests — analogous to a [ApiController] in ASP.NET Core.
// Each method maps to a route registered in the router.
type TodoItemHandler struct {
	svc service.TodoItemService
}

// NewTodoItemHandler creates the handler with its service dependency injected.
func NewTodoItemHandler(svc service.TodoItemService) *TodoItemHandler {
	return &TodoItemHandler{svc: svc}
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func queryInt(c *gin.Context, key string, fallback int) int {
	if v := c.Query(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil && i > 0 {
			return i
		}
	}
	return fallback
}

func handleServiceError(c *gin.Context, err error) {
	if errors.Is(err, service.ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
}

// ── Endpoints ─────────────────────────────────────────────────────────────────

// GetAll godoc
// @Summary      List all todo items
// @Tags         Todo Items
// @Produce      json
// @Param        page     query  int  false  "Page number"   default(1)
// @Param        pageSize query  int  false  "Page size"     default(20)
// @Success      200  {object}  dto.PaginatedResponse[dto.TodoItemResponse]
// @Router       /api/todo-items [get]
func (h *TodoItemHandler) GetAll(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)

	result, err := h.svc.GetAll(page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetIncomplete godoc
// @Summary      List incomplete todo items
// @Tags         Todo Items
// @Produce      json
// @Param        page     query  int  false  "Page number"   default(1)
// @Param        pageSize query  int  false  "Page size"     default(20)
// @Success      200  {object}  dto.PaginatedResponse[dto.TodoItemResponse]
// @Router       /api/todo-items/incomplete [get]
func (h *TodoItemHandler) GetIncomplete(c *gin.Context) {
	page := queryInt(c, "page", 1)
	pageSize := queryInt(c, "pageSize", 20)

	result, err := h.svc.GetIncomplete(page, pageSize)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// GetByID godoc
// @Summary      Get a todo item by ID
// @Tags         Todo Items
// @Produce      json
// @Param        id   path  int  true  "Todo item ID"
// @Success      200  {object}  dto.TodoItemResponse
// @Failure      404  {object}  map[string]string
// @Router       /api/todo-items/{id} [get]
func (h *TodoItemHandler) GetByID(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	result, err := h.svc.GetByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Create godoc
// @Summary      Create a todo item
// @Tags         Todo Items
// @Accept       json
// @Produce      json
// @Param        body  body  dto.CreateTodoItemRequest  true  "Request body"
// @Success      201   {object}  dto.TodoItemResponse
// @Failure      400   {object}  map[string]string
// @Router       /api/todo-items [post]
func (h *TodoItemHandler) Create(c *gin.Context) {
	var req dto.CreateTodoItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.Create(req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, result)
}

// Update godoc
// @Summary      Update a todo item
// @Tags         Todo Items
// @Accept       json
// @Produce      json
// @Param        id    path  int                        true  "Todo item ID"
// @Param        body  body  dto.UpdateTodoItemRequest  true  "Request body"
// @Success      200   {object}  dto.TodoItemResponse
// @Failure      400   {object}  map[string]string
// @Failure      404   {object}  map[string]string
// @Router       /api/todo-items/{id} [put]
func (h *TodoItemHandler) Update(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	var req dto.UpdateTodoItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.Update(id, req)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// MarkComplete godoc
// @Summary      Mark a todo item as complete
// @Tags         Todo Items
// @Produce      json
// @Param        id  path  int  true  "Todo item ID"
// @Success      200  {object}  dto.TodoItemResponse
// @Failure      404  {object}  map[string]string
// @Router       /api/todo-items/{id}/complete [patch]
func (h *TodoItemHandler) MarkComplete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	result, err := h.svc.MarkComplete(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	c.JSON(http.StatusOK, result)
}

// Delete godoc
// @Summary      Delete a todo item
// @Tags         Todo Items
// @Param        id  path  int  true  "Todo item ID"
// @Success      204
// @Failure      404  {object}  map[string]string
// @Router       /api/todo-items/{id} [delete]
func (h *TodoItemHandler) Delete(c *gin.Context) {
	id, err := parseID(c)
	if err != nil {
		return
	}
	if err := h.svc.Delete(id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseID(c *gin.Context) (uint, error) {
	raw := c.Param("id")
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, err
	}
	return uint(id), nil
}
