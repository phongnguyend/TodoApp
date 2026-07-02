package repository

import (
	"github.com/todo/backend/go/internal/models"
	"gorm.io/gorm"
)

// fileRepository is the GORM-backed implementation of FileRepository.
// Mirrors a repository class backed by EF Core's DbSet<File>.
type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new repository - called from the DI composition root.
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) FindAll(skip, limit int) (FilePaginatedResult, error) {
	var items []models.File
	var total int64

	if err := r.db.Model(&models.File{}).Count(&total).Error; err != nil {
		return FilePaginatedResult{}, err
	}
	if err := r.db.Offset(skip).Limit(limit).Order("created_at desc").Find(&items).Error; err != nil {
		return FilePaginatedResult{}, err
	}
	return FilePaginatedResult{Items: items, Total: total}, nil
}

func (r *fileRepository) FindByID(id uint) (*models.File, error) {
	var file models.File
	if err := r.db.First(&file, id).Error; err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) Create(file *models.File) (*models.File, error) {
	if err := r.db.Create(file).Error; err != nil {
		return nil, err
	}
	return file, nil
}

func (r *fileRepository) Delete(file *models.File) error {
	return r.db.Delete(file).Error
}
