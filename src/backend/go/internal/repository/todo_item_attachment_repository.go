package repository

import (
	"github.com/todo/backend/go/internal/models"
	"gorm.io/gorm"
)

type todoItemAttachmentRepository struct{ db *gorm.DB }

func NewTodoItemAttachmentRepository(db *gorm.DB) TodoItemAttachmentRepository {
	return &todoItemAttachmentRepository{db: db}
}
func (r *todoItemAttachmentRepository) FindAllByTodoItemID(todoItemID uint) ([]models.TodoItemAttachment, error) {
	var items []models.TodoItemAttachment
	err := r.db.Where("todo_item_id = ?", todoItemID).Order("created_at asc").Find(&items).Error
	return items, err
}
func (r *todoItemAttachmentRepository) FindByIDAndTodoItemID(id, todoItemID uint) (*models.TodoItemAttachment, error) {
	var item models.TodoItemAttachment
	err := r.db.Where("id = ? AND todo_item_id = ?", id, todoItemID).First(&item).Error
	return &item, err
}
func (r *todoItemAttachmentRepository) FindByTodoItemIDAndFileID(todoItemID, fileID uint) (*models.TodoItemAttachment, error) {
	var item models.TodoItemAttachment
	err := r.db.Where("todo_item_id = ? AND file_id = ?", todoItemID, fileID).First(&item).Error
	return &item, err
}
func (r *todoItemAttachmentRepository) Create(item *models.TodoItemAttachment) (*models.TodoItemAttachment, error) {
	return item, r.db.Create(item).Error
}
func (r *todoItemAttachmentRepository) Update(item *models.TodoItemAttachment) (*models.TodoItemAttachment, error) {
	return item, r.db.Save(item).Error
}
func (r *todoItemAttachmentRepository) Delete(item *models.TodoItemAttachment) error {
	return r.db.Delete(item).Error
}
