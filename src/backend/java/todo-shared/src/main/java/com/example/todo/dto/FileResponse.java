package com.example.todo.dto;

import com.example.todo.entity.FileEntity;

import java.time.Instant;

/**
 * Response DTO - mirrors a FileDto / view model in C#.
 * Note: the on-disk {@code location} is intentionally not exposed to clients;
 * content is retrieved via the dedicated download endpoint instead.
 */
public record FileResponse(
        Long id,
        String name,
        String extension,
        long size,
        String contentType,
        Instant createdAt,
        Long createdByUserId,
        Instant updatedAt,
        Long updatedByUserId
) {
    public FileResponse(Long id, String name, String extension, long size, String contentType,
            Instant createdAt, Instant updatedAt) {
        this(id, name, extension, size, contentType, createdAt, null, updatedAt, null);
    }

    public static FileResponse from(FileEntity entity) {
        return new FileResponse(
                entity.getId(),
                entity.getName(),
                entity.getExtension(),
                entity.getSize(),
                entity.getContentType(),
                entity.getCreatedAt(),
                entity.getCreatedByUserId(),
                entity.getUpdatedAt(),
                entity.getUpdatedByUserId()
        );
    }
}
