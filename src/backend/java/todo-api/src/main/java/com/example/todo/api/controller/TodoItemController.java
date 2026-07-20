package com.example.todo.api.controller;

import com.example.todo.api.config.AuditActor;
import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.ImportResult;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import com.example.todo.service.TodoItemService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

/**
 * REST controller - analogous to an [ApiController] class in ASP.NET Core.
 *
 * @RestController = [ApiController] + [Route]
 * @RequestMapping = [Route("api/todo-items")]
 * @GetMapping etc. = [HttpGet] / [HttpPost] etc.
 * @Valid = automatic model validation (like the ASP.NET Core model-binding
 *        pipeline)
 * @RequiredArgsConstructor = constructor injection (mirrors DI in ASP.NET Core
 *                          controllers)
 */
@RestController
@RequestMapping("/api/todo-items")
@RequiredArgsConstructor
@Tag(name = "Todo Items")
public class TodoItemController {

    private final TodoItemService service;

    @GetMapping
    @Operation(summary = "Get all todo items (paginated)")
    public PaginatedResponse<TodoItemResponse> getAll(
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        return service.getAll(page, pageSize);
    }

    @GetMapping("/incomplete")
    @Operation(summary = "Get incomplete todo items (paginated)")
    public PaginatedResponse<TodoItemResponse> getIncomplete(
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        return service.getIncomplete(page, pageSize);
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get a todo item by ID")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public TodoItemResponse getById(@PathVariable Long id) {
        return service.getById(id);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Create a todo item")
    public TodoItemResponse create(@Valid @RequestBody CreateTodoItemRequest request) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.create(request) : service.create(request, actor);
    }

    @PutMapping("/{id}")
    @Operation(summary = "Update a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public TodoItemResponse update(
            @PathVariable Long id,
            @Valid @RequestBody UpdateTodoItemRequest request) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.update(id, request) : service.update(id, request, actor);
    }

    @PatchMapping("/{id}/complete")
    @Operation(summary = "Mark a todo item as complete")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public TodoItemResponse markComplete(@PathVariable Long id) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.markComplete(id) : service.markComplete(id, actor);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @Operation(summary = "Delete a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public void delete(@PathVariable Long id) {
        service.delete(id);
    }

    @PostMapping(value = "/import/csv", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @Operation(summary = "Import todo items from a CSV file")
    @io.swagger.v3.oas.annotations.parameters.RequestBody(content = @Content(mediaType = MediaType.MULTIPART_FORM_DATA_VALUE))
    public ImportResult importCsv(
            @Parameter(schema = @Schema(type = "string", format = "binary")) @RequestParam("file") MultipartFile file) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.importCsv(file) : service.importCsv(file, actor);
    }

    @GetMapping("/export/csv")
    @Operation(summary = "Export todo items as a CSV file")
    public ResponseEntity<String> exportCsv() {
        String content = service.exportCsv();
        return ResponseEntity.ok()
                .contentType(MediaType.parseMediaType("text/csv"))
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"todo_items.csv\"")
                .body(content);
    }

    @PostMapping(value = "/import/excel", consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @Operation(summary = "Import todo items from an Excel file")
    @io.swagger.v3.oas.annotations.parameters.RequestBody(content = @Content(mediaType = MediaType.MULTIPART_FORM_DATA_VALUE))
    public ImportResult importExcel(
            @Parameter(schema = @Schema(type = "string", format = "binary")) @RequestParam("file") MultipartFile file) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.importExcel(file) : service.importExcel(file, actor);
    }

    @GetMapping("/export/excel")
    @Operation(summary = "Export todo items as an Excel file")
    public ResponseEntity<byte[]> exportExcel() {
        byte[] content = service.exportExcel();
        return ResponseEntity.ok()
                .contentType(
                        MediaType.parseMediaType("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"))
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"todo_items.xlsx\"")
                .body(content);
    }
}
