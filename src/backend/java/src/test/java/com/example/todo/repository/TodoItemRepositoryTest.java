package com.example.todo.repository;

import com.example.todo.entity.TodoItem;
import org.junit.jupiter.api.Test;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.autoconfigure.orm.jpa.DataJpaTest;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;

import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;

@DataJpaTest
class TodoItemRepositoryTest {

    @Autowired
    private TodoItemRepository repository;

    // ── save / findById ───────────────────────────────────────────────────────

    @Test
    void save_persistsItemAndAssignsId() {
        TodoItem saved = repository.save(new TodoItem("Buy groceries", "Milk, eggs"));

        assertThat(saved.getId()).isNotNull().isPositive();
        assertThat(saved.getTitle()).isEqualTo("Buy groceries");
        assertThat(saved.getDescription()).isEqualTo("Milk, eggs");
        assertThat(saved.isCompleted()).isFalse();
        assertThat(saved.getCreatedAt()).isNotNull();
    }

    @Test
    void findById_existingItem_returnsItem() {
        TodoItem saved = repository.save(new TodoItem("Test item", null));

        Optional<TodoItem> found = repository.findById(saved.getId());

        assertThat(found).isPresent();
        assertThat(found.get().getTitle()).isEqualTo("Test item");
        assertThat(found.get().getDescription()).isNull();
    }

    @Test
    void findById_nonExistingId_returnsEmpty() {
        Optional<TodoItem> found = repository.findById(999L);

        assertThat(found).isEmpty();
    }

    // ── findAll ───────────────────────────────────────────────────────────────

    @Test
    void findAll_returnsAllSavedItems() {
        repository.save(new TodoItem("Item 1", null));
        repository.save(new TodoItem("Item 2", null));
        repository.save(new TodoItem("Item 3", null));

        Page<TodoItem> page = repository.findAll(PageRequest.of(0, 10));

        assertThat(page.getTotalElements()).isEqualTo(3);
        assertThat(page.getContent()).hasSize(3);
    }

    @Test
    void findAll_respectsPageSizeLimit() {
        repository.save(new TodoItem("Item 1", null));
        repository.save(new TodoItem("Item 2", null));
        repository.save(new TodoItem("Item 3", null));

        Page<TodoItem> page = repository.findAll(PageRequest.of(0, 2));

        assertThat(page.getContent()).hasSize(2);
        assertThat(page.getTotalElements()).isEqualTo(3);
        assertThat(page.getTotalPages()).isEqualTo(2);
    }

    // ── findByCompletedFalse ──────────────────────────────────────────────────

    @Test
    void findByCompletedFalse_returnsOnlyIncompleteItems() {
        TodoItem complete = new TodoItem("Done", null);
        complete.setCompleted(true);
        repository.save(complete);
        repository.save(new TodoItem("Pending 1", null));
        repository.save(new TodoItem("Pending 2", null));

        Page<TodoItem> result = repository.findByCompletedFalse(PageRequest.of(0, 10));

        assertThat(result.getTotalElements()).isEqualTo(2);
        assertThat(result.getContent()).allMatch(i -> !i.isCompleted());
    }

    @Test
    void findByCompletedFalse_allCompleted_returnsEmpty() {
        TodoItem item1 = new TodoItem("Done 1", null);
        item1.setCompleted(true);
        TodoItem item2 = new TodoItem("Done 2", null);
        item2.setCompleted(true);
        repository.save(item1);
        repository.save(item2);

        Page<TodoItem> result = repository.findByCompletedFalse(PageRequest.of(0, 10));

        assertThat(result.getContent()).isEmpty();
        assertThat(result.getTotalElements()).isZero();
    }

    @Test
    void findByCompletedFalse_respectsPageSizeLimit() {
        repository.save(new TodoItem("Pending 1", null));
        repository.save(new TodoItem("Pending 2", null));
        repository.save(new TodoItem("Pending 3", null));

        Page<TodoItem> result = repository.findByCompletedFalse(
                PageRequest.of(0, 2, Sort.by("id").ascending()));

        assertThat(result.getContent()).hasSize(2);
        assertThat(result.getTotalElements()).isEqualTo(3);
    }

    // ── delete ────────────────────────────────────────────────────────────────

    @Test
    void delete_removesItemFromDatabase() {
        TodoItem saved = repository.save(new TodoItem("Delete me", null));
        Long id = saved.getId();

        repository.delete(saved);

        assertThat(repository.findById(id)).isEmpty();
    }

    // ── update ────────────────────────────────────────────────────────────────

    @Test
    void save_updatedItem_persistsChanges() {
        TodoItem saved = repository.save(new TodoItem("Original title", null));
        saved.setTitle("Updated title");
        saved.setCompleted(true);

        repository.save(saved);
        TodoItem updated = repository.findById(saved.getId()).orElseThrow();

        assertThat(updated.getTitle()).isEqualTo("Updated title");
        assertThat(updated.isCompleted()).isTrue();
    }
}
