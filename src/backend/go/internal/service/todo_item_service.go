package service

import (
	"errors"
	"math"

	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"gorm.io/gorm"
)

// ErrNotFound is returned when a requested entity does not exist.
// Mirrors a NotFoundException / 404 pattern in ASP.NET Core services.
var ErrNotFound = errors.New("todo item not found")

// TodoItemService defines the business-logic contract.
// Mirrors ITodoItemService in C#.
type TodoItemService interface {
	GetAll(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error)
	GetIncomplete(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error)
	GetByID(id uint) (dto.TodoItemResponse, error)
	Create(req dto.CreateTodoItemRequest) (dto.TodoItemResponse, error)
	Update(id uint, req dto.UpdateTodoItemRequest) (dto.TodoItemResponse, error)
	MarkComplete(id uint) (dto.TodoItemResponse, error)
	Delete(id uint) error
}

type todoItemService struct {
	repo repository.TodoItemRepository
}

// NewTodoItemService constructs the service with its repository dependency injected.
func NewTodoItemService(repo repository.TodoItemRepository) TodoItemService {
	return &todoItemService{repo: repo}
}

// ── Mapping ───────────────────────────────────────────────────────────────────

func toResponse(m *models.TodoItem) dto.TodoItemResponse {
	return dto.TodoItemResponse{
		ID:          m.ID,
		Title:       m.Title,
		Description: m.Description,
		IsCompleted: m.IsCompleted,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toPaginated(result repository.PaginatedResult, page, pageSize int) dto.PaginatedResponse[dto.TodoItemResponse] {
	items := make([]dto.TodoItemResponse, len(result.Items))
	for i := range result.Items {
		items[i] = toResponse(&result.Items[i])
	}
	return dto.PaginatedResponse[dto.TodoItemResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(result.Total) / float64(pageSize))),
	}
}

func (s *todoItemService) getOrNotFound(id uint) (*models.TodoItem, error) {
	item, err := s.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return item, err
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (s *todoItemService) GetAll(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
	result, err := s.repo.FindAll((page-1)*pageSize, pageSize)
	if err != nil {
		return dto.PaginatedResponse[dto.TodoItemResponse]{}, err
	}
	return toPaginated(result, page, pageSize), nil
}

func (s *todoItemService) GetIncomplete(page, pageSize int) (dto.PaginatedResponse[dto.TodoItemResponse], error) {
	result, err := s.repo.FindIncomplete((page-1)*pageSize, pageSize)
	if err != nil {
		return dto.PaginatedResponse[dto.TodoItemResponse]{}, err
	}
	return toPaginated(result, page, pageSize), nil
}

func (s *todoItemService) GetByID(id uint) (dto.TodoItemResponse, error) {
	item, err := s.getOrNotFound(id)
	if err != nil {
		return dto.TodoItemResponse{}, err
	}
	return toResponse(item), nil
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (s *todoItemService) Create(req dto.CreateTodoItemRequest) (dto.TodoItemResponse, error) {
	item := &models.TodoItem{Title: req.Title, Description: req.Description}
	created, err := s.repo.Create(item)
	if err != nil {
		return dto.TodoItemResponse{}, err
	}
	return toResponse(created), nil
}

func (s *todoItemService) Update(id uint, req dto.UpdateTodoItemRequest) (dto.TodoItemResponse, error) {
	item, err := s.getOrNotFound(id)
	if err != nil {
		return dto.TodoItemResponse{}, err
	}
	if req.Title != nil {
		item.Title = *req.Title
	}
	if req.Description != nil {
		item.Description = req.Description
	}
	if req.IsCompleted != nil {
		item.IsCompleted = *req.IsCompleted
	}
	updated, err := s.repo.Update(item)
	if err != nil {
		return dto.TodoItemResponse{}, err
	}
	return toResponse(updated), nil
}

func (s *todoItemService) MarkComplete(id uint) (dto.TodoItemResponse, error) {
	item, err := s.getOrNotFound(id)
	if err != nil {
		return dto.TodoItemResponse{}, err
	}
	item.IsCompleted = true
	updated, err := s.repo.Update(item)
	if err != nil {
		return dto.TodoItemResponse{}, err
	}
	return toResponse(updated), nil
}

func (s *todoItemService) Delete(id uint) error {
	item, err := s.getOrNotFound(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(item)
}
