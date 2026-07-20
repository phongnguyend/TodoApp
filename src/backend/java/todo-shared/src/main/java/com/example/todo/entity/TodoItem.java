package com.example.todo.entity;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.Instant;

/**
 * JPA entity - analogous to an EF Core entity class mapped via Data Annotations or Fluent API.
 * Hibernate manages the ORM mapping; Flyway manages the schema.
 */
@Entity
@Table(name = "todo_items")
@Getter
@Setter
@NoArgsConstructor
public class TodoItem {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(nullable = false, length = 200)
    private String title;

    @Column(columnDefinition = "TEXT")
    private String description;

    @Column(name = "is_completed", nullable = false)
    private boolean completed = false;

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

    public TodoItem(String title, String description) {
        this.title = title;
        this.description = description;
    }
}
