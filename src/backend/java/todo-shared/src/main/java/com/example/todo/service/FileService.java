package com.example.todo.service;

import com.example.todo.dto.FileDownloadTarget;
import com.example.todo.dto.FileResponse;
import com.example.todo.dto.PaginatedResponse;
import org.springframework.web.multipart.MultipartFile;

/**
 * Service interface - mirrors IFileService in C#.
 * Defines the contract for the business-logic layer that manages uploaded files.
 */
public interface FileService {

    PaginatedResponse<FileResponse> getAll(int page, int pageSize);

    FileResponse getById(Long id);

    FileResponse upload(MultipartFile file);

    FileDownloadTarget getDownloadTarget(Long id);

    void delete(Long id);
}
