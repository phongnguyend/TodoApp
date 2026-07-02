package service_test

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"github.com/todo/backend/go/internal/service"
	"gorm.io/gorm"
)

// ── mock repository ───────────────────────────────────────────────────────────

type mockFileRepo struct {
	findAllFn  func(skip, limit int) (repository.FilePaginatedResult, error)
	findByIDFn func(id uint) (*models.File, error)
	createFn   func(file *models.File) (*models.File, error)
	deleteFn   func(file *models.File) error
}

func (m *mockFileRepo) FindAll(skip, limit int) (repository.FilePaginatedResult, error) {
	return m.findAllFn(skip, limit)
}
func (m *mockFileRepo) FindByID(id uint) (*models.File, error) {
	return m.findByIDFn(id)
}
func (m *mockFileRepo) Create(file *models.File) (*models.File, error) {
	return m.createFn(file)
}
func (m *mockFileRepo) Delete(file *models.File) error {
	return m.deleteFn(file)
}

// ── helpers ───────────────────────────────────────────────────────────────────

func sampleFile() *models.File {
	ct := "text/plain"
	return &models.File{
		ID:          1,
		Name:        "notes.txt",
		Extension:   "txt",
		Size:        5,
		ContentType: &ct,
		Location:    "notes.txt", // overridden per-test as needed
		CreatedAt:   time.Now(),
	}
}

func newFileService(repo repository.FileRepository, storageDir string) service.FileService {
	cfg := &config.Config{FileStoragePath: storageDir, MaxUploadSizeBytes: 1024}
	return service.NewFileService(repo, cfg)
}

// ── GetAll ────────────────────────────────────────────────────────────────────

func TestFileGetAll_ReturnsPaginatedItems(t *testing.T) {
	file := sampleFile()
	repo := &mockFileRepo{
		findAllFn: func(skip, limit int) (repository.FilePaginatedResult, error) {
			assert.Equal(t, 0, skip)
			assert.Equal(t, 10, limit)
			return repository.FilePaginatedResult{Items: []models.File{*file}, Total: 25}, nil
		},
	}
	svc := newFileService(repo, t.TempDir())

	result, err := svc.GetAll(1, 10)

	require.NoError(t, err)
	assert.Equal(t, 1, result.Page)
	assert.Equal(t, 10, result.PageSize)
	assert.Equal(t, int64(25), result.Total)
	assert.Equal(t, 3, result.TotalPages)
	assert.Len(t, result.Items, 1)
	assert.Equal(t, file.Name, result.Items[0].Name)
}

func TestFileGetAll_RepoError_Propagates(t *testing.T) {
	repo := &mockFileRepo{
		findAllFn: func(skip, limit int) (repository.FilePaginatedResult, error) {
			return repository.FilePaginatedResult{}, errors.New("db error")
		},
	}
	svc := newFileService(repo, t.TempDir())

	_, err := svc.GetAll(1, 10)
	require.Error(t, err)
}

// ── GetByID ───────────────────────────────────────────────────────────────────

func TestFileGetByID_ReturnsFile(t *testing.T) {
	file := sampleFile()
	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) {
			assert.Equal(t, uint(1), id)
			return file, nil
		},
	}
	svc := newFileService(repo, t.TempDir())

	result, err := svc.GetByID(1)

	require.NoError(t, err)
	assert.Equal(t, file.Name, result.Name)
	assert.Equal(t, file.Extension, result.Extension)
}

func TestFileGetByID_NotFound_ReturnsErrFileNotFound(t *testing.T) {
	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) {
			return nil, gorm.ErrRecordNotFound
		},
	}
	svc := newFileService(repo, t.TempDir())

	_, err := svc.GetByID(99)
	require.ErrorIs(t, err, service.ErrFileNotFound)
}

// ── Upload ────────────────────────────────────────────────────────────────────

func TestUpload_StoresContentAndMetadata(t *testing.T) {
	dir := t.TempDir()
	var created *models.File
	repo := &mockFileRepo{
		createFn: func(file *models.File) (*models.File, error) {
			file.ID = 1
			created = file
			return file, nil
		},
	}
	svc := newFileService(repo, dir)

	content := "hello world"
	result, err := svc.Upload(service.UploadInput{
		OriginalName: "notes.txt",
		ContentType:  "text/plain",
		Size:         int64(len(content)),
		Reader:       strings.NewReader(content),
	})

	require.NoError(t, err)
	assert.Equal(t, "notes.txt", result.Name)
	assert.Equal(t, "txt", result.Extension)
	assert.Equal(t, int64(len(content)), result.Size)
	require.NotNil(t, result.ContentType)
	assert.Equal(t, "text/plain", *result.ContentType)

	// Content was written to disk under the configured storage directory.
	require.NotNil(t, created)
	assert.True(t, strings.HasPrefix(created.Location, dir))
	data, readErr := os.ReadFile(created.Location)
	require.NoError(t, readErr)
	assert.Equal(t, content, string(data))
}

