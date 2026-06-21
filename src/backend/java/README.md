# Todo API — Java / Spring Boot

A RESTful API for managing todo items built with **Spring Boot**, **Spring Data JPA**, **Hibernate**, and **Flyway** — the Java equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF | Java equivalent |
|---|---|
| ASP.NET Core | **Spring Boot** |
| Entity Framework Core | **Spring Data JPA** + Hibernate |
| EF Migrations | **Flyway** (`V1__*.sql` scripts) |
| `DbContext` | `JpaRepository` / `EntityManager` |
| `[ApiController]` / Controllers | `@RestController` |
| Services (`IService` / `Service`) | `@Service` interface + `@Service` impl |
| DTOs / Data Annotations | Java records + **Jakarta Bean Validation** |
| `appsettings.json` / `IConfiguration` | `application.yml` + Spring Environment |
| Dependency Injection | Spring DI (`@Autowired` / constructor injection) |
| Swagger / OpenAPI | **SpringDoc OpenAPI** at `/swagger` |
| `ProblemDetails` middleware | `@RestControllerAdvice` + `ProblemDetail` |
| `Program.cs` / `Startup.cs` | `TodoApplication.java` + auto-configuration |

## Project structure

```
src/backend/java/
├── pom.xml                                        # Maven build (like .csproj)
└── src/
    ├── main/
    │   ├── java/com/example/todo/
    │   │   ├── TodoApplication.java               # Entry point (@SpringBootApplication)
    │   │   ├── config/
    │   │   │   └── OpenApiConfig.java             # Swagger config
    │   │   ├── controller/
    │   │   │   └── TodoItemController.java        # @RestController
    │   │   ├── dto/
    │   │   │   ├── CreateTodoItemRequest.java     # record + @Valid
    │   │   │   ├── UpdateTodoItemRequest.java     # record + @Valid
    │   │   │   ├── TodoItemResponse.java          # record with factory method
    │   │   │   └── PaginatedResponse.java         # generic record
    │   │   ├── entity/
    │   │   │   └── TodoItem.java                  # @Entity (JPA model)
    │   │   ├── exception/
    │   │   │   └── GlobalExceptionHandler.java   # @RestControllerAdvice
    │   │   ├── repository/
    │   │   │   └── TodoItemRepository.java        # JpaRepository<TodoItem, Long>
    │   │   └── service/
    │   │       ├── TodoItemService.java           # interface
    │   │       └── TodoItemServiceImpl.java       # @Service implementation
    │   └── resources/
    │       ├── application.yml                    # appsettings.json equivalent
    │       └── db/migration/
    │           └── V1__create_todo_items.sql      # Flyway migration
    └── test/
        └── java/com/example/todo/
            └── TodoApplicationTests.java
```

## Getting started

### Prerequisites

- Java 21+
- Maven 3.9+

### 1. Build the project

```bash
cd src/backend/java
mvn clean package -DskipTests
```

### 2. Run the application

```bash
mvn spring-boot:run
```

Or run the packaged JAR:

```bash
java -jar target/todo-api-1.0.0.jar
```

The application starts on <http://localhost:8080>.  
Swagger UI → <http://localhost:8080/swagger>  
H2 Console (dev) → <http://localhost:8080/h2-console>

Flyway automatically applies `V1__create_todo_items.sql` on startup (like `dotnet ef database update`).

## API endpoints

| Method | URL | Description |
|--------|-----|-------------|
| `GET` | `/api/todo-items` | List all todo items (paginated) |
| `GET` | `/api/todo-items/incomplete` | List incomplete items (paginated) |
| `GET` | `/api/todo-items/{id}` | Get a single todo item |
| `POST` | `/api/todo-items` | Create a todo item |
| `PUT` | `/api/todo-items/{id}` | Update a todo item |
| `PATCH` | `/api/todo-items/{id}/complete` | Mark a todo item as complete |
| `DELETE` | `/api/todo-items/{id}` | Delete a todo item |

### Pagination query parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `page` | `1` | Page number (1-based) |
| `pageSize` | `20` | Items per page |

## Switching databases

### PostgreSQL

1. Uncomment the PostgreSQL driver in `pom.xml`.
2. Start the app with the `postgres` profile:

```bash
mvn spring-boot:run -Dspring-boot.run.profiles=postgres \
  -Dspring-boot.run.arguments="--DB_USERNAME=myuser --DB_PASSWORD=mypass"
```

### MySQL

1. Uncomment the MySQL driver in `pom.xml`.
2. Change `V1__create_todo_items.sql` — `BIGINT AUTO_INCREMENT` is already MySQL-compatible.
3. Add a MySQL datasource profile to `application.yml`.

## Docker

### Build the image

```bash
# Run from src/backend/java/
docker build -t todo-api-java .
```

### Run the container

```bash
docker run -d -p 8080:8080 --name todo-api-java todo-api-java
```

The API is available at <http://localhost:8080>.  
Swagger UI: <http://localhost:8080/swagger>  
H2 Console (embedded DB): <http://localhost:8080/h2-console>

### Persist the H2 database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8080:8080 -v todo-java-data:/app --name todo-api-java todo-api-java
```

### Stop and remove the container

```bash
docker stop todo-api-java
docker rm todo-api-java
```
