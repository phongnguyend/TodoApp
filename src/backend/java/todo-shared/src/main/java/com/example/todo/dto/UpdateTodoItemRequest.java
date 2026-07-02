package com.example.todo.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.Size;

/**
 * Request DTO for updating a todo item - all fields are optional (PATCH semantics).
 * Mirrors an UpdateTodoItemRequest in C#.
 */
public record UpdateTodoItemRequest(

        @Schema(example = "Buy groceries")
        @Size(max = 200, message = "Title must not exceed 200 characters")
        String title,

        @Schema(example = "Milk, eggs, bread")
        @Size(max = 2000, message = "Description must not exceed 2000 characters")
        String description,

        Boolean isCompleted
) {}
