package com.example.todo.dto;

import io.swagger.v3.oas.annotations.media.Schema;
import jakarta.validation.constraints.NotNull;
import jakarta.validation.constraints.Positive;

public record SaveTodoItemAttachmentRequest(
        @Schema(example = "5")
        @NotNull(message = "File ID is required")
        @Positive(message = "File ID must be positive")
        Long fileId
) {}
