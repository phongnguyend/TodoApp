package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"github.com/todo/backend/go/internal/service"
	"gorm.io/gorm"
)

type attachmentRepoMock struct {
	findAll  func(uint) ([]models.TodoItemAttachment, error)
	findID   func(uint, uint) (*models.TodoItemAttachment, error)
	findFile func(uint, uint) (*models.TodoItemAttachment, error)
	create   func(*models.TodoItemAttachment) (*models.TodoItemAttachment, error)
	update   func(*models.TodoItemAttachment) (*models.TodoItemAttachment, error)
	delete   func(*models.TodoItemAttachment) error
}

func (m *attachmentRepoMock) FindAllByTodoItemID(id uint) ([]models.TodoItemAttachment, error) {
	return m.findAll(id)
}
func (m *attachmentRepoMock) FindByIDAndTodoItemID(id, todo uint) (*models.TodoItemAttachment, error) {
	return m.findID(id, todo)
}
func (m *attachmentRepoMock) FindByTodoItemIDAndFileID(todo, file uint) (*models.TodoItemAttachment, error) {
	return m.findFile(todo, file)
}
func (m *attachmentRepoMock) Create(a *models.TodoItemAttachment) (*models.TodoItemAttachment, error) {
	return m.create(a)
}
func (m *attachmentRepoMock) Update(a *models.TodoItemAttachment) (*models.TodoItemAttachment, error) {
	return m.update(a)
}
func (m *attachmentRepoMock) Delete(a *models.TodoItemAttachment) error { return m.delete(a) }

type attachmentTodoRepoMock struct {
	repository.TodoItemRepository
	find func(uint) (*models.TodoItem, error)
}

func (m *attachmentTodoRepoMock) FindByID(id uint) (*models.TodoItem, error) { return m.find(id) }

type attachmentFileRepoMock struct {
	repository.FileRepository
	find func(uint) (*models.File, error)
}

func (m *attachmentFileRepoMock) FindByID(id uint) (*models.File, error) { return m.find(id) }

func attachmentService(repo *attachmentRepoMock, todoErr, fileErr error) service.TodoItemAttachmentService {
	todos := &attachmentTodoRepoMock{find: func(id uint) (*models.TodoItem, error) { return &models.TodoItem{ID: id}, todoErr }}
	files := &attachmentFileRepoMock{find: func(id uint) (*models.File, error) { return &models.File{ID: id}, fileErr }}
	return service.NewTodoItemAttachmentService(repo, todos, files)
}

func TestAttachmentGetAllScopesToExistingTodo(t *testing.T) {
	repo := &attachmentRepoMock{findAll: func(id uint) ([]models.TodoItemAttachment, error) {
		assert.Equal(t, uint(10), id)
		return []models.TodoItemAttachment{{ID: 1, TodoItemID: id, FileID: 5}}, nil
	}}
	result, err := attachmentService(repo, nil, nil).GetAll(10)
	require.NoError(t, err)
	require.Len(t, result, 1)
	assert.Equal(t, uint(5), result[0].FileID)
}

func TestAttachmentGetAllMissingTodo(t *testing.T) {
	repo := &attachmentRepoMock{}
	_, err := attachmentService(repo, gorm.ErrRecordNotFound, nil).GetAll(99)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestAttachmentCreateRejectsMissingFile(t *testing.T) {
	repo := &attachmentRepoMock{}
	_, err := attachmentService(repo, nil, gorm.ErrRecordNotFound).Create(10, 99)
	require.ErrorIs(t, err, service.ErrFileNotFound)
}

func TestAttachmentCreateIsIdempotent(t *testing.T) {
	existing := &models.TodoItemAttachment{ID: 3, TodoItemID: 10, FileID: 5}
	repo := &attachmentRepoMock{findFile: func(todo, file uint) (*models.TodoItemAttachment, error) { return existing, nil }}
	result, err := attachmentService(repo, nil, nil).Create(10, 5)
	require.NoError(t, err)
	assert.Equal(t, uint(3), result.ID)
}

func TestAttachmentCreatePersistsNewLink(t *testing.T) {
	repo := &attachmentRepoMock{
		findFile: func(todo, file uint) (*models.TodoItemAttachment, error) { return nil, gorm.ErrRecordNotFound },
		create: func(a *models.TodoItemAttachment) (*models.TodoItemAttachment, error) {
			assert.Equal(t, uint(10), a.TodoItemID)
			a.ID = 7
			return a, nil
		},
	}
	result, err := attachmentService(repo, nil, nil).Create(10, 5)
	require.NoError(t, err)
	assert.Equal(t, uint(7), result.ID)
}

func TestAttachmentGetByIDRejectsOtherTodoAttachment(t *testing.T) {
	repo := &attachmentRepoMock{findID: func(id, todo uint) (*models.TodoItemAttachment, error) { return nil, gorm.ErrRecordNotFound }}
	_, err := attachmentService(repo, nil, nil).GetByID(10, 3)
	require.ErrorIs(t, err, service.ErrAttachmentNotFound)
}

func TestAttachmentUpdateChangesFile(t *testing.T) {
	current := &models.TodoItemAttachment{ID: 3, TodoItemID: 10, FileID: 5}
	repo := &attachmentRepoMock{
		findID:   func(id, todo uint) (*models.TodoItemAttachment, error) { return current, nil },
		findFile: func(todo, file uint) (*models.TodoItemAttachment, error) { return nil, gorm.ErrRecordNotFound },
		update:   func(a *models.TodoItemAttachment) (*models.TodoItemAttachment, error) { return a, nil },
	}
	result, err := attachmentService(repo, nil, nil).Update(10, 3, 6)
	require.NoError(t, err)
	assert.Equal(t, uint(6), result.FileID)
}

func TestAttachmentDeleteRemovesScopedLink(t *testing.T) {
	deleted := false
	repo := &attachmentRepoMock{
		findID: func(id, todo uint) (*models.TodoItemAttachment, error) {
			return &models.TodoItemAttachment{ID: id, TodoItemID: todo}, nil
		},
		delete: func(a *models.TodoItemAttachment) error { deleted = true; return nil },
	}
	require.NoError(t, attachmentService(repo, nil, nil).Delete(10, 3))
	assert.True(t, deleted)
}
