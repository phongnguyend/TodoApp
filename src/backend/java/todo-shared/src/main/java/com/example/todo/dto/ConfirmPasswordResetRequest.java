package com.example.todo.dto;

import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record ConfirmPasswordResetRequest(
        @NotBlank String token,
        @NotBlank @Size(min = 8, max = 128) String newPassword) {
}
