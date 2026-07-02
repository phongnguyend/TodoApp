package com.example.todo.api.controller;

import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import com.example.todo.service.TodoItemService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

/**
 * REST controller - analogous to an [ApiController] class in ASP.NET Core.
 *
 * @RestController   = [ApiController] + [Route]
 * @RequestMapping   = [Route("api/todo-items")]
 * @GetMapping etc.  = [HttpGet] / [HttpPost] etc.
 * @Valid            = automatic model validation (like the ASP.NET Core model-binding pipeline)
 * @RequiredArgsConstructor = constructor injection (mirrors DI in ASP.NET Core controllers)
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
            @RequestParam(defaultValue = "1")  int page,
            @RequestParam(defaultValue = "20") int pageSize
    ) {
        return service.getAll(page, pageSize);
    }

    @GetMapping("/incomplete")
    @Operation(summary = "Get incomplete todo items (paginated)")
    public PaginatedResponse<TodoItemResponse> getIncomplete(
            @RequestParam(defaultValue = "1")  int page,
            @RequestParam(defaultValue = "20") int pageSize
    ) {
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
        return service.create(request);
    }

    @PutMapping("/{id}")
    @Operation(summary = "Update a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public TodoItemResponse update(
            @PathVariable Long id,
            @Valid @RequestBody UpdateTodoItemRequest request
    ) {
        return service.update(id, request);
    }

    @PatchMapping("/{id}/complete")
    @Operation(summary = "Mark a todo item as complete")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public TodoItemResponse markComplete(@PathVariable Long id) {
        return service.markComplete(id);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @Operation(summary = "Delete a todo item")
    @ApiResponse(responseCode = "404", description = "Todo item not found")
    public void delete(@PathVariable Long id) {
        service.delete(id);
    }
}
