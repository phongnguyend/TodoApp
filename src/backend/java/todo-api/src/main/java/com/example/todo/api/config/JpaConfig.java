package com.example.todo.api.config;

import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.context.annotation.Configuration;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

/**
 * JPA configuration — separates repository/entity scanning from the main application class
 * so that @WebMvcTest slice contexts (which have no entityManagerFactory) are not affected.
 *
 * @EnableJpaRepositories and @EntityScan are required because @AutoConfigurationPackage
 * defaults to the application class's package (com.example.todo.api) and would otherwise
 * miss repositories and entities from the todo-shared module (com.example.todo).
 */
@Configuration
@EnableJpaRepositories(basePackages = "com.example.todo.repository")
@EntityScan(basePackages = "com.example.todo.entity")
public class JpaConfig {
}
