package com.example.todo.dto;

import com.example.todo.entity.TodoItem;

import java.time.Instant;

/**
 * Response DTO - mirrors a TodoItemDto / view model in C#.
 * Uses a static factory method instead of AutoMapper for explicit mapping.
 */
public record TodoItemResponse(
        Long id,
        String title,
        String description,
        boolean isCompleted,
        Instant createdAt,
        Instant updatedAt
) {
    public static TodoItemResponse from(TodoItem entity) {
        return new TodoItemResponse(
                entity.getId(),
                entity.getTitle(),
                entity.getDescription(),
                entity.isCompleted(),
                entity.getCreatedAt(),
                entity.getUpdatedAt()
        );
    }
}
