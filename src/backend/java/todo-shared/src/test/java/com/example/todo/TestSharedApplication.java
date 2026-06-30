package com.example.todo;

import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Minimal Spring Boot configuration class used exclusively by test slices
 * (e.g., {@code @DataJpaTest}) within the todo-shared module.
 *
 * Without a {@code @SpringBootApplication} in the module, Spring Boot's test
 * bootstrapper cannot determine the base package for entity and repository
 * scanning. Placing this class in {@code com.example.todo} causes Spring Boot
 * to scan all sub-packages, ensuring {@code @Entity} classes in
 * {@code com.example.todo.entity} and {@code @Repository} interfaces in
 * {@code com.example.todo.repository} are discovered correctly.
 */
@SpringBootApplication
class TestSharedApplication {
}
