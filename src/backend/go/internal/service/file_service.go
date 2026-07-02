package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/todo/backend/go/internal/config"
	"github.com/todo/backend/go/internal/dto"
	"github.com/todo/backend/go/internal/models"
	"github.com/todo/backend/go/internal/repository"
	"gorm.io/gorm"
)

// ErrFileNotFound is returned when a requested file does not exist.
var ErrFileNotFound = errors.New("file not found")

// ErrFileTooLarge is returned when an uploaded file exceeds the configured maximum size.
var ErrFileTooLarge = errors.New("file exceeds the maximum allowed upload size")

// ErrFileContentMissing is returned when a file's metadata exists but its content is missing on disk.
var ErrFileContentMissing = errors.New("file content not found on disk")

// UploadInput carries the raw data needed to persist an uploaded file - decouples
// the service from Gin's multipart.FileHeader type.
type UploadInput struct {
	OriginalName string
	ContentType  string
	Size         int64
	Reader       io.Reader
}

// DownloadTarget carries everything needed to stream a file back to the client.
type DownloadTarget struct {
	Path        string
	Name        string
	ContentType string
}

// FileService defines the business-logic contract for uploaded files.
// Mirrors IFileService in C#.
type FileService interface {
	GetAll(page, pageSize int) (dto.PaginatedResponse[dto.FileResponse], error)
	GetByID(id uint) (dto.FileResponse, error)
	Upload(input UploadInput) (dto.FileResponse, error)
	GetDownloadTarget(id uint) (DownloadTarget, error)
	Delete(id uint) error
}

type fileService struct {
	repo          repository.FileRepository
	storageDir    string
	maxUploadSize int64
}

// NewFileService constructs the service with its repository dependency injected.
func NewFileService(repo repository.FileRepository, cfg *config.Config) FileService {
	return &fileService{
		repo:          repo,
		storageDir:    cfg.FileStoragePath,
		maxUploadSize: cfg.MaxUploadSizeBytes,
	}
}

// ── Mapping ───────────────────────────────────────────────────────────────────

func fileToResponse(m *models.File) dto.FileResponse {
	return dto.FileResponse{
		ID:          m.ID,
		Name:        m.Name,
		Extension:   m.Extension,
		Size:        m.Size,
		ContentType: m.ContentType,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toFilePaginated(result repository.FilePaginatedResult, page, pageSize int) dto.PaginatedResponse[dto.FileResponse] {
	items := make([]dto.FileResponse, len(result.Items))
	for i := range result.Items {
		items[i] = fileToResponse(&result.Items[i])
	}
	return dto.PaginatedResponse[dto.FileResponse]{
		Items:      items,
		Total:      result.Total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int(math.Ceil(float64(result.Total) / float64(pageSize))),
	}
}

func (s *fileService) getOrNotFound(id uint) (*models.File, error) {
	file, err := s.repo.FindByID(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrFileNotFound
	}
	return file, err
}

// ── Queries ───────────────────────────────────────────────────────────────────

func (s *fileService) GetAll(page, pageSize int) (dto.PaginatedResponse[dto.FileResponse], error) {
	result, err := s.repo.FindAll((page-1)*pageSize, pageSize)
	if err != nil {
		return dto.PaginatedResponse[dto.FileResponse]{}, err
	}
	return toFilePaginated(result, page, pageSize), nil
}

func (s *fileService) GetByID(id uint) (dto.FileResponse, error) {
	file, err := s.getOrNotFound(id)
	if err != nil {
		return dto.FileResponse{}, err
	}
	return fileToResponse(file), nil
}

// ── Commands ──────────────────────────────────────────────────────────────────

func (s *fileService) Upload(input UploadInput) (dto.FileResponse, error) {
	if input.Size > s.maxUploadSize {
		return dto.FileResponse{}, ErrFileTooLarge
	}

	// Strip any directory components from the client-supplied name to prevent path traversal.
	originalName := filepath.Base(input.OriginalName)
	if originalName == "" || originalName == "." || originalName == string(filepath.Separator) {
		originalName = "unnamed"
	}
	extension := strings.TrimPrefix(filepath.Ext(originalName), ".")

	if err := os.MkdirAll(s.storageDir, 0o755); err != nil {
		return dto.FileResponse{}, err
	}

	prefix, err := randomHex(16)
	if err != nil {
		return dto.FileResponse{}, err
	}
	// A random prefix avoids collisions/overwrites between uploads that share a name.
	location := filepath.Join(s.storageDir, prefix+"_"+originalName)

	out, err := os.Create(location)
	if err != nil {
		return dto.FileResponse{}, err
	}
	defer out.Close()

	written, err := io.Copy(out, input.Reader)
	if err != nil {
		return dto.FileResponse{}, err
	}

	var contentType *string
	if input.ContentType != "" {
		contentType = &input.ContentType
	}

	file := &models.File{
		Name:        originalName,
		Extension:   extension,
		Size:        written,
		ContentType: contentType,
		Location:    location,
	}
	created, err := s.repo.Create(file)
	if err != nil {
		return dto.FileResponse{}, err
	}
	return fileToResponse(created), nil
}

func (s *fileService) GetDownloadTarget(id uint) (DownloadTarget, error) {
	file, err := s.getOrNotFound(id)
	if err != nil {
		return DownloadTarget{}, err
	}
	if _, err := os.Stat(file.Location); err != nil {
		return DownloadTarget{}, ErrFileContentMissing
	}

	contentType := "application/octet-stream"
	if file.ContentType != nil && *file.ContentType != "" {
		contentType = *file.ContentType
	}
	return DownloadTarget{Path: file.Location, Name: file.Name, ContentType: contentType}, nil
}

func (s *fileService) Delete(id uint) error {
	file, err := s.getOrNotFound(id)
	if err != nil {
		return err
	}
	if err := s.repo.Delete(file); err != nil {
		return err
	}
	if _, statErr := os.Stat(file.Location); statErr == nil {
		return os.Remove(file.Location)
	}
	return nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
