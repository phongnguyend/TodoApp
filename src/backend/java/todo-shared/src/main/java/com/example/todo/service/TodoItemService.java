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
    default TodoItemResponse create(CreateTodoItemRequest request, Long actorUserId) { return create(request); }

    TodoItemResponse update(Long id, UpdateTodoItemRequest request);
    default TodoItemResponse update(Long id, UpdateTodoItemRequest request, Long actorUserId) { return update(id, request); }

    TodoItemResponse markComplete(Long id);
    default TodoItemResponse markComplete(Long id, Long actorUserId) { return markComplete(id); }

    void delete(Long id);

    ImportResult importCsv(MultipartFile file);
    default ImportResult importCsv(MultipartFile file, Long actorUserId) { return importCsv(file); }

    String exportCsv();

    ImportResult importExcel(MultipartFile file);
    default ImportResult importExcel(MultipartFile file, Long actorUserId) { return importExcel(file); }

    byte[] exportExcel();
}
