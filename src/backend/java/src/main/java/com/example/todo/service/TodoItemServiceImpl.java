package com.example.todo.service;

import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import com.example.todo.entity.TodoItem;
import com.example.todo.repository.TodoItemRepository;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.util.List;

/**
 * Service implementation — registered as a Spring-managed bean via @Service.
 * Mirrors a service class injected through ASP.NET Core's DI container.
 *
 * @Transactional mirrors [UnitOfWork] / SaveChanges() semantics in EF Core.
 */
@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TodoItemServiceImpl implements TodoItemService {

    private final TodoItemRepository repository;

    // ── Helpers ────────────────────────────────────────────────────────────────

    private TodoItem getOrThrow(Long id) {
        return repository.findById(id)
                .orElseThrow(() -> new EntityNotFoundException("Todo item " + id + " not found."));
    }

    private static Pageable toPageable(int page, int pageSize) {
        return PageRequest.of(page - 1, pageSize, Sort.by("createdAt").descending());
    }

    private static PaginatedResponse<TodoItemResponse> toPaginated(Page<TodoItem> pageResult, int page, int pageSize) {
        List<TodoItemResponse> items = pageResult.getContent()
                .stream()
                .map(TodoItemResponse::from)
                .toList();
        return new PaginatedResponse<>(items, pageResult.getTotalElements(), page, pageSize, pageResult.getTotalPages());
    }

    // ── Queries ────────────────────────────────────────────────────────────────

    @Override
    public PaginatedResponse<TodoItemResponse> getAll(int page, int pageSize) {
        Page<TodoItem> result = repository.findAll(toPageable(page, pageSize));
        return toPaginated(result, page, pageSize);
    }

    @Override
    public PaginatedResponse<TodoItemResponse> getIncomplete(int page, int pageSize) {
        Page<TodoItem> result = repository.findByCompletedFalse(toPageable(page, pageSize));
        return toPaginated(result, page, pageSize);
    }

    @Override
    public TodoItemResponse getById(Long id) {
        return TodoItemResponse.from(getOrThrow(id));
    }

    // ── Commands ───────────────────────────────────────────────────────────────

    @Override
    @Transactional
    public TodoItemResponse create(CreateTodoItemRequest request) {
        TodoItem item = new TodoItem(request.title(), request.description());
        return TodoItemResponse.from(repository.save(item));
    }

    @Override
    @Transactional
    public TodoItemResponse update(Long id, UpdateTodoItemRequest request) {
        TodoItem item = getOrThrow(id);
        if (request.title() != null)       item.setTitle(request.title());
        if (request.description() != null) item.setDescription(request.description());
        if (request.isCompleted() != null) item.setCompleted(request.isCompleted());
        return TodoItemResponse.from(repository.save(item));
    }

    @Override
    @Transactional
    public TodoItemResponse markComplete(Long id) {
        TodoItem item = getOrThrow(id);
        item.setCompleted(true);
        return TodoItemResponse.from(repository.save(item));
    }

    @Override
    @Transactional
    public void delete(Long id) {
        TodoItem item = getOrThrow(id);
        repository.delete(item);
    }
}
