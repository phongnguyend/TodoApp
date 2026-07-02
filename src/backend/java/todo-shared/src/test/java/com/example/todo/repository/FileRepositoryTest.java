package com.example.todo.repository;

import com.example.todo.entity.FileEntity;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
class FileRepositoryTest {

    @Autowired
    private FileRepository repository;

    private static FileEntity newFile(String name) {
        return new FileEntity(name, "txt", 123L, "text/plain", "/tmp/uploads/" + name);
    }

    // ── save / findById ───────────────────────────────────────────────────────

    @Test
    void save_persistsFileAndAssignsId() {
        FileEntity saved = repository.save(newFile("report.txt"));

        assertThat(saved.getId()).isNotNull().isPositive();
        assertThat(saved.getName()).isEqualTo("report.txt");
        assertThat(saved.getExtension()).isEqualTo("txt");
        assertThat(saved.getSize()).isEqualTo(123L);
        assertThat(saved.getContentType()).isEqualTo("text/plain");
        assertThat(saved.getCreatedAt()).isNotNull();
    }

    @Test
    void findById_existingFile_returnsFile() {
        FileEntity saved = repository.save(newFile("notes.txt"));

        Optional<FileEntity> found = repository.findById(saved.getId());

        assertThat(found).isPresent();
        assertThat(found.get().getName()).isEqualTo("notes.txt");
    }

    @Test
    void findById_nonExistingId_returnsEmpty() {
        Optional<FileEntity> found = repository.findById(999L);

        assertThat(found).isEmpty();
    }

    // ── findAll ───────────────────────────────────────────────────────────────

    @Test
    void findAll_returnsAllSavedFiles() {
        repository.save(newFile("a.txt"));
        repository.save(newFile("b.txt"));
        repository.save(newFile("c.txt"));

        Page<FileEntity> page = repository.findAll(PageRequest.of(0, 10));

        assertThat(page.getTotalElements()).isEqualTo(3);
        assertThat(page.getContent()).hasSize(3);
    }

    @Test
    void findAll_respectsPageSizeLimit() {
        repository.save(newFile("a.txt"));
        repository.save(newFile("b.txt"));
        repository.save(newFile("c.txt"));

        Page<FileEntity> page = repository.findAll(PageRequest.of(0, 2));

        assertThat(page.getContent()).hasSize(2);
        assertThat(page.getTotalElements()).isEqualTo(3);
        assertThat(page.getTotalPages()).isEqualTo(2);
    }

    // ── delete ────────────────────────────────────────────────────────────────

    @Test
    void delete_removesFileFromDatabase() {
        FileEntity saved = repository.save(newFile("delete-me.txt"));
        Long id = saved.getId();

        repository.delete(saved);

        assertThat(repository.findById(id)).isEmpty();
    }
}
