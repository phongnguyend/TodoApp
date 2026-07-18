package com.example.todo.repository;

import com.example.todo.entity.TodoItemAttachment;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;
import java.util.Optional;

@Repository
public interface TodoItemAttachmentRepository extends JpaRepository<TodoItemAttachment, Long> {
    List<TodoItemAttachment> findByTodoItemIdOrderByCreatedAtAsc(Long todoItemId);

    Optional<TodoItemAttachment> findByIdAndTodoItemId(Long id, Long todoItemId);

    Optional<TodoItemAttachment> findByTodoItemIdAndFileId(Long todoItemId, Long fileId);
}
