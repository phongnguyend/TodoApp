package repository

import (
	"time"

	"github.com/todo/backend/go/internal/models"
	"gorm.io/gorm"
)

type emailLogRepository struct {
	db *gorm.DB
}

// NewEmailLogRepository creates a new EmailLog repository.
func NewEmailLogRepository(db *gorm.DB) EmailLogRepository {
	return &emailLogRepository{db: db}
}

func (r *emailLogRepository) Create(log *models.EmailLog) (*models.EmailLog, error) {
	if err := r.db.Create(log).Error; err != nil {
		return nil, err
	}
	return log, nil
}

func (r *emailLogRepository) MarkSent(log *models.EmailLog) error {
	now := time.Now().UTC()
	return r.db.Model(log).Updates(map[string]any{
		"status":  "sent",
		"sent_at": now,
	}).Error
}

func (r *emailLogRepository) MarkFailed(log *models.EmailLog, errMsg string) error {
	return r.db.Model(log).Updates(map[string]any{
		"status":        "failed",
		"error_message": errMsg,
	}).Error
}
