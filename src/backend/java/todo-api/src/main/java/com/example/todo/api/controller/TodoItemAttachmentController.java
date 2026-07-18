package com.example.todo.api.controller;

import com.example.todo.dto.SaveTodoItemAttachmentRequest;
import com.example.todo.dto.TodoItemAttachmentResponse;
import com.example.todo.service.TodoItemAttachmentService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

import java.util.List;

@RestController
@RequestMapping("/api/todo-items/{todoItemId}/attachments")
@RequiredArgsConstructor
@Tag(name = "Todo Item Attachments")
public class TodoItemAttachmentController {

    private final TodoItemAttachmentService service;

    @GetMapping
    @Operation(summary = "Get all attachments for a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public List<TodoItemAttachmentResponse> getAll(@PathVariable Long todoItemId) {
        return service.getAll(todoItemId);
    }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Attach a file to a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item or file not found")
    public TodoItemAttachmentResponse create(
            @PathVariable Long todoItemId,
            @Valid @RequestBody SaveTodoItemAttachmentRequest request) {
        return service.create(todoItemId, request);
    }

    @GetMapping("/{attachmentId}")
    @Operation(summary = "Get an attachment by ID")
    @ApiResponse(responseCode = "404", description = "Todo item or attachment not found")
    public TodoItemAttachmentResponse getById(
            @PathVariable Long todoItemId, @PathVariable Long attachmentId) {
        return service.getById(todoItemId, attachmentId);
    }

    @PutMapping("/{attachmentId}")
    @Operation(summary = "Update an attachment")
    @ApiResponse(responseCode = "404", description = "Todo item, file, or attachment not found")
    public TodoItemAttachmentResponse update(
            @PathVariable Long todoItemId,
            @PathVariable Long attachmentId,
            @Valid @RequestBody SaveTodoItemAttachmentRequest request) {
        return service.update(todoItemId, attachmentId, request);
    }

    @DeleteMapping("/{attachmentId}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @Operation(summary = "Remove an attachment from a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item or attachment not found")
    public void delete(@PathVariable Long todoItemId, @PathVariable Long attachmentId) {
        service.delete(todoItemId, attachmentId);
    }
}
