package com.example.todo.repository;

import com.example.todo.entity.EmailLog;
import org.springframework.data.jpa.repository.JpaRepository;
import org.springframework.stereotype.Repository;

/**
 * EmailLogRepository - CRUD access for the email_logs table.
 * Spring Data JPA generates the implementation at runtime.
 */
@Repository
public interface EmailLogRepository extends JpaRepository<EmailLog, Long> {
}
