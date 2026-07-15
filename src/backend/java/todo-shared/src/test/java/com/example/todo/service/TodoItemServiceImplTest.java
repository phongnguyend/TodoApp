package com.example.todo.service;

import com.example.todo.dto.CreateTodoItemRequest;
import com.example.todo.dto.ImportResult;
import com.example.todo.dto.PaginatedResponse;
import com.example.todo.dto.TodoItemResponse;
import com.example.todo.dto.UpdateTodoItemRequest;
import com.example.todo.entity.TodoItem;
import com.example.todo.repository.TodoItemRepository;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.InjectMocks;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Pageable;
import org.springframework.mock.web.MockMultipartFile;

import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class TodoItemServiceImplTest {

    @Mock
    private TodoItemRepository repository;

    @InjectMocks
    private TodoItemServiceImpl service;

    private TodoItem item;

    @BeforeEach
    void setUp() {
        item = new TodoItem("Buy groceries", "Milk, eggs, bread");
        item.setId(1L);
    }

    // ── getAll ────────────────────────────────────────────────────────────────

    @Test
    void getAll_returnsPaginatedResponse() {
        Page<TodoItem> page = new PageImpl<>(List.of(item), PageRequest.of(0, 20), 1);
        when(repository.findAll(any(Pageable.class))).thenReturn(page);

        PaginatedResponse<TodoItemResponse> result = service.getAll(1, 20);

        assertThat(result.items()).hasSize(1);
        assertThat(result.total()).isEqualTo(1L);
        assertThat(result.page()).isEqualTo(1);
        assertThat(result.pageSize()).isEqualTo(20);
        assertThat(result.totalPages()).isEqualTo(1);
    }

    @Test
    void getAll_emptyRepository_returnsEmptyItems() {
        Page<TodoItem> empty = new PageImpl<>(List.of(), PageRequest.of(0, 20), 0);
        when(repository.findAll(any(Pageable.class))).thenReturn(empty);

        PaginatedResponse<TodoItemResponse> result = service.getAll(1, 20);

        assertThat(result.items()).isEmpty();
        assertThat(result.total()).isZero();
    }

    // ── getIncomplete ─────────────────────────────────────────────────────────

    @Test
    void getIncomplete_returnsPaginatedResponse() {
        Page<TodoItem> page = new PageImpl<>(List.of(item), PageRequest.of(0, 10), 1);
        when(repository.findByCompletedFalse(any(Pageable.class))).thenReturn(page);

        PaginatedResponse<TodoItemResponse> result = service.getIncomplete(1, 10);

        assertThat(result.items()).hasSize(1);
        assertThat(result.total()).isEqualTo(1L);
        assertThat(result.page()).isEqualTo(1);
        assertThat(result.pageSize()).isEqualTo(10);
    }

    // ── getById ───────────────────────────────────────────────────────────────

    @Test
    void getById_existingId_returnsResponse() {
        when(repository.findById(1L)).thenReturn(Optional.of(item));

        TodoItemResponse result = service.getById(1L);

        assertThat(result.id()).isEqualTo(1L);
        assertThat(result.title()).isEqualTo("Buy groceries");
        assertThat(result.description()).isEqualTo("Milk, eggs, bread");
    }

    @Test
    void getById_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.getById(99L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }

    // ── create ────────────────────────────────────────────────────────────────

    @Test
    void create_validRequest_savesAndReturnsResponse() {
        CreateTodoItemRequest request = new CreateTodoItemRequest("Buy groceries", "Milk, eggs");
        when(repository.save(any(TodoItem.class))).thenReturn(item);

        TodoItemResponse result = service.create(request);

        assertThat(result.id()).isEqualTo(1L);
        assertThat(result.title()).isEqualTo("Buy groceries");
        verify(repository).save(any(TodoItem.class));
    }

    @Test
    void create_nullDescription_savesSuccessfully() {
        CreateTodoItemRequest request = new CreateTodoItemRequest("Title only", null);
        TodoItem saved = new TodoItem("Title only", null);
        saved.setId(2L);
        when(repository.save(any(TodoItem.class))).thenReturn(saved);

        TodoItemResponse result = service.create(request);

        assertThat(result.title()).isEqualTo("Title only");
        assertThat(result.description()).isNull();
    }

    // ── update ────────────────────────────────────────────────────────────────

    @Test
    void update_allFields_updatesAndReturnsResponse() {
        when(repository.findById(1L)).thenReturn(Optional.of(item));
        when(repository.save(any(TodoItem.class))).thenAnswer(inv -> inv.getArgument(0));

        UpdateTodoItemRequest request = new UpdateTodoItemRequest("New title", "New description", true);
        TodoItemResponse result = service.update(1L, request);

        assertThat(result.title()).isEqualTo("New title");
        assertThat(result.description()).isEqualTo("New description");
        assertThat(result.isCompleted()).isTrue();
    }

    @Test
    void update_nullFields_keepsExistingValues() {
        when(repository.findById(1L)).thenReturn(Optional.of(item));
        when(repository.save(any(TodoItem.class))).thenAnswer(inv -> inv.getArgument(0));

        UpdateTodoItemRequest request = new UpdateTodoItemRequest(null, null, null);
        TodoItemResponse result = service.update(1L, request);

        assertThat(result.title()).isEqualTo("Buy groceries");
        assertThat(result.description()).isEqualTo("Milk, eggs, bread");
        assertThat(result.isCompleted()).isFalse();
    }

    @Test
    void update_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.update(99L, new UpdateTodoItemRequest("T", null, null)))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }

    // ── markComplete ──────────────────────────────────────────────────────────

    @Test
    void markComplete_existingId_setsCompletedTrue() {
        when(repository.findById(1L)).thenReturn(Optional.of(item));
        when(repository.save(any(TodoItem.class))).thenAnswer(inv -> inv.getArgument(0));

        TodoItemResponse result = service.markComplete(1L);

        assertThat(result.isCompleted()).isTrue();
        verify(repository).save(item);
    }

    @Test
    void markComplete_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.markComplete(99L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }

    // ── delete ────────────────────────────────────────────────────────────────

    @Test
    void delete_existingId_deletesItem() {
        when(repository.findById(1L)).thenReturn(Optional.of(item));

        service.delete(1L);

        verify(repository).delete(item);
    }

    @Test
    void delete_nonExistingId_throwsEntityNotFoundException() {
        when(repository.findById(99L)).thenReturn(Optional.empty());

        assertThatThrownBy(() -> service.delete(99L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("99");
    }

    // ── importCsv ──────────────────────────────────────────────────────

    @Test
    void importCsv_validRows_importsAndReturnsSummary() {
        MockMultipartFile file = new MockMultipartFile(
                "file", "todo_items.csv", "text/csv",
                "title,description,is_completed\nBuy milk,Whole milk,true\nBuy eggs,,false\n".getBytes());
        when(repository.save(any(TodoItem.class))).thenAnswer(inv -> inv.getArgument(0));

        ImportResult result = service.importCsv(file);

        assertThat(result.imported()).isEqualTo(2);
        assertThat(result.failed()).isZero();
        assertThat(result.errors()).isEmpty();
        verify(repository, org.mockito.Mockito.times(2)).save(any(TodoItem.class));
    }

    @Test
    void importCsv_blankTitle_recordsRowError() {
        MockMultipartFile file = new MockMultipartFile(
                "file", "todo_items.csv", "text/csv",
                "title,description,is_completed\n,,false\n".getBytes());

        ImportResult result = service.importCsv(file);

        assertThat(result.imported()).isZero();
        assertThat(result.failed()).isEqualTo(1);
        assertThat(result.errors()).hasSize(1);
        assertThat(result.errors().get(0).row()).isEqualTo(2);
        assertThat(result.errors().get(0).error()).isEqualTo("Title is required.");
    }

    @Test
    void importCsv_emptyFile_returnsEmptySummary() {
        MockMultipartFile file = new MockMultipartFile("file", "todo_items.csv", "text/csv", new byte[0]);

        ImportResult result = service.importCsv(file);

        assertThat(result.imported()).isZero();
        assertThat(result.failed()).isZero();
        assertThat(result.errors()).isEmpty();
    }

    // ── exportCsv ───────────────────────────────────────────────────

    @Test
    void exportCsv_returnsHeaderAndRows() {
        when(repository.findAllByOrderByCreatedAtDesc()).thenReturn(List.of(item));

        String csv = service.exportCsv();
        List<String> lines = List.of(csv.split("\r\n"));

        assertThat(lines.get(0)).isEqualTo("id,title,description,is_completed,created_at,updated_at");
        assertThat(lines.get(1)).contains("1,Buy groceries,\"Milk, eggs, bread\",false");
    }

    @Test
    void exportCsv_noItems_returnsHeaderOnly() {
        when(repository.findAllByOrderByCreatedAtDesc()).thenReturn(List.of());

        String csv = service.exportCsv();

        assertThat(csv).isEqualTo("id,title,description,is_completed,created_at,updated_at\r\n");
    }
}
