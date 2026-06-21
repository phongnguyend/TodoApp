package dto

import "time"

// CreateTodoItemRequest mirrors a CreateTodoItemRequest command DTO in C#.
// The `binding` tag drives Gin's validation (analogous to [Required] / FluentValidation).
type CreateTodoItemRequest struct {
	Title       string  `json:"title"       binding:"required,min=1,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
}

// UpdateTodoItemRequest mirrors an UpdateTodoItemRequest DTO — all fields optional.
type UpdateTodoItemRequest struct {
	Title       *string `json:"title"       binding:"omitempty,min=1,max=200"`
	Description *string `json:"description" binding:"omitempty,max=2000"`
	IsCompleted *bool   `json:"isCompleted"`
}

// TodoItemResponse mirrors a TodoItemDto / view model returned from controllers.
type TodoItemResponse struct {
	ID          uint       `json:"id"`
	Title       string     `json:"title"`
	Description *string    `json:"description"`
	IsCompleted bool       `json:"isCompleted"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   *time.Time `json:"updatedAt"`
}

// PaginatedResponse is a generic paginated list wrapper — mirrors PagedResult<T> in C#.
type PaginatedResponse[T any] struct {
	Items      []T   `json:"items"`
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	PageSize   int   `json:"pageSize"`
	TotalPages int   `json:"totalPages"`
}
