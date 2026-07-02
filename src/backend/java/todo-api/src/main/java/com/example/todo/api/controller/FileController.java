package com.example.todo.api.controller;

import com.example.todo.dto.FileDownloadTarget;
import com.example.todo.dto.FileResponse;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.service.FileService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.Parameter;
import io.swagger.v3.oas.annotations.media.Content;
import io.swagger.v3.oas.annotations.media.Schema;
import io.swagger.v3.oas.annotations.parameters.RequestBody;
import io.swagger.v3.oas.annotations.responses.ApiResponse;
import io.swagger.v3.oas.annotations.tags.Tag;
import lombok.RequiredArgsConstructor;
import org.springframework.core.io.FileSystemResource;
import org.springframework.core.io.Resource;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.MediaType;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.*;
import org.springframework.web.multipart.MultipartFile;

/**
 * REST controller for uploaded files - analogous to a FilesController : ControllerBase in ASP.NET Core.
 *
 * @RestController   = [ApiController] + [Route]
 * @RequestMapping   = [Route("api/files")]
 * @RequiredArgsConstructor = constructor injection (mirrors DI in ASP.NET Core controllers)
 */
@RestController
@RequestMapping("/api/files")
@RequiredArgsConstructor
@Tag(name = "Files")
public class FileController {

    private final FileService service;

    @GetMapping
    @Operation(summary = "Get all uploaded files (paginated)")
    public PaginatedResponse<FileResponse> getAll(
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize
    ) {
        return service.getAll(page, pageSize);
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get a file's metadata by ID")
    @ApiResponse(responseCode = "404", description = "File not found")
    public FileResponse getById(@PathVariable Long id) {
        return service.getById(id);
    }

    @GetMapping("/{id}/download")
    @Operation(summary = "Download a file's content")
    @ApiResponse(responseCode = "404", description = "File or its content not found")
    public ResponseEntity<Resource> download(@PathVariable Long id) {
        FileDownloadTarget target = service.getDownloadTarget(id);
        Resource resource = new FileSystemResource(target.path());

        MediaType mediaType;
        try {
            mediaType = MediaType.parseMediaType(target.contentType());
        } catch (IllegalArgumentException e) {
            mediaType = MediaType.APPLICATION_OCTET_STREAM;
        }

        return ResponseEntity.ok()
                .contentType(mediaType)
                .header(HttpHeaders.CONTENT_DISPOSITION, "attachment; filename=\"" + target.name() + "\"")
                .body(resource);
    }

    @PostMapping(consumes = MediaType.MULTIPART_FORM_DATA_VALUE)
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Upload a file")
    @RequestBody(content = @Content(mediaType = MediaType.MULTIPART_FORM_DATA_VALUE))
    @ApiResponse(responseCode = "413", description = "File exceeds the maximum allowed size")
    public FileResponse upload(
            @Parameter(schema = @Schema(type = "string", format = "binary"))
            @RequestParam("file") MultipartFile file
    ) {
        return service.upload(file);
    }

    @DeleteMapping("/{id}")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @Operation(summary = "Delete a file")
    @ApiResponse(responseCode = "404", description = "File not found")
    public void delete(@PathVariable Long id) {
        service.delete(id);
    }
}
