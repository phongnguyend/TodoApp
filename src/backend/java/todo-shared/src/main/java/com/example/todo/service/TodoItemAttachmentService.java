package com.example.todo.service;

import com.example.todo.dto.SaveTodoItemAttachmentRequest;
import com.example.todo.dto.TodoItemAttachmentResponse;

import java.util.List;

public interface TodoItemAttachmentService {
    List<TodoItemAttachmentResponse> getAll(Long todoItemId);

    TodoItemAttachmentResponse getById(Long todoItemId, Long attachmentId);

    TodoItemAttachmentResponse create(Long todoItemId, SaveTodoItemAttachmentRequest request);

    TodoItemAttachmentResponse update(Long todoItemId, Long attachmentId, SaveTodoItemAttachmentRequest request);

    void delete(Long todoItemId, Long attachmentId);
}
