package com.example.todo.service;

import com.example.todo.dto.SaveTodoItemAttachmentRequest;
import com.example.todo.entity.TodoItemAttachment;
import com.example.todo.repository.FileRepository;
import com.example.todo.repository.TodoItemAttachmentRepository;
import com.example.todo.repository.TodoItemRepository;
import jakarta.persistence.EntityNotFoundException;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.never;
import static org.mockito.Mockito.verify;
import static org.mockito.Mockito.when;

@ExtendWith(MockitoExtension.class)
class TodoItemAttachmentServiceImplTest {

    @Mock private TodoItemAttachmentRepository repository;
    @Mock private TodoItemRepository todoItems;
    @Mock private FileRepository files;
    private TodoItemAttachmentServiceImpl service;

    @BeforeEach
    void setUp() {
        service = new TodoItemAttachmentServiceImpl(repository, todoItems, files);
    }

    @Test
    void getAll_requiresExistingTodoItem() {
        when(todoItems.existsById(10L)).thenReturn(false);

        assertThatThrownBy(() -> service.getAll(10L))
                .isInstanceOf(EntityNotFoundException.class)
                .hasMessageContaining("10");
        verify(repository, never()).findByTodoItemIdOrderByCreatedAtAsc(any());
    }

    @Test
    void create_savesAttachmentWhenReferencesExist() {
        when(todoItems.existsById(10L)).thenReturn(true);
        when(files.existsById(5L)).thenReturn(true);
        when(repository.findByTodoItemIdAndFileId(10L, 5L)).thenReturn(Optional.empty());
        when(repository.save(any())).thenAnswer(invocation -> {
            TodoItemAttachment value = invocation.getArgument(0);
            value.setId(3L);
            return value;
        });

        var result = service.create(10L, new SaveTodoItemAttachmentRequest(5L));

        assertThat(result.id()).isEqualTo(3L);
        assertThat(result.todoItemId()).isEqualTo(10L);
        assertThat(result.fileId()).isEqualTo(5L);
    }

    @Test
    void create_returnsExistingDuplicateWithoutSaving() {
        TodoItemAttachment existing = attachment(3L, 10L, 5L);
        when(todoItems.existsById(10L)).thenReturn(true);
        when(files.existsById(5L)).thenReturn(true);
        when(repository.findByTodoItemIdAndFileId(10L, 5L)).thenReturn(Optional.of(existing));

        assertThat(service.create(10L, new SaveTodoItemAttachmentRequest(5L)).id()).isEqualTo(3L);
        verify(repository, never()).save(any());
    }

    @Test
    void update_changesFileForAttachmentOwnedByTodoItem() {
        TodoItemAttachment existing = attachment(3L, 10L, 5L);
        when(todoItems.existsById(10L)).thenReturn(true);
        when(files.existsById(6L)).thenReturn(true);
        when(repository.findByIdAndTodoItemId(3L, 10L)).thenReturn(Optional.of(existing));
        when(repository.findByTodoItemIdAndFileId(10L, 6L)).thenReturn(Optional.empty());
        when(repository.save(existing)).thenReturn(existing);

        assertThat(service.update(10L, 3L, new SaveTodoItemAttachmentRequest(6L)).fileId()).isEqualTo(6L);
    }

    @Test
    void delete_removesOnlyAttachmentOwnedByTodoItem() {
        TodoItemAttachment existing = attachment(3L, 10L, 5L);
        when(todoItems.existsById(10L)).thenReturn(true);
        when(repository.findByIdAndTodoItemId(3L, 10L)).thenReturn(Optional.of(existing));

        service.delete(10L, 3L);

        verify(repository).delete(existing);
    }

    private static TodoItemAttachment attachment(Long id, Long todoItemId, Long fileId) {
        TodoItemAttachment value = new TodoItemAttachment(todoItemId, fileId);
        value.setId(id);
        return value;
    }
}
