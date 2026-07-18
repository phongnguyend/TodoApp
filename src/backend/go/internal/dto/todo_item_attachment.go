package dto

import "time"

type SaveTodoItemAttachmentRequest struct {
	FileID uint `json:"fileId" binding:"required"`
}

type TodoItemAttachmentResponse struct {
	ID         uint       `json:"id"`
	TodoItemID uint       `json:"todoItemId"`
	FileID     uint       `json:"fileId"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  *time.Time `json:"updatedAt"`
}
