package com.example.todo.worker;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Worker application entry point — analogous to Program.cs in an ASP.NET Core Worker Service project.
 *
 * scanBasePackages includes the shared {@code com.example.todo} package so Spring
 * discovers entities, repositories, and services from the todo-shared module.
 * The web server is disabled via {@code spring.main.web-application-type=none}
 * in {@code application.yml}.
 */
@SpringBootApplication(scanBasePackages = "com.example.todo")
public class TodoWorkerApplication {

    public static void main(String[] args) {
        SpringApplication.run(TodoWorkerApplication.class, args);
    }
}
