package com.example.todo.repository;

import com.example.todo.entity.TodoItem;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.Pageable;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

import java.util.List;

/**
 * TodoItemRepository extends JpaRepository - analogous to inheriting from EF's DbSet<TodoItem>
 * with built-in CRUD operations. Spring Data JPA generates the implementation at runtime.
 *
 * Derived query method {@code findByCompletedFalse} mirrors a LINQ Where() clause in EF.
 */
@Repository
public interface TodoItemRepository extends JpaRepository<TodoItem, Long> {

    Page<TodoItem> findByCompletedFalse(Pageable pageable);

    /** Returns all incomplete todos ordered by creation date - used by the background worker. */
    List<TodoItem> findByCompletedFalseOrderByCreatedAtAsc();
}
