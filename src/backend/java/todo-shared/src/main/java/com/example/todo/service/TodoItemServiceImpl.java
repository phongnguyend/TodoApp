package com.example.todo.service;

import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.ImportResult;
import com.example.todo.dto.ImportRowError;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import com.example.todo.entity.TodoItem;
import com.example.todo.repository.TodoItemRepository;
import com.example.todo.util.CsvUtil;
import com.example.todo.util.ExcelUtil;
import jakarta.persistence.EntityNotFoundException;
import lombok.RequiredArgsConstructor;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.charset.StandardCharsets;
import java.util.ArrayList;
import java.util.List;
import java.util.Map;
import java.util.Set;
import java.util.stream.Collectors;
import java.util.stream.IntStream;

/**
 * Service implementation - registered as a Spring-managed bean via @Service.
 * Mirrors a service class injected through ASP.NET Core's DI container.
 *
 * @Transactional mirrors [UnitOfWork] / SaveChanges() semantics in EF Core.
 */
@Service
@RequiredArgsConstructor
@Transactional(readOnly = true)
public class TodoItemServiceImpl implements TodoItemService {

    private static final List<String> CSV_HEADER = List.of("id", "title", "description", "is_completed", "created_at",
            "updated_at");
    private static final Set<String> TRUE_VALUES = Set.of("1", "true", "yes", "y");

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
        return new PaginatedResponse<>(items, pageResult.getTotalElements(), page, pageSize,
                pageResult.getTotalPages());
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
        return create(request, null);
    }

    @Override
    @Transactional
    public TodoItemResponse create(CreateTodoItemRequest request, Long actorUserId) {
        TodoItem item = new TodoItem(request.title(), request.description());
        item.setCreatedByUserId(actorUserId);
        return TodoItemResponse.from(repository.save(item));
    }

    @Override
    @Transactional
    public TodoItemResponse update(Long id, UpdateTodoItemRequest request) {
        return update(id, request, null);
    }

    @Override
    @Transactional
    public TodoItemResponse update(Long id, UpdateTodoItemRequest request, Long actorUserId) {
        TodoItem item = getOrThrow(id);
        if (request.title() != null)
            item.setTitle(request.title());
        if (request.description() != null)
            item.setDescription(request.description());
        if (request.isCompleted() != null)
            item.setCompleted(request.isCompleted());
        item.setUpdatedByUserId(actorUserId);
        return TodoItemResponse.from(repository.save(item));
    }

    @Override
    @Transactional
    public TodoItemResponse markComplete(Long id) {
        return markComplete(id, null);
    }

    @Override
    @Transactional
    public TodoItemResponse markComplete(Long id, Long actorUserId) {
        TodoItem item = getOrThrow(id);
        item.setCompleted(true);
        item.setUpdatedByUserId(actorUserId);
        return TodoItemResponse.from(repository.save(item));
    }

    @Override
    @Transactional
    public void delete(Long id) {
        TodoItem item = getOrThrow(id);
        repository.delete(item);
    }

    // ── CSV import/export ─────────────────────────────────────────────────────

    @Override
    @Transactional
    public ImportResult importCsv(MultipartFile file) {
        return importCsv(file, null);
    }

    @Override
    @Transactional
    public ImportResult importCsv(MultipartFile file, Long actorUserId) {
        List<List<String>> rows = CsvUtil.parse(readAsUtf8(file));
        return importRows(rows, actorUserId);
    }

    @Override
    public String exportCsv() {
        List<TodoItem> items = repository.findAllByOrderByCreatedAtDesc();

        StringBuilder sb = new StringBuilder();
        sb.append(CsvUtil.toCsvRow(CSV_HEADER));
        for (TodoItem item : items) {
            sb.append(CsvUtil.toCsvRow(List.of(
                    item.getId(),
                    item.getTitle(),
                    item.getDescription() != null ? item.getDescription() : "",
                    item.isCompleted(),
                    item.getCreatedAt() != null ? item.getCreatedAt().toString() : "",
                    item.getUpdatedAt() != null ? item.getUpdatedAt().toString() : "")));
        }
        return sb.toString();
    }

    // ── Excel import/export ───────────────────────────────────────────────────

    @Override
    @Transactional
    public ImportResult importExcel(MultipartFile file) {
        return importExcel(file, null);
    }

    @Override
    @Transactional
    public ImportResult importExcel(MultipartFile file, Long actorUserId) {
        List<List<String>> rows;
        try {
            rows = ExcelUtil.parse(file.getInputStream());
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
        return importRows(rows, actorUserId);
    }

    @Override
    public byte[] exportExcel() {
        List<TodoItem> items = repository.findAllByOrderByCreatedAtDesc();

        List<List<Object>> rows = items.stream()
                .<List<Object>>map(item -> List.of(
                        item.getId(),
                        item.getTitle(),
                        item.getDescription() != null ? item.getDescription() : "",
                        item.isCompleted(),
                        item.getCreatedAt() != null ? item.getCreatedAt().toString() : "",
                        item.getUpdatedAt() != null ? item.getUpdatedAt().toString() : ""))
                .toList();

        return ExcelUtil.write(CSV_HEADER, rows);
    }

    // ── Shared import row processing (CSV & Excel) ────────────────────────────

    private ImportResult importRows(List<List<String>> rows, Long actorUserId) {
        if (rows.isEmpty()) {
            return new ImportResult(0, 0, List.of());
        }

        List<String> header = rows.get(0).stream().map(col -> col.trim().toLowerCase()).toList();
        Map<String, Integer> colIndex = IntStream.range(0, header.size())
                .boxed()
                .collect(Collectors.toMap(header::get, i -> i, (first, second) -> first));

        int imported = 0;
        List<ImportRowError> errors = new ArrayList<>();

        for (int i = 1; i < rows.size(); i++) {
            int rowNumber = i + 1; // header occupies row 1
            List<String> row = rows.get(i);

            String title = getCell(row, colIndex, "title").trim();
            if (title.isEmpty()) {
                errors.add(new ImportRowError(rowNumber, "Title is required."));
                continue;
            }

            String description = getCell(row, colIndex, "description").trim();
            boolean isCompleted = TRUE_VALUES.contains(getCell(row, colIndex, "is_completed").trim().toLowerCase());

            TodoItem item = new TodoItem(title, description.isEmpty() ? null : description);
            item.setCompleted(isCompleted);
            item.setCreatedByUserId(actorUserId);
            repository.save(item);
            imported++;
        }

        return new ImportResult(imported, errors.size(), errors);
    }

    private static String getCell(List<String> row, Map<String, Integer> colIndex, String name) {
        Integer idx = colIndex.get(name);
        if (idx == null || idx >= row.size()) {
            return "";
        }
        String value = row.get(idx);
        return value != null ? value : "";
    }

    private static String readAsUtf8(MultipartFile file) {
        try {
            String text = new String(file.getBytes(), StandardCharsets.UTF_8);
            return text.startsWith("\uFEFF") ? text.substring(1) : text;
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }
}
