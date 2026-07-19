package com.example.todo.dto;

import com.example.todo.entity.User;

import java.time.Instant;

public record UserResponse(
        Long id,
        String username,
        String email,
        boolean isActive,
        Instant createdAt,
        Instant updatedAt) {
    public static UserResponse from(User user) {
        return new UserResponse(user.getId(), user.getUsername(), user.getEmail(), user.isActive(),
                user.getCreatedAt(), user.getUpdatedAt());
    }
}
