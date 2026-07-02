package com.example.todo.api.controller;

import com.example.todo.dto.FileDownloadTarget;
import com.example.todo.dto.FileResponse;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.exception.PayloadTooLargeException;
import com.example.todo.service.FileService;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.web.servlet.WebMvcTest;
import org.springframework.mock.web.MockMultipartFile;
import org.springframework.test.context.bean.override.mockito.MockitoBean;
import org.springframework.test.web.servlet.MockMvc;

import org.junit.jupiter.api.io.TempDir;

import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Instant;
import java.util.List;

import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.doNothing;
import static org.mockito.Mockito.doThrow;
import static org.mockito.Mockito.when;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.delete;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.get;
import static org.springframework.test.web.servlet.request.MockMvcRequestBuilders.multipart;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.header;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.jsonPath;
import static org.springframework.test.web.servlet.result.MockMvcResultMatchers.status;

@WebMvcTest(FileController.class)
class FileControllerTest {

    @Autowired
    private MockMvc mockMvc;

    @MockitoBean
    private FileService service;

    private FileResponse sampleResponse() {
        return new FileResponse(1L, "report.txt", "txt", 11L, "text/plain", Instant.now(), Instant.now());
    }

    private PaginatedResponse<FileResponse> pagedResponse(FileResponse item) {
        return new PaginatedResponse<>(List.of(item), 1L, 1, 20, 1);
    }

    // ── GET /api/files ────────────────────────────────────────────────────────

    @Test
    void getAll_returnsOkWithPaginatedFiles() throws Exception {
        when(service.getAll(1, 20)).thenReturn(pagedResponse(sampleResponse()));

        mockMvc.perform(get("/api/files"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.items").isArray())
                .andExpect(jsonPath("$.items[0].id").value(1))
                .andExpect(jsonPath("$.items[0].name").value("report.txt"))
                .andExpect(jsonPath("$.total").value(1))
                .andExpect(jsonPath("$.page").value(1))
                .andExpect(jsonPath("$.pageSize").value(20));
    }

    @Test
    void getAll_withCustomPagination_passesParamsToService() throws Exception {
        when(service.getAll(2, 5)).thenReturn(new PaginatedResponse<>(List.of(), 0L, 2, 5, 0));

        mockMvc.perform(get("/api/files").param("page", "2").param("pageSize", "5"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.page").value(2))
                .andExpect(jsonPath("$.pageSize").value(5));
    }

    // ── GET /api/files/{id} ───────────────────────────────────────────────────

    @Test
    void getById_existingId_returnsOkWithFile() throws Exception {
        when(service.getById(1L)).thenReturn(sampleResponse());

        mockMvc.perform(get("/api/files/1"))
                .andExpect(status().isOk())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.name").value("report.txt"));
    }

    @Test
    void getById_nonExistingId_returnsNotFound() throws Exception {
        when(service.getById(99L)).thenThrow(new EntityNotFoundException("File 99 not found."));

        mockMvc.perform(get("/api/files/99"))
                .andExpect(status().isNotFound());
    }

    // ── GET /api/files/{id}/download ──────────────────────────────────────────

    @Test
    void download_existingFile_returnsOkWithContentDisposition(@TempDir Path tempDir) throws Exception {
        Path location = tempDir.resolve("report.txt");
        Files.writeString(location, "hello world");
        when(service.getDownloadTarget(1L)).thenReturn(new FileDownloadTarget(location.toString(), "report.txt", "text/plain"));

        mockMvc.perform(get("/api/files/1/download"))
                .andExpect(status().isOk())
                .andExpect(header().string("Content-Disposition", "attachment; filename=\"report.txt\""));
    }

    @Test
    void download_nonExistingId_returnsNotFound() throws Exception {
        when(service.getDownloadTarget(99L)).thenThrow(new EntityNotFoundException("File 99 not found."));

        mockMvc.perform(get("/api/files/99/download"))
                .andExpect(status().isNotFound());
    }

    // ── POST /api/files ───────────────────────────────────────────────────────

    @Test
    void upload_validFile_returnsCreated() throws Exception {
        MockMultipartFile multipartFile = new MockMultipartFile(
                "file", "report.txt", "text/plain", "hello world".getBytes());
        when(service.upload(any())).thenReturn(sampleResponse());

        mockMvc.perform(multipart("/api/files").file(multipartFile))
                .andExpect(status().isCreated())
                .andExpect(jsonPath("$.id").value(1))
                .andExpect(jsonPath("$.name").value("report.txt"));
    }

    @Test
    void upload_missingFile_returnsBadRequest() throws Exception {
        mockMvc.perform(multipart("/api/files"))
                .andExpect(status().isBadRequest());
    }

    @Test
    void upload_exceedsMaxSize_returnsPayloadTooLarge() throws Exception {
        MockMultipartFile multipartFile = new MockMultipartFile(
                "file", "big.bin", "application/octet-stream", new byte[10]);
        when(service.upload(any())).thenThrow(new PayloadTooLargeException("File exceeds the maximum allowed size of 1024 bytes."));

        mockMvc.perform(multipart("/api/files").file(multipartFile))
                .andExpect(status().isPayloadTooLarge());
    }

    // ── DELETE /api/files/{id} ────────────────────────────────────────────────

    @Test
    void delete_existingId_returnsNoContent() throws Exception {
        doNothing().when(service).delete(1L);

        mockMvc.perform(delete("/api/files/1"))
                .andExpect(status().isNoContent());
    }

    @Test
    void delete_nonExistingId_returnsNotFound() throws Exception {
        doThrow(new EntityNotFoundException("File 99 not found.")).when(service).delete(99L);

        mockMvc.perform(delete("/api/files/99"))
                .andExpect(status().isNotFound());
    }
}
