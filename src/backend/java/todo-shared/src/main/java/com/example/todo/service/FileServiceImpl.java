package com.example.todo.service;

import com.example.todo.dto.FileDownloadTarget;
import com.example.todo.dto.FileResponse;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.entity.FileEntity;
import com.example.todo.exception.PayloadTooLargeException;
import com.example.todo.repository.FileRepository;
import jakarta.persistence.EntityNotFoundException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;
import org.springframework.util.StringUtils;
import org.springframework.web.multipart.MultipartFile;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;
import java.util.UUID;

/**
 * Service implementation - registered as a Spring-managed bean via @Service.
 * Mirrors a service class injected through ASP.NET Core's DI container.
 *
 * @Transactional mirrors [UnitOfWork] / SaveChanges() semantics in EF Core.
 */
@Service
@Transactional(readOnly = true)
public class FileServiceImpl implements FileService {

    private final FileRepository repository;
    private final Path storageDir;
    private final long maxUploadSizeBytes;

    public FileServiceImpl(
            FileRepository repository,
            @Value("${app.file.storage-path:./uploads}") String storagePath,
            @Value("${app.file.max-upload-size-bytes:10485760}") long maxUploadSizeBytes
    ) {
        this.repository = repository;
        this.storageDir = Paths.get(storagePath);
        this.maxUploadSizeBytes = maxUploadSizeBytes;
    }

    // ── Helpers ────────────────────────────────────────────────────────────────

    private FileEntity getOrThrow(Long id) {
        return repository.findById(id)
                .orElseThrow(() -> new EntityNotFoundException("File " + id + " not found."));
    }

    private static Pageable toPageable(int page, int pageSize) {
        return PageRequest.of(page - 1, pageSize, Sort.by("createdAt").descending());
    }

    private static PaginatedResponse<FileResponse> toPaginated(Page<FileEntity> pageResult, int page, int pageSize) {
        List<FileResponse> items = pageResult.getContent()
                .stream()
                .map(FileResponse::from)
                .toList();
        return new PaginatedResponse<>(items, pageResult.getTotalElements(), page, pageSize, pageResult.getTotalPages());
    }

    // ── Queries ────────────────────────────────────────────────────────────────

    @Override
    public PaginatedResponse<FileResponse> getAll(int page, int pageSize) {
        Page<FileEntity> result = repository.findAll(toPageable(page, pageSize));
        return toPaginated(result, page, pageSize);
    }

    @Override
    public FileResponse getById(Long id) {
        return FileResponse.from(getOrThrow(id));
    }

    @Override
    public FileDownloadTarget getDownloadTarget(Long id) {
        FileEntity entity = getOrThrow(id);
        if (!Files.isRegularFile(Paths.get(entity.getLocation()))) {
            throw new EntityNotFoundException("File " + id + " content not found on disk.");
        }
        String contentType = entity.getContentType() != null ? entity.getContentType() : "application/octet-stream";
        return new FileDownloadTarget(entity.getLocation(), entity.getName(), contentType);
    }

    // ── Commands ───────────────────────────────────────────────────────────────

    @Override
    @Transactional
    public FileResponse upload(MultipartFile file) {
        return upload(file, null);
    }

    @Override
    @Transactional
    public FileResponse upload(MultipartFile file, Long actorUserId) {
        if (file.getSize() > maxUploadSizeBytes) {
            throw new PayloadTooLargeException(
                    "File exceeds the maximum allowed size of " + maxUploadSizeBytes + " bytes.");
        }

        // Strip any directory components from the client-supplied name to prevent path traversal.
        String originalName = Paths.get(StringUtils.cleanPath(
                StringUtils.hasText(file.getOriginalFilename()) ? file.getOriginalFilename() : "unnamed"
        )).getFileName().toString();
        String extension = originalName.contains(".")
                ? originalName.substring(originalName.lastIndexOf('.') + 1).toLowerCase()
                : "";

        try {
            Files.createDirectories(storageDir);

            // A random prefix avoids collisions/overwrites between uploads that share a name.
            String storedName = UUID.randomUUID() + "_" + originalName;
            Path location = storageDir.resolve(storedName);
            file.transferTo(location);

            FileEntity entity = new FileEntity(originalName, extension, file.getSize(), file.getContentType(), location.toString());
            entity.setCreatedByUserId(actorUserId);
            return FileResponse.from(repository.save(entity));
        } catch (IOException e) {
            throw new UncheckedIOException("Failed to store uploaded file.", e);
        }
    }

    @Override
    @Transactional
    public void delete(Long id) {
        FileEntity entity = getOrThrow(id);
        repository.delete(entity);
        try {
            Files.deleteIfExists(Paths.get(entity.getLocation()));
        } catch (IOException e) {
            // Content already missing or inaccessible on disk - nothing left to clean up.
        }
    }
}