func TestUpload_StripsDirectoryComponents(t *testing.T) {
	dir := t.TempDir()
	var created *models.File
	repo := &mockFileRepo{
		createFn: func(file *models.File) (*models.File, error) {
			file.ID = 1
			created = file
			return file, nil
		},
	}
	svc := newFileService(repo, dir)

	content := "data"
	result, err := svc.Upload(service.UploadInput{
		OriginalName: "../../etc/passwd",
		Size:         int64(len(content)),
		Reader:       strings.NewReader(content),
	})

	require.NoError(t, err)
	assert.Equal(t, "passwd", result.Name)
	require.NotNil(t, created)
	assert.Equal(t, dir, filepath.Dir(created.Location))
}

func TestUpload_ExceedsMaxSize_ReturnsErrFileTooLarge(t *testing.T) {
	repo := &mockFileRepo{}
	svc := newFileService(repo, t.TempDir())

	_, err := svc.Upload(service.UploadInput{
		OriginalName: "big.bin",
		Size:         2048, // exceeds the 1024-byte test limit
		Reader:       strings.NewReader("irrelevant"),
	})

	require.ErrorIs(t, err, service.ErrFileTooLarge)
}

func TestUpload_RepoError_Propagates(t *testing.T) {
	repo := &mockFileRepo{
		createFn: func(file *models.File) (*models.File, error) {
			return nil, errors.New("db error")
		},
	}
	svc := newFileService(repo, t.TempDir())

	_, err := svc.Upload(service.UploadInput{
		OriginalName: "notes.txt",
		Size:         4,
		Reader:       strings.NewReader("data"),
	})
	require.Error(t, err)
}

// ── GetDownloadTarget ─────────────────────────────────────────────────────────

func TestGetDownloadTarget_ReturnsPathAndMetadata(t *testing.T) {
	dir := t.TempDir()
	location := filepath.Join(dir, "stored.txt")
	require.NoError(t, os.WriteFile(location, []byte("content"), 0o644))

	file := sampleFile()
	file.Location = location

	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) { return file, nil },
	}
	svc := newFileService(repo, dir)

	target, err := svc.GetDownloadTarget(1)

	require.NoError(t, err)
	assert.Equal(t, location, target.Path)
	assert.Equal(t, file.Name, target.Name)
	assert.Equal(t, "text/plain", target.ContentType)
}

func TestGetDownloadTarget_NotFound_ReturnsErrFileNotFound(t *testing.T) {
	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) { return nil, gorm.ErrRecordNotFound },
	}
	svc := newFileService(repo, t.TempDir())

	_, err := svc.GetDownloadTarget(99)
	require.ErrorIs(t, err, service.ErrFileNotFound)
}

func TestGetDownloadTarget_MissingOnDisk_ReturnsErrFileContentMissing(t *testing.T) {
	file := sampleFile()
	file.Location = filepath.Join(t.TempDir(), "missing.txt")

	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) { return file, nil },
	}
	svc := newFileService(repo, t.TempDir())

	_, err := svc.GetDownloadTarget(1)
	require.ErrorIs(t, err, service.ErrFileContentMissing)
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestFileDelete_RemovesRowAndDiskContent(t *testing.T) {
	dir := t.TempDir()
	location := filepath.Join(dir, "stored.txt")
	require.NoError(t, os.WriteFile(location, []byte("content"), 0o644))

	file := sampleFile()
	file.Location = location
	deleted := false

	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) { return file, nil },
		deleteFn: func(in *models.File) error {
			deleted = true
			return nil
		},
	}
	svc := newFileService(repo, dir)

	err := svc.Delete(1)

	require.NoError(t, err)
	assert.True(t, deleted)
	_, statErr := os.Stat(location)
	assert.True(t, os.IsNotExist(statErr))
}

func TestFileDelete_NotFound_ReturnsErrFileNotFound(t *testing.T) {
	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) { return nil, gorm.ErrRecordNotFound },
	}
	svc := newFileService(repo, t.TempDir())

	err := svc.Delete(99)
	require.ErrorIs(t, err, service.ErrFileNotFound)
}

func TestFileDelete_RepoDeleteError_Propagates(t *testing.T) {
	file := sampleFile()
	file.Location = filepath.Join(t.TempDir(), "stored.txt")

	repo := &mockFileRepo{
		findByIDFn: func(id uint) (*models.File, error) { return file, nil },
		deleteFn: func(in *models.File) error {
			return errors.New("db error")
		},
	}
	svc := newFileService(repo, t.TempDir())

	err := svc.Delete(1)
	require.Error(t, err)
}
