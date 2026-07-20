package com.example.todo.dto;

import com.example.todo.entity.TodoItemAttachment;

import java.time.Instant;

public record TodoItemAttachmentResponse(
        Long id,
        Long todoItemId,
        Long fileId,
        Instant createdAt,
        Long createdByUserId,
        Instant updatedAt,
        Long updatedByUserId
) {
    public TodoItemAttachmentResponse(Long id, Long todoItemId, Long fileId,
            Instant createdAt, Instant updatedAt) {
        this(id, todoItemId, fileId, createdAt, null, updatedAt, null);
    }

    public static TodoItemAttachmentResponse from(TodoItemAttachment attachment) {
        return new TodoItemAttachmentResponse(
                attachment.getId(), attachment.getTodoItemId(), attachment.getFileId(),
                attachment.getCreatedAt(), attachment.getCreatedByUserId(),
                attachment.getUpdatedAt(), attachment.getUpdatedByUserId());
    }
}
