package repository

import (
	"github.com/todo/backend/go/internal/models"
	"gorm.io/gorm"
)

// todoItemRepository is the GORM-backed implementation of TodoItemRepository.
// Mirrors a repository class backed by EF Core's DbSet<TodoItem>.
type todoItemRepository struct {
	db *gorm.DB
}

// NewTodoItemRepository creates a new repository — called from the DI composition root.
func NewTodoItemRepository(db *gorm.DB) TodoItemRepository {
	return &todoItemRepository{db: db}
}

func (r *todoItemRepository) FindAll(skip, limit int) (PaginatedResult, error) {
	var items []models.TodoItem
	var total int64

	if err := r.db.Model(&models.TodoItem{}).Count(&total).Error; err != nil {
		return PaginatedResult{}, err
	}
	if err := r.db.Offset(skip).Limit(limit).Order("created_at desc").Find(&items).Error; err != nil {
		return PaginatedResult{}, err
	}
	return PaginatedResult{Items: items, Total: total}, nil
}

func (r *todoItemRepository) FindIncomplete(skip, limit int) (PaginatedResult, error) {
	var items []models.TodoItem
	var total int64

	q := r.db.Model(&models.TodoItem{}).Where("is_completed = ?", false)
	if err := q.Count(&total).Error; err != nil {
		return PaginatedResult{}, err
	}
	if err := q.Offset(skip).Limit(limit).Order("created_at desc").Find(&items).Error; err != nil {
		return PaginatedResult{}, err
	}
	return PaginatedResult{Items: items, Total: total}, nil
}

func (r *todoItemRepository) FindByID(id uint) (*models.TodoItem, error) {
	var item models.TodoItem
	if err := r.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *todoItemRepository) Create(item *models.TodoItem) (*models.TodoItem, error) {
	if err := r.db.Create(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *todoItemRepository) Update(item *models.TodoItem) (*models.TodoItem, error) {
	if err := r.db.Save(item).Error; err != nil {
		return nil, err
	}
	return item, nil
}

func (r *todoItemRepository) Delete(item *models.TodoItem) error {
	return r.db.Delete(item).Error
}
