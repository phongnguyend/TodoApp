package com.example.todo.dto;

import jakarta.validation.constraints.Email;
import jakarta.validation.constraints.Pattern;
import jakarta.validation.constraints.Size;

public record UpdateUserRequest(
        @Pattern(regexp = ".*\\S.*", message = "must not be blank") @Size(max = 50) String username,
        @Email @Size(max = 255) String email,
        @Size(min = 8, max = 128) String password) {
}
