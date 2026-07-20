package com.example.todo.entity;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.Instant;

/**
 * JPA entity for an uploaded file's metadata - analogous to an EF Core entity class.
 * Named {@code FileEntity} (rather than {@code File}) to avoid clashing with {@link java.io.File}.
 * The actual file content is stored on disk at the path recorded in {@link #location};
 * only metadata is persisted in the {@code files} table.
 */
@Entity
@Table(name = "files")
@Getter
@Setter
@NoArgsConstructor
public class FileEntity {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, length = 255)
    private String name;

    @Column(nullable = false, length = 20)
    private String extension;

    @Column(nullable = false)
    private long size;

    @Column(name = "content_type", length = 100)
    private String contentType;

    @Column(nullable = false, length = 500)
    private String location;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @Column(name = "created_by_user_id")
    private Long createdByUserId;

    @UpdateTimestamp
    @Column(name = "updated_at")
    private Instant updatedAt;

    @Column(name = "updated_by_user_id")
    private Long updatedByUserId;

    public FileEntity(String name, String extension, long size, String contentType, String location) {
        this.name = name;
        this.extension = extension;
        this.size = size;
        this.contentType = contentType;
        this.location = location;
    }
}
