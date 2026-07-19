package com.example.todo.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.NotBlank;
import jakarta.validation.constraints.Size;

public record TokenRequest(
        @NotBlank @Email @Size(max = 255) String email,
        @NotBlank @Size(max = 128) String password) {}
