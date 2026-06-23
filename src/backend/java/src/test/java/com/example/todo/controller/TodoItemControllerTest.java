package com.example.todo.controller;

import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import com.example.todo.service.TodoItemService;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.time.Instant;
import java.util.List;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.delete;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.patch;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.post;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.put;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(TodoItemController.class)
class TodoItemControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @Autowired
    private ObjectMapper objectMapper;

    @MockitoBean
    private TodoItemService service;

    private TodoItemResponse sampleResponse() {
        return new TodoItemResponse(1L, "Buy groceries", "Milk, eggs, bread", false, Instant.now(), Instant.now());
    }

    private PaginatedResponse<TodoItemResponse> pagedResponse(TodoItemResponse item) {
        return new PaginatedResponse<>(List.of(item), 1L, 1, 20, 1);
    }

    // ── GET /api/todo-items ───────────────────────────────────────────────────

    @Test
    void getAll_returnsOkWithPaginatedItems() throws Exception {
        when(service.getAll(1, 20)).thenReturn(pagedResponse(sampleResponse()));

        mockMvc.perform(get("/api/todo-items"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.items").isArray())
                .andExpect(jsonPath("$.items[0].id").value(1))
                .andExpect(jsonPath("$.items[0].title").value("Buy groceries"))
                .andExpect(jsonPath("$.total").value(1))
                .andExpect(jsonPath("$.page").value(1))
                .andExpect(jsonPath("$.pageSize").value(20))
                .andExpect(jsonPath("$.totalPages").value(1));
    }

    @Test
    void getAll_withCustomPagination_passesParamsToService() throws Exception {
        when(service.getAll(2, 5)).thenReturn(new PaginatedResponse<>(List.of(), 0L, 2, 5, 0));

        mockMvc.perform(get("/api/todo-items").param("page", "2").param("pageSize", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.page").value(2))
                .andExpect(jsonPath("$.pageSize").value(5));
    }

    // ── GET /api/todo-items/incomplete ────────────────────────────────────────

    @Test
    void getIncomplete_returnsOkWithItems() throws Exception {
        when(service.getIncomplete(1, 20)).thenReturn(pagedResponse(sampleResponse()));

        mockMvc.perform(get("/api/todo-items/incomplete"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.items").isArray())
                .andExpect(jsonPath("$.items[0].title").value("Buy groceries"));
    }

    // ── GET /api/todo-items/{id} ──────────────────────────────────────────────

    @Test
    void getById_existingId_returnsOkWithItem() throws Exception {
        when(service.getById(1L)).thenReturn(sampleResponse());

        mockMvc.perform(get("/api/todo-items/1"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.title").value("Buy groceries"));
    }

    @Test
    void getById_nonExistingId_returnsNotFound() throws Exception {
        when(service.getById(99L)).thenThrow(new EntityNotFoundException("Todo item 99 not found."));

        mockMvc.perform(get("/api/todo-items/99"))
                .andExpect(status().isNotFound());
    }

    // ── POST /api/todo-items ──────────────────────────────────────────────────

    @Test
    void create_validRequest_returnsCreated() throws Exception {
        CreateTodoItemRequest request = new CreateTodoItemRequest("Buy groceries", "Milk, eggs");
        when(service.create(any(CreateTodoItemRequest.class))).thenReturn(sampleResponse());

        mockMvc.perform(post("/api/todo-items")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.title").value("Buy groceries"));
    }

    @Test
    void create_blankTitle_returnsBadRequest() throws Exception {
        CreateTodoItemRequest request = new CreateTodoItemRequest("", null);

        mockMvc.perform(post("/api/todo-items")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }

    @Test
    void create_titleExceeds200Chars_returnsBadRequest() throws Exception {
        CreateTodoItemRequest request = new CreateTodoItemRequest("A".repeat(201), null);

        mockMvc.perform(post("/api/todo-items")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isBadRequest());
    }

    @Test
    void create_missingBody_returnsBadRequest() throws Exception {
        mockMvc.perform(post("/api/todo-items")
                        .contentType(MediaType.APPLICATION_JSON))
                .andExpect(status().isBadRequest());
    }

    // ── PUT /api/todo-items/{id} ──────────────────────────────────────────────

    @Test
    void update_validRequest_returnsOk() throws Exception {
        UpdateTodoItemRequest request = new UpdateTodoItemRequest("Updated title", null, null);
        when(service.update(eq(1L), any(UpdateTodoItemRequest.class))).thenReturn(sampleResponse());

        mockMvc.perform(put("/api/todo-items/1")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1));
    }

    @Test
    void update_nonExistingId_returnsNotFound() throws Exception {
        UpdateTodoItemRequest request = new UpdateTodoItemRequest("Title", null, null);
        when(service.update(eq(99L), any(UpdateTodoItemRequest.class)))
                .thenThrow(new EntityNotFoundException("Todo item 99 not found."));

        mockMvc.perform(put("/api/todo-items/99")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isNotFound());
    }

    // ── PATCH /api/todo-items/{id}/complete ───────────────────────────────────

    @Test
    void markComplete_existingId_returnsOk() throws Exception {
        when(service.markComplete(1L)).thenReturn(sampleResponse());

        mockMvc.perform(patch("/api/todo-items/1/complete"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1));
    }

    @Test
    void markComplete_nonExistingId_returnsNotFound() throws Exception {
        when(service.markComplete(99L)).thenThrow(new EntityNotFoundException("Todo item 99 not found."));

        mockMvc.perform(patch("/api/todo-items/99/complete"))
                .andExpect(status().isNotFound());
    }

    // ── DELETE /api/todo-items/{id} ───────────────────────────────────────────

    @Test
    void delete_existingId_returnsNoContent() throws Exception {
        doNothing().when(service).delete(1L);

        mockMvc.perform(delete("/api/todo-items/1"))
                .andExpect(status().isNoContent());
    }

    @Test
    void delete_nonExistingId_returnsNotFound() throws Exception {
        doThrow(new EntityNotFoundException("Todo item 99 not found.")).when(service).delete(99L);

        mockMvc.perform(delete("/api/todo-items/99"))
                .andExpect(status().isNotFound());
    }
}
