package com.example.todo.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

/**
 * Request DTO for creating a todo item.
 * Mirrors a CreateTodoItemRequest command model in C#.
 * Jakarta Bean Validation annotations mirror [Required] / [MaxLength] Data Annotations.
 */
public record CreateTodoItemRequest(

        @Schema(example = "Buy groceries")
        @NotBlank(message = "Title must not be blank")
        @Size(max = 200, message = "Title must not exceed 200 characters")
        String title,

        @Schema(example = "Milk, eggs, bread")
        @Size(max = 2000, message = "Description must not exceed 2000 characters")
        String description
) {}
