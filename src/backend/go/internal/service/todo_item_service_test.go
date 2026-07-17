package service_test

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"github.com/todo/backend/go/internal/service"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// ── mock repository ───────────────────────────────────────────────────────────

type mockRepo struct {
	findAllFn        func(skip, limit int) (repository.PaginatedResult, error)
	findIncompleteFn func(skip, limit int) (repository.PaginatedResult, error)
	findAllItemsFn   func() ([]models.TodoItem, error)
	findByIDFn       func(id uint) (*models.TodoItem, error)
	createFn         func(item *models.TodoItem) (*models.TodoItem, error)
	updateFn         func(item *models.TodoItem) (*models.TodoItem, error)
	deleteFn         func(item *models.TodoItem) error
}

func (m *mockRepo) FindAll(skip, limit int) (repository.PaginatedResult, error) {
	return m.findAllFn(skip, limit)
}
func (m *mockRepo) FindIncomplete(skip, limit int) (repository.PaginatedResult, error) {
	return m.findIncompleteFn(skip, limit)
}
func (m *mockRepo) FindAllItems() ([]models.TodoItem, error) {
	return m.findAllItemsFn()
}
func (m *mockRepo) FindByID(id uint) (*models.TodoItem, error) {
	return m.findByIDFn(id)
}
func (m *mockRepo) Create(item *models.TodoItem) (*models.TodoItem, error) {
	return m.createFn(item)
}
func (m *mockRepo) Update(item *models.TodoItem) (*models.TodoItem, error) {
	return m.updateFn(item)
}
func (m *mockRepo) Delete(item *models.TodoItem) error {
	return m.deleteFn(item)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func sampleItem() *models.TodoItem {
	desc := "a description"
	return &models.TodoItem{
		ID:          1,
		Title:       "Buy milk",
		Description: &desc,
		IsCompleted: false,
		CreatedAt:   time.Now(),
	}
}

// ── GetAll ────────────────────────────────────────────────────────────────────

func TestGetAll_ReturnsPaginatedItems(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findAllFn: func(skip, limit int) (repository.PaginatedResult, error) {
			assert.Equal(t, 0, skip) // page 1 → skip = 0
			assert.Equal(t, 10, limit)
			return repository.PaginatedResult{Items: []models.TodoItem{*item}, Total: 25}, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	result, err := svc.GetAll(1, 10)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, int64(25), result.Total)
	assert.Equal(t, 3, result.TotalPages) // ceil(25/10) = 3
	assert.Len(t, result.Items, 1)
	assert.Equal(t, item.Title, result.Items[0].Title)
}

func TestGetAll_Page2_CorrectSkip(t *testing.T) {
	repo := &mockRepo{
		findAllFn: func(skip, limit int) (repository.PaginatedResult, error) {
			assert.Equal(t, 10, skip) // page 2 → skip = 10
			return repository.PaginatedResult{Items: []models.TodoItem{}, Total: 25}, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.GetAll(2, 10)
	require.NoError(t, err)
}

func TestGetAll_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		findAllFn: func(skip, limit int) (repository.PaginatedResult, error) {
			return repository.PaginatedResult{}, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.GetAll(1, 10)
	require.Error(t, err)
}

// ── GetIncomplete ─────────────────────────────────────────────────────────────

func TestGetIncomplete_ReturnsPaginatedItems(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findIncompleteFn: func(skip, limit int) (repository.PaginatedResult, error) {
			return repository.PaginatedResult{Items: []models.TodoItem{*item}, Total: 1}, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	result, err := svc.GetIncomplete(1, 20)

	require.NoError(t, err)
	assert.Len(t, result.Items, 1)
	assert.False(t, result.Items[0].IsCompleted)
}

func TestGetIncomplete_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		findIncompleteFn: func(skip, limit int) (repository.PaginatedResult, error) {
			return repository.PaginatedResult{}, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.GetIncomplete(1, 20)
	require.Error(t, err)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestGetByID_ReturnsItem(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) {
			assert.Equal(t, uint(1), id)
			return item, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	result, err := svc.GetByID(1)

	require.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, item.Title, result.Title)
	assert.Equal(t, item.Description, result.Description)
}

func TestGetByID_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.GetByID(99)

	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestGetByID_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.GetByID(1)

	require.Error(t, err)
	assert.NotErrorIs(t, err, service.ErrNotFound)
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreate_ReturnsCreatedItem(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		createFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			assert.Equal(t, item.Title, in.Title)
			assert.Equal(t, item.Description, in.Description)
			return item, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	result, err := svc.Create(dto.CreateTodoItemRequest{Title: item.Title, Description: item.Description})

	require.NoError(t, err)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, item.Title, result.Title)
}

func TestCreate_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		createFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			return nil, errors.New("constraint error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.Create(dto.CreateTodoItemRequest{Title: "Test"})
	require.Error(t, err)
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdate_UpdatesFields(t *testing.T) {
	item := sampleItem()
	newTitle := "Updated Title"
	newDesc := "Updated Desc"
	updated := *item
	updated.Title = newTitle
	updated.Description = &newDesc

	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		updateFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			assert.Equal(t, newTitle, in.Title)
			assert.Equal(t, newDesc, *in.Description)
			return &updated, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	result, err := svc.Update(1, dto.UpdateTodoItemRequest{Title: &newTitle, Description: &newDesc})

	require.NoError(t, err)
	assert.Equal(t, newTitle, result.Title)
	assert.Equal(t, newDesc, *result.Description)
}

func TestUpdate_PartialUpdate_OnlyTitle(t *testing.T) {
	item := sampleItem()
	newTitle := "New Title"
	originalDesc := *item.Description

	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		updateFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			assert.Equal(t, newTitle, in.Title)
			assert.Equal(t, originalDesc, *in.Description) // unchanged
			return in, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.Update(1, dto.UpdateTodoItemRequest{Title: &newTitle})
	require.NoError(t, err)
}

func TestUpdate_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.Update(99, dto.UpdateTodoItemRequest{})
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestUpdate_RepoUpdateError_Propagates(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		updateFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.Update(1, dto.UpdateTodoItemRequest{})
	require.Error(t, err)
}

// ── MarkComplete ──────────────────────────────────────────────────────────────

func TestMarkComplete_SetsIsCompleted(t *testing.T) {
	item := sampleItem()
	completed := *item
	completed.IsCompleted = true

	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		updateFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			assert.True(t, in.IsCompleted)
			return &completed, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	result, err := svc.MarkComplete(1)

	require.NoError(t, err)
	assert.True(t, result.IsCompleted)
}

func TestMarkComplete_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.MarkComplete(99)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestMarkComplete_RepoUpdateError_Propagates(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		updateFn: func(in *models.TodoItem) (*models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.MarkComplete(1)
	require.Error(t, err)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	item := sampleItem()
	deleted := false

	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		deleteFn: func(in *models.TodoItem) error {
			assert.Equal(t, item.ID, in.ID)
			deleted = true
			return nil
		},
	}
	svc := service.NewTodoItemService(repo)

	err := svc.Delete(1)

	require.NoError(t, err)
	assert.True(t, deleted)
}

func TestDelete_NotFound_ReturnsErrNotFound(t *testing.T) {
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := service.NewTodoItemService(repo)

	err := svc.Delete(99)
	require.ErrorIs(t, err, service.ErrNotFound)
}

func TestDelete_RepoDeleteError_Propagates(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findByIDFn: func(id uint) (*models.TodoItem, error) { return item, nil },
		deleteFn: func(in *models.TodoItem) error {
			return errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	err := svc.Delete(1)
	require.Error(t, err)
}

// ── ImportCSV ─────────────────────────────────────────────────────────────────

func TestImportCSV_CreatesValidRows(t *testing.T) {
	var created []*models.TodoItem
	repo := &mockRepo{
		createFn: func(item *models.TodoItem) (*models.TodoItem, error) {
			created = append(created, item)
			return item, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	csvContent := "title,description,is_completed\nBuy milk,Whole milk,true\nBuy eggs,,false\n"
	result, err := svc.ImportCSV(strings.NewReader(csvContent))

	require.NoError(t, err)
	assert.Equal(t, 2, result.Imported)
	assert.Equal(t, 0, result.Failed)
	assert.Empty(t, result.Errors)
	require.Len(t, created, 2)
	assert.Equal(t, "Buy milk", created[0].Title)
	assert.Equal(t, "Whole milk", *created[0].Description)
	assert.True(t, created[0].IsCompleted)
	assert.Equal(t, "Buy eggs", created[1].Title)
	assert.Nil(t, created[1].Description)
	assert.False(t, created[1].IsCompleted)
}

func TestImportCSV_MissingTitle_RecordsError(t *testing.T) {
	repo := &mockRepo{
		createFn: func(item *models.TodoItem) (*models.TodoItem, error) { return item, nil },
	}
	svc := service.NewTodoItemService(repo)

	csvContent := "title,description,is_completed\n,No title,false\nBuy milk,,false\n"
	result, err := svc.ImportCSV(strings.NewReader(csvContent))

	require.NoError(t, err)
	assert.Equal(t, 1, result.Imported)
	assert.Equal(t, 1, result.Failed)
	require.Len(t, result.Errors, 1)
	assert.Equal(t, 2, result.Errors[0].Row)
}

func TestImportCSV_EmptyFile_ReturnsZeroResult(t *testing.T) {
	svc := service.NewTodoItemService(&mockRepo{})

	result, err := svc.ImportCSV(strings.NewReader(""))

	require.NoError(t, err)
	assert.Equal(t, 0, result.Imported)
	assert.Equal(t, 0, result.Failed)
}

func TestImportCSV_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		createFn: func(item *models.TodoItem) (*models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.ImportCSV(strings.NewReader("title,description,is_completed\nBuy milk,,false\n"))
	require.Error(t, err)
}

// ── ExportCSV ─────────────────────────────────────────────────────────────────

func TestExportCSV_WritesHeaderAndRows(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findAllItemsFn: func() ([]models.TodoItem, error) {
			return []models.TodoItem{*item}, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	content, err := svc.ExportCSV()

	require.NoError(t, err)
	assert.Contains(t, content, "id,title,description,is_completed,created_at,updated_at")
	assert.Contains(t, content, item.Title)
}

func TestExportCSV_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		findAllItemsFn: func() ([]models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.ExportCSV()
	require.Error(t, err)
}

// ── ImportExcel ───────────────────────────────────────────────────────────────

func buildExcelFile(t *testing.T, rows [][]string) []byte {
	t.Helper()
	f := excelize.NewFile()
	defer f.Close()
	for r, row := range rows {
		values := make([]interface{}, len(row))
		for c, v := range row {
			values[c] = v
		}
		require.NoError(t, f.SetSheetRow("Sheet1", fmt.Sprintf("A%d", r+1), &values))
	}
	buf, err := f.WriteToBuffer()
	require.NoError(t, err)
	return buf.Bytes()
}

func TestImportExcel_CreatesValidRows(t *testing.T) {
	var created []*models.TodoItem
	repo := &mockRepo{
		createFn: func(item *models.TodoItem) (*models.TodoItem, error) {
			created = append(created, item)
			return item, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	content := buildExcelFile(t, [][]string{
		{"title", "description", "is_completed"},
		{"Buy milk", "Whole milk", "true"},
		{"Buy eggs", "", "false"},
	})
	result, err := svc.ImportExcel(bytes.NewReader(content))

	require.NoError(t, err)
	assert.Equal(t, 2, result.Imported)
	assert.Equal(t, 0, result.Failed)
	assert.Empty(t, result.Errors)
	require.Len(t, created, 2)
	assert.Equal(t, "Buy milk", created[0].Title)
	assert.Equal(t, "Whole milk", *created[0].Description)
	assert.True(t, created[0].IsCompleted)
	assert.Equal(t, "Buy eggs", created[1].Title)
	assert.Nil(t, created[1].Description)
	assert.False(t, created[1].IsCompleted)
}

func TestImportExcel_MissingTitle_RecordsError(t *testing.T) {
	repo := &mockRepo{
		createFn: func(item *models.TodoItem) (*models.TodoItem, error) { return item, nil },
	}
	svc := service.NewTodoItemService(repo)

	content := buildExcelFile(t, [][]string{
		{"title", "description", "is_completed"},
		{"", "No title", "false"},
		{"Buy milk", "", "false"},
	})
	result, err := svc.ImportExcel(bytes.NewReader(content))

	require.NoError(t, err)
	assert.Equal(t, 1, result.Imported)
	assert.Equal(t, 1, result.Failed)
	require.Len(t, result.Errors, 1)
	assert.Equal(t, 2, result.Errors[0].Row)
}

func TestImportExcel_EmptyFile_ReturnsZeroResult(t *testing.T) {
	svc := service.NewTodoItemService(&mockRepo{})

	content := buildExcelFile(t, [][]string{{"title", "description", "is_completed"}})
	result, err := svc.ImportExcel(bytes.NewReader(content))

	require.NoError(t, err)
	assert.Equal(t, 0, result.Imported)
	assert.Equal(t, 0, result.Failed)
}

func TestImportExcel_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		createFn: func(item *models.TodoItem) (*models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	content := buildExcelFile(t, [][]string{
		{"title", "description", "is_completed"},
		{"Buy milk", "", "false"},
	})
	_, err := svc.ImportExcel(bytes.NewReader(content))
	require.Error(t, err)
}

// ── ExportExcel ───────────────────────────────────────────────────────────────

func TestExportExcel_WritesHeaderAndRows(t *testing.T) {
	item := sampleItem()
	repo := &mockRepo{
		findAllItemsFn: func() ([]models.TodoItem, error) {
			return []models.TodoItem{*item}, nil
		},
	}
	svc := service.NewTodoItemService(repo)

	content, err := svc.ExportExcel()

	require.NoError(t, err)
	f, err := excelize.OpenReader(bytes.NewReader(content))
	require.NoError(t, err)
	defer f.Close()
	rows, err := f.GetRows(f.GetSheetName(0))
	require.NoError(t, err)
	require.Len(t, rows, 2)
	assert.Equal(t, []string{"id", "title", "description", "is_completed", "created_at", "updated_at"}, rows[0])
	assert.Equal(t, item.Title, rows[1][1])
}

func TestExportExcel_RepoError_Propagates(t *testing.T) {
	repo := &mockRepo{
		findAllItemsFn: func() ([]models.TodoItem, error) {
			return nil, errors.New("db error")
		},
	}
	svc := service.NewTodoItemService(repo)

	_, err := svc.ExportExcel()
	require.Error(t, err)
}
