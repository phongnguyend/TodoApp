package com.example.todo.dto;

import com.example.todo.entity.TodoItemAttachment;

import java.time.Instant;

public record TodoItemAttachmentResponse(
        Long id,
        Long todoItemId,
        Long fileId,
        Instant createdAt,
        Instant updatedAt
) {
    public static TodoItemAttachmentResponse from(TodoItemAttachment attachment) {
        return new TodoItemAttachmentResponse(
                attachment.getId(), attachment.getTodoItemId(), attachment.getFileId(),
                attachment.getCreatedAt(), attachment.getUpdatedAt());
    }
}
