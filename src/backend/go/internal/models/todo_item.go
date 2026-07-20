package models

import (
	"time"

	"gorm.io/gorm"
)

// TodoItem is the GORM entity - analogous to an EF Core entity class.
// GORM maps struct fields to database columns using conventions and struct tags.
type TodoItem struct {
	ID              uint           `gorm:"primarykey;autoIncrement"          json:"id"`
	Title           string         `gorm:"type:varchar(200);not null"        json:"title"`
	Description     *string        `gorm:"type:text"                         json:"description"`
	IsCompleted     bool           `gorm:"default:false;not null"            json:"isCompleted"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"                    json:"createdAt"`
	CreatedByUserID *uint          `gorm:"index"                             json:"createdByUserId"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime"                    json:"updatedAt"`
	UpdatedByUserID *uint          `gorm:"index"                             json:"updatedByUserId"`
	CreatedByUser   *User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:CreatedByUserID" json:"-"`
	UpdatedByUser   *User          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UpdatedByUserID" json:"-"`
	DeletedAt       gorm.DeletedAt `gorm:"index"                             json:"-"` // soft delete
}

// TableName overrides the default table name - mirrors [Table("todo_items")] in EF.
func (TodoItem) TableName() string {
	return "todo_items"
}
