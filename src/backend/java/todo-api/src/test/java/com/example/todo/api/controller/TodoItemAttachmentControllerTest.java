package com.example.todo.api.controller;

import com.example.todo.dto.SaveTodoItemAttachmentRequest;
import com.example.todo.dto.TodoItemAttachmentResponse;
import com.example.todo.service.TodoItemAttachmentService;
import com.fasterxml.jackson.databind.ObjectMapper;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.boot.test.autoconfigure.web.servlet.AutoConfigureMockMvc;
import org.springframework.http.MediaType;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import java.time.Instant;
import java.util.List;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.*;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(TodoItemAttachmentController.class)
@AutoConfigureMockMvc(addFilters = false)
class TodoItemAttachmentControllerTest {

    @Autowired private MockMvc mockMvc;
    @Autowired private ObjectMapper objectMapper;
    @MockitoBean private TodoItemAttachmentService service;

    @Test
    void supportsAttachmentCrudEndpoints() throws Exception {
        TodoItemAttachmentResponse response = new TodoItemAttachmentResponse(
                3L, 10L, 5L, Instant.parse("2026-07-18T00:00:00Z"), null);
        SaveTodoItemAttachmentRequest request = new SaveTodoItemAttachmentRequest(5L);
        when(service.getAll(10L)).thenReturn(List.of(response));
        when(service.getById(10L, 3L)).thenReturn(response);
        when(service.create(any(), any())).thenReturn(response);
        when(service.update(any(), any(), any())).thenReturn(response);
        doNothing().when(service).delete(10L, 3L);

        mockMvc.perform(get("/api/todo-items/10/attachments"))
                .andExpect(status().isOk()).andExpect(jsonPath("$[0].fileId").value(5));
        mockMvc.perform(get("/api/todo-items/10/attachments/3"))
                .andExpect(status().isOk()).andExpect(jsonPath("$.id").value(3));
        mockMvc.perform(post("/api/todo-items/10/attachments")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isCreated());
        mockMvc.perform(put("/api/todo-items/10/attachments/3")
                        .contentType(MediaType.APPLICATION_JSON)
                        .content(objectMapper.writeValueAsString(request)))
                .andExpect(status().isOk());
        mockMvc.perform(delete("/api/todo-items/10/attachments/3"))
                .andExpect(status().isNoContent());
    }

    @Test
    void rejectsMissingFileId() throws Exception {
        mockMvc.perform(post("/api/todo-items/10/attachments")
                        .contentType(MediaType.APPLICATION_JSON).content("{}"))
                .andExpect(status().isBadRequest());
    }

    @Test
    void mapsMissingTodoItemToNotFound() throws Exception {
        when(service.getAll(99L)).thenThrow(new EntityNotFoundException("Todo item 99 not found."));

        mockMvc.perform(get("/api/todo-items/99/attachments"))
                .andExpect(status().isNotFound());
    }
}
