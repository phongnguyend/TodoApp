package com.example.todo.service;

import com.example.todo.dto.FileDownloadTarget;
import com.example.todo.dto.FileResponse;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.entity.FileEntity;
import com.example.todo.exception.PayloadTooLargeException;
import com.example.todo.repository.FileRepository;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.mock.web.MockMultipartFile;

import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class FileServiceImplTest {

    @Mock
    private FileRepository repository;

    @TempDir
    private Path tempDir;

    private FileServiceImpl service;
    private FileEntity file;

    @BeforeEach
    void setUp() {
        service = new FileServiceImpl(repository, tempDir.toString(), 1024L);
        file = new FileEntity("report.txt", "txt", 11L, "text/plain", tempDir.resolve("report.txt").toString());
        file.setId(1L);
    }

    // ── getAll ────────────────────────────────────────────────────────────────

    @Test
    void getAll_returnsPaginatedResponse() {
        Page<FileEntity> page = new PageImpl<>(List.of(file), PageRequest.of(0, 20), 1);
        when(repository.findAll(any(Pageable.class))).thenReturn(page);

        PaginatedResponse<FileResponse> result = service.getAll(1, 20);

        assertThat(result.items()).hasSize(1);
        assertThat(result.total()).isEqualTo(1L);
        assertThat(result.page()).isEqualTo(1);
        assertThat(result.pageSize()).isEqualTo(20);
        assertThat(result.totalPages()).isEqualTo(1);
    }

    @Test
    void getAll_emptyRepository_returnsEmptyItems() {
        Page<FileEntity> empty = new PageImpl<>(List.of(), PageRequest.of(0, 20), 0);
        when(repository.findAll(any(Pageable.class))).thenReturn(empty);

        PaginatedResponse<FileResponse> result = service.getAll(1, 20);

        assertThat(result.items()).isEmpty();
        assertThat(result.total()).isZero();
    }

    // ── getById ───────────────────────────────────────────────────────────────

    @Test
    void getById_existingId_returnsResponse() {
        when(repository.findById(1L)).thenReturn(Optional.of(file));

        FileResponse result = service.getById(1L);

        assertThat(result.id()).isEqualTo(1L);
        assertThat(result.name()).isEqualTo("report.txt");
        assertThat(result.extension()).isEqualTo("txt");
    }

    @Test
    void getById_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.getById(99L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }

    // ── upload ────────────────────────────────────────────────────────────────

    @Test
    void upload_validFile_savesContentAndReturnsResponse() {
        MockMultipartFile multipartFile = new MockMultipartFile(
                "file", "document.pdf", "application/pdf", "hello world".getBytes());
        when(repository.save(any(FileEntity.class))).thenAnswer(inv -> {
            FileEntity saved = inv.getArgument(0);
            saved.setId(2L);
            return saved;
        });

        FileResponse result = service.upload(multipartFile);

        assertThat(result.id()).isEqualTo(2L);
        assertThat(result.name()).isEqualTo("document.pdf");
        assertThat(result.extension()).isEqualTo("pdf");
        assertThat(result.size()).isEqualTo(11L);
        assertThat(result.contentType()).isEqualTo("application/pdf");
        verify(repository).save(any(FileEntity.class));
    }

    @Test
    void upload_pathTraversalInFileName_stripsDirectoryComponents() {
        MockMultipartFile multipartFile = new MockMultipartFile(
                "file", "../../etc/passwd", "text/plain", "data".getBytes());
        when(repository.save(any(FileEntity.class))).thenAnswer(inv -> inv.getArgument(0));

        FileResponse result = service.upload(multipartFile);

        assertThat(result.name()).isEqualTo("passwd");
    }

    @Test
    void upload_exceedsMaxSize_throwsPayloadTooLargeException() {
        MockMultipartFile multipartFile = new MockMultipartFile(
                "file", "big.bin", "application/octet-stream", new byte[2048]);

        assertThatThrownBy(() -> service.upload(multipartFile))
                .isInstanceOf(PayloadTooLargeException.class);
    }

    // ── getDownloadTarget ─────────────────────────────────────────────────────

    @Test
    void getDownloadTarget_existingFileOnDisk_returnsTarget() throws Exception {
        Path location = tempDir.resolve("report.txt");
        Files.writeString(location, "content");
        file.setLocation(location.toString());
        when(repository.findById(1L)).thenReturn(Optional.of(file));

        FileDownloadTarget target = service.getDownloadTarget(1L);

        assertThat(target.path()).isEqualTo(location.toString());
        assertThat(target.name()).isEqualTo("report.txt");
        assertThat(target.contentType()).isEqualTo("text/plain");
    }

    @Test
    void getDownloadTarget_missingOnDisk_throwsEntityNotFoundException() {
        when(repository.findById(1L)).thenReturn(Optional.of(file));

        assertThatThrownBy(() -> service.getDownloadTarget(1L))
                .isInstanceOf(EntityNotFoundException.class);
    }

    @Test
    void getDownloadTarget_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.getDownloadTarget(99L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }

    // ── delete ────────────────────────────────────────────────────────────────

    @Test
    void delete_existingId_deletesFileAndRemovesContent() throws Exception {
        Path location = tempDir.resolve("report.txt");
        Files.writeString(location, "content");
        file.setLocation(location.toString());
        when(repository.findById(1L)).thenReturn(Optional.of(file));

        service.delete(1L);

        verify(repository).delete(file);
        assertThat(Files.exists(location)).isFalse();
    }

    @Test
    void delete_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.delete(99L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }
}
