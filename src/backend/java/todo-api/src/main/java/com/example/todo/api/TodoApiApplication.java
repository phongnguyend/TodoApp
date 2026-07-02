package com.example.todo.api;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * API application entry point - analogous to Program.cs in an ASP.NET Core Web API project.
 *
 * scanBasePackages includes the shared {@code com.example.todo} package so Spring
 * discovers entities, repositories, and services from the todo-shared module.
 * JPA repository and entity scanning is configured in {@link com.example.todo.api.config.JpaConfig}
 * to keep it out of @WebMvcTest slice contexts.
 */
@SpringBootApplication(scanBasePackages = "com.example.todo")
public class TodoApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(TodoApiApplication.class, args);
    }
}
