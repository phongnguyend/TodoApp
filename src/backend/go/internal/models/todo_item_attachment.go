package models

import "time"

// TodoItemAttachment links an existing uploaded file to a todo item.
type TodoItemAttachment struct {
	ID         uint       `gorm:"primarykey;autoIncrement" json:"id"`
	TodoItemID uint       `gorm:"not null;uniqueIndex:uq_todo_item_attachments_todo_file" json:"todoItemId"`
	FileID     uint       `gorm:"not null;uniqueIndex:uq_todo_item_attachments_todo_file" json:"fileId"`
	CreatedAt  time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  *time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	TodoItem   TodoItem   `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:TodoItemID" json:"-"`
	File       File       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:FileID" json:"-"`
}

func (TodoItemAttachment) TableName() string { return "todo_item_attachments" }
