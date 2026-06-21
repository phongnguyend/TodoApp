package com.example.todo;

import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

/**
 * Application entry point — analogous to Program.cs in ASP.NET Core.
 *
 * @SpringBootApplication combines:
 *   @Configuration    — like a DI service registrations class
 *   @EnableAutoConfiguration — auto-wires Spring Boot starters
 *   @ComponentScan    — discovers @Service, @Repository, @Controller beans
 */
@SpringBootApplication
public class TodoApplication {

    public static void main(String[] args) {
        SpringApplication.run(TodoApplication.class, args);
    }
}
