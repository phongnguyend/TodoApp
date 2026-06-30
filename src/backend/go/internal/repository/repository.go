package repository

import "github.com/todo/backend/go/internal/models"

// PaginatedResult is the repository-level paginated result — no DTO dependency here.
type PaginatedResult struct {
	Items []models.TodoItem
	Total int64
}

// TodoItemRepository defines the data-access contract.
// Mirrors IRepository<TodoItem> / ITodoItemRepository in C#.
type TodoItemRepository interface {
	FindAll(skip, limit int) (PaginatedResult, error)
	FindIncomplete(skip, limit int) (PaginatedResult, error)
	FindByID(id uint) (*models.TodoItem, error)
	Create(item *models.TodoItem) (*models.TodoItem, error)
	Update(item *models.TodoItem) (*models.TodoItem, error)
	Delete(item *models.TodoItem) error
}

// EmailLogRepository defines the data-access contract for email audit logs.
type EmailLogRepository interface {
	Create(log *models.EmailLog) (*models.EmailLog, error)
	MarkSent(log *models.EmailLog) error
	MarkFailed(log *models.EmailLog, errMsg string) error
}
