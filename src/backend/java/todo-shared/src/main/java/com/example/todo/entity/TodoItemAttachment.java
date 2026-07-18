package com.example.todo.entity;

import jakarta.persistence.*;
import lombok.Getter;
import lombok.NoArgsConstructor;
import lombok.Setter;
import org.hibernate.annotations.CreationTimestamp;
import org.hibernate.annotations.UpdateTimestamp;

import java.time.Instant;

@Entity
@Table(name = "todo_item_attachments",
        uniqueConstraints = @UniqueConstraint(
                name = "uq_todo_item_attachments_todo_file",
                columnNames = {"todo_item_id", "file_id"}))
@Getter
@Setter
@NoArgsConstructor
public class TodoItemAttachment {

    @Id
    @GeneratedValue(strategy = GenerationType.IDENTITY)
    private Long id;

    @Column(name = "todo_item_id", nullable = false)
    private Long todoItemId;

    @Column(name = "file_id", nullable = false)
    private Long fileId;

    @CreationTimestamp
    @Column(name = "created_at", nullable = false, updatable = false)
    private Instant createdAt;

    @UpdateTimestamp
    @Column(name = "updated_at")
    private Instant updatedAt;

    public TodoItemAttachment(Long todoItemId, Long fileId) {
        this.todoItemId = todoItemId;
        this.fileId = fileId;
    }
}
