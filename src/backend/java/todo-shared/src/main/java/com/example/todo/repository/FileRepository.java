package com.example.todo.repository;

import com.example.todo.entity.FileEntity;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * FileRepository extends JpaRepository - analogous to inheriting from EF's DbSet<File>
 * with built-in CRUD operations. Spring Data JPA generates the implementation at runtime.
 */
@Repository
public interface FileRepository extends JpaRepository<FileEntity, Long> {
}
