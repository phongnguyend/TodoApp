package com.example.todo.service;

import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.ImportResult;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import org.springframework.web.multipart.MultipartFile;

/**
 * Service interface - mirrors ITodoItemService in C#.
 * Defines the contract for the business-logic layer.
 */
public interface TodoItemService {

    PaginatedResponse<TodoItemResponse> getAll(int page, int pageSize);

    PaginatedResponse<TodoItemResponse> getIncomplete(int page, int pageSize);

    TodoItemResponse getById(Long id);

    TodoItemResponse create(CreateTodoItemRequest request);

    TodoItemResponse update(Long id, UpdateTodoItemRequest request);

    TodoItemResponse markComplete(Long id);

    void delete(Long id);

    ImportResult importCsv(MultipartFile file);

    String exportCsv();
}
