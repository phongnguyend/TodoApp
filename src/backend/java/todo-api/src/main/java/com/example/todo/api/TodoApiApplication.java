package com.example.todo.api;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;
import org.springframework.boot.autoconfigure.domain.EntityScan;
import org.springframework.data.jpa.repository.config.EnableJpaRepositories;

/**
 * API application entry point — analogous to Program.cs in an ASP.NET Core Web API project.
 *
 * scanBasePackages includes the shared {@code com.example.todo} package so Spring
 * discovers entities, repositories, and services from the todo-shared module.
 * @EnableJpaRepositories and @EntityScan are required because @AutoConfigurationPackage
 * defaults to this class's package and would otherwise miss the shared module.
 */
@SpringBootApplication(scanBasePackages = "com.example.todo")
@EnableJpaRepositories(basePackages = "com.example.todo.repository")
@EntityScan(basePackages = "com.example.todo.entity")
public class TodoApiApplication {

    public static void main(String[] args) {
        SpringApplication.run(TodoApiApplication.class, args);
    }
}
