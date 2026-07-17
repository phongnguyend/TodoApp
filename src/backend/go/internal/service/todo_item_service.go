package service

import (
	"encoding/csv"
	"errors"
	"io"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"gorm.io/gorm"
)

// ErrNotFound is returned when a requested entity does not exist.
// Mirrors a NotFoundException / 404 pattern in ASP.NET Core services.
var ErrNotFound = errors.New("todo item not found")

// csvFieldNames is the CSV header row - matches the Python/PHP/Java implementations exactly.
var csvFieldNames = []string{"id", "title", "description", "is_completed", "created_at", "updated_at"}

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
	ImportCSV(r io.Reader) (dto.ImportResult, error)
	ExportCSV() (string, error)
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

// ── CSV import/export ─────────────────────────────────────────────────────────

func parseBool(value string) bool {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "yes", "y":
		return true
	default:
		return false
	}
}

// ImportCSV reads a CSV file (header: title, description, is_completed) and creates a todo
// item for every valid row. Mirrors import_csv in the Python/PHP/Java implementations.
func (s *todoItemService) ImportCSV(r io.Reader) (dto.ImportResult, error) {
	reader := csv.NewReader(r)
	reader.FieldsPerRecord = -1

	header, err := reader.Read()
	if err != nil {
		if err == io.EOF {
			return dto.ImportResult{Errors: []dto.ImportRowError{}}, nil
		}
		return dto.ImportResult{}, err
	}
	colIndex := make(map[string]int, len(header))
	for i, name := range header {
		colIndex[strings.ToLower(strings.TrimSpace(name))] = i
	}
	cell := func(row []string, name string) string {
		idx, ok := colIndex[name]
		if !ok || idx >= len(row) {
			return ""
		}
		return row[idx]
	}

	imported := 0
	errs := []dto.ImportRowError{}
	rowNumber := 1 // header is row 1
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return dto.ImportResult{}, err
		}
		rowNumber++

		title := strings.TrimSpace(cell(row, "title"))
		if title == "" {
			errs = append(errs, dto.ImportRowError{Row: rowNumber, Error: "Title is required."})
			continue
		}

		item := &models.TodoItem{
			Title:       title,
			IsCompleted: parseBool(cell(row, "is_completed")),
		}
		if description := strings.TrimSpace(cell(row, "description")); description != "" {
			item.Description = &description
		}

		if _, err := s.repo.Create(item); err != nil {
			return dto.ImportResult{}, err
		}
		imported++
	}

	return dto.ImportResult{Imported: imported, Failed: len(errs), Errors: errs}, nil
}

// ExportCSV renders every todo item as CSV text. Mirrors export_csv in the
// Python/PHP/Java implementations (same header row and column order).
func (s *todoItemService) ExportCSV() (string, error) {
	items, err := s.repo.FindAllItems()
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	w := csv.NewWriter(&buf)
	if err := w.Write(csvFieldNames); err != nil {
		return "", err
	}
	for _, item := range items {
		description := ""
		if item.Description != nil {
			description = *item.Description
		}
		updatedAt := ""
		if item.UpdatedAt != nil {
			updatedAt = item.UpdatedAt.Format(time.RFC3339Nano)
		}
		record := []string{
			strconv.FormatUint(uint64(item.ID), 10),
			item.Title,
			description,
			strconv.FormatBool(item.IsCompleted),
			item.CreatedAt.Format(time.RFC3339Nano),
			updatedAt,
		}
		if err := w.Write(record); err != nil {
			return "", err
		}
	}
	w.Flush()
	if err := w.Error(); err != nil {
		return "", err
	}
	return buf.String(), nil
}
