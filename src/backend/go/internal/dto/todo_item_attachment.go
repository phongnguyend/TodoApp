package dto

import "time"

type SaveTodoItemAttachmentRequest struct {
	FileID uint `json:"fileId" binding:"required"`
}

type TodoItemAttachmentResponse struct {
	ID              uint       `json:"id"`
	TodoItemID      uint       `json:"todoItemId"`
	FileID          uint       `json:"fileId"`
	CreatedAt       time.Time  `json:"createdAt"`
	CreatedByUserID *uint      `json:"createdByUserId"`
	UpdatedAt       *time.Time `json:"updatedAt"`
	UpdatedByUserID *uint      `json:"updatedByUserId"`
}
