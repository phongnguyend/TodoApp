package com.example.todo.dto;

import com.example.todo.entity.User;

import java.time.Instant;

public record UserResponse(
        Long id,
        String username,
        String email,
        boolean isActive,
        Instant createdAt,
        Long createdByUserId,
        Instant updatedAt,
        Long updatedByUserId) {
    public UserResponse(Long id, String username, String email, boolean isActive,
            Instant createdAt, Instant updatedAt) {
        this(id, username, email, isActive, createdAt, null, updatedAt, null);
    }

    public static UserResponse from(User user) {
        return new UserResponse(user.getId(), user.getUsername(), user.getEmail(), user.isActive(),
                user.getCreatedAt(), user.getCreatedByUserId(), user.getUpdatedAt(), user.getUpdatedByUserId());
    }
}
