package service

import (
	"errors"

	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"gorm.io/gorm"
)

var ErrAttachmentNotFound = errors.New("attachment not found")

type TodoItemAttachmentService interface {
	GetAll(todoItemID uint) ([]dto.TodoItemAttachmentResponse, error)
	GetByID(todoItemID, attachmentID uint) (dto.TodoItemAttachmentResponse, error)
	Create(todoItemID, fileID uint, actorUserID ...*uint) (dto.TodoItemAttachmentResponse, error)
	Update(todoItemID, attachmentID, fileID uint, actorUserID ...*uint) (dto.TodoItemAttachmentResponse, error)
	Delete(todoItemID, attachmentID uint) error
}

type todoItemAttachmentService struct {
	repo  repository.TodoItemAttachmentRepository
	todos repository.TodoItemRepository
	files repository.FileRepository
}

func NewTodoItemAttachmentService(repo repository.TodoItemAttachmentRepository, todos repository.TodoItemRepository, files repository.FileRepository) TodoItemAttachmentService {
	return &todoItemAttachmentService{repo: repo, todos: todos, files: files}
}

func attachmentResponse(a *models.TodoItemAttachment) dto.TodoItemAttachmentResponse {
	return dto.TodoItemAttachmentResponse{ID: a.ID, TodoItemID: a.TodoItemID, FileID: a.FileID,
		CreatedAt: a.CreatedAt, CreatedByUserID: a.CreatedByUserID,
		UpdatedAt: a.UpdatedAt, UpdatedByUserID: a.UpdatedByUserID}
}
func (s *todoItemAttachmentService) requireTodo(id uint) error {
	_, err := s.todos.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
func (s *todoItemAttachmentService) requireFile(id uint) error {
	_, err := s.files.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrFileNotFound
	}
	return err
}
func (s *todoItemAttachmentService) requireAttachment(todoID, id uint) (*models.TodoItemAttachment, error) {
	a, err := s.repo.FindByIDAndTodoItemID(id, todoID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrAttachmentNotFound
	}
	return a, err
}
func (s *todoItemAttachmentService) GetAll(todoID uint) ([]dto.TodoItemAttachmentResponse, error) {
	if err := s.requireTodo(todoID); err != nil {
		return nil, err
	}
	items, err := s.repo.FindAllByTodoItemID(todoID)
	if err != nil {
		return nil, err
	}
	result := make([]dto.TodoItemAttachmentResponse, len(items))
	for i := range items {
		result[i] = attachmentResponse(&items[i])
	}
	return result, nil
}
func (s *todoItemAttachmentService) GetByID(todoID, id uint) (dto.TodoItemAttachmentResponse, error) {
	if err := s.requireTodo(todoID); err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	a, err := s.requireAttachment(todoID, id)
	if err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	return attachmentResponse(a), nil
}
func (s *todoItemAttachmentService) Create(todoID, fileID uint, actorUserID ...*uint) (dto.TodoItemAttachmentResponse, error) {
	if err := s.requireTodo(todoID); err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	if err := s.requireFile(fileID); err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	a, err := s.repo.FindByTodoItemIDAndFileID(todoID, fileID)
	if err == nil {
		return attachmentResponse(a), nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return dto.TodoItemAttachmentResponse{}, err
	}
	a, err = s.repo.Create(&models.TodoItemAttachment{TodoItemID: todoID, FileID: fileID,
		CreatedByUserID: firstActor(actorUserID)})
	if err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	return attachmentResponse(a), nil
}
func (s *todoItemAttachmentService) Update(todoID, id, fileID uint, actorUserID ...*uint) (dto.TodoItemAttachmentResponse, error) {
	if err := s.requireTodo(todoID); err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	if err := s.requireFile(fileID); err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	a, err := s.requireAttachment(todoID, id)
	if err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	existing, findErr := s.repo.FindByTodoItemIDAndFileID(todoID, fileID)
	if findErr == nil && existing.ID != id {
		return attachmentResponse(existing), nil
	}
	if findErr != nil && !errors.Is(findErr, gorm.ErrRecordNotFound) {
		return dto.TodoItemAttachmentResponse{}, findErr
	}
	a.FileID = fileID
	a.UpdatedByUserID = firstActor(actorUserID)
	a, err = s.repo.Update(a)
	if err != nil {
		return dto.TodoItemAttachmentResponse{}, err
	}
	return attachmentResponse(a), nil
}
func (s *todoItemAttachmentService) Delete(todoID, id uint) error {
	if err := s.requireTodo(todoID); err != nil {
		return err
	}
	a, err := s.requireAttachment(todoID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(a)
}
