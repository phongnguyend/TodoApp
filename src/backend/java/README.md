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
| xUnit / NUnit | **JUnit 5** + **Mockito** + **Spring Boot Test** (`@WebMvcTest`, `@DataJpaTest`) |

## Project structure

The project follows a multi-module Maven layout (analogous to a .NET solution with separate `TodoApi` and `TodoWorker` projects sharing a common library):

```
src/backend/java/
├── pom.xml                                            # Parent POM (like a .NET solution file)
├── todo-shared/                                       # Shared library (analogous to a shared .csproj)
│   ├── pom.xml
│   └── src/
│       ├── main/
│       │   ├── java/com/example/todo/
│       │   │   ├── dto/
│       │   │   │   ├── CreateTodoItemRequest.java     # record + @Valid
│       │   │   │   ├── UpdateTodoItemRequest.java     # record + @Valid
│       │   │   │   ├── TodoItemResponse.java          # record with factory method
│       │   │   │   └── PaginatedResponse.java         # generic record
│       │   │   ├── entity/
│       │   │   │   ├── TodoItem.java                  # @Entity (JPA model)
│       │   │   │   └── EmailLog.java                  # @Entity — email audit log
│       │   │   ├── repository/
│       │   │   │   ├── TodoItemRepository.java        # JpaRepository<TodoItem, Long>
│       │   │   │   └── EmailLogRepository.java        # JpaRepository<EmailLog, Long>
│       │   │   └── service/
│       │   │       ├── TodoItemService.java           # interface
│       │   │       └── TodoItemServiceImpl.java       # @Service implementation
│       │   └── resources/db/migration/
│       │       ├── V1__create_todo_items.sql          # Flyway migration
│       │       └── V2__create_email_logs.sql          # Flyway migration
│       └── test/
│           ├── java/com/example/todo/
│           │   ├── repository/TodoItemRepositoryTest.java  # @DataJpaTest — JPA slice
│           │   └── service/TodoItemServiceImplTest.java    # Mockito — pure unit tests
│           └── resources/application.yml                   # test config (H2 in-memory)
├── todo-api/                                          # REST API (analogous to TodoApi.csproj)
│   ├── pom.xml
│   ├── Dockerfile
│   └── src/
│       ├── main/
│       │   ├── java/com/example/todo/api/
│       │   │   ├── TodoApiApplication.java            # Entry point (@SpringBootApplication)
│       │   │   ├── config/OpenApiConfig.java          # Swagger config
│       │   │   ├── controller/TodoItemController.java # @RestController
│       │   │   └── exception/GlobalExceptionHandler.java  # @RestControllerAdvice
│       │   └── resources/application.yml             # API config (H2 + Swagger)
│       └── test/
│           ├── java/com/example/todo/api/
│           │   ├── TodoApiApplicationTests.java       # context load smoke test
│           │   └── controller/TodoItemControllerTest.java  # @WebMvcTest — HTTP layer
│           └── resources/application.yml             # test config (H2 in-memory)
└── todo-worker/                                       # Background worker (analogous to TodoWorker.csproj)
    ├── pom.xml
    ├── Dockerfile
    └── src/main/
        ├── java/com/example/todo/worker/
        │   ├── TodoWorkerApplication.java             # Entry point (no web server)
        │   ├── IncompleteTodosEmailJob.java           # Job: query → build email → persist → send
        │   └── WorkerRunner.java                      # ApplicationRunner: scheduling loop
        └── resources/application.yml                 # Worker config (no web server, SMTP)
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

### 2. Run unit tests

```bash
# Run all tests across all modules
mvn test

# Run tests for a specific module
mvn test -pl todo-shared
mvn test -pl todo-api

# Run a specific test class
mvn test -pl todo-api -Dtest=TodoItemControllerTest
mvn test -pl todo-shared -Dtest=TodoItemServiceImplTest
```

### 3. Run the API

```bash
mvn spring-boot:run -pl todo-api
```

Or run the packaged JAR:

```bash
java -jar todo-api/target/todo-api-1.0.0.jar
```

The application starts on <http://localhost:8080>.  
Swagger UI → <http://localhost:8080/swagger>  
H2 Console (dev) → <http://localhost:8080/h2-console>

Flyway automatically applies `V1__create_todo_items.sql` on startup (like `dotnet ef database update`).

## API endpoints

See the [shared API contract](../README.md#api-endpoints) in the backend README.

## Background worker

The worker is a separate process (separate container) that periodically sends an email digest of all incomplete todo items.

### How it works

1. Queries all incomplete `todo_items` ordered by `created_at`.
2. Builds a plain-text + HTML email body.
3. Inserts an `email_logs` row with `status = 'pending'`.
4. Sends the email via SMTP (`JavaMailSender`).
5. Updates the `email_logs` row to `status = 'sent'` (or `'failed'` + `error_message`).

The worker runs as a standalone Spring Boot application (`todo-worker` module) with `spring.main.web-application-type=none` in its `application.yml` — no web server is started.

### Environment variables

| Variable | Default | Description |
|---|---|---|
| `WORKER_INTERVAL_MINUTES` | `60` | How often the job runs |
| `EMAIL_RECIPIENT` | `admin@example.com` | Destination address for the digest |
| `EMAIL_SENDER` | `noreply@example.com` | From address |
| `SMTP_HOST` | `localhost` | SMTP server hostname |
| `SMTP_PORT` | `587` | SMTP server port |
| `SMTP_USERNAME` | _(empty)_ | SMTP auth username |
| `SMTP_PASSWORD` | _(empty)_ | SMTP auth password |
| `SMTP_AUTH` | `false` | Enable SMTP authentication |
| `SMTP_STARTTLS` | `false` | Enable STARTTLS |

### Run the worker locally

```bash
mvn spring-boot:run -pl todo-worker
```

### Run with custom interval and SMTP

```bash
java -jar todo-worker/target/todo-worker-1.0.0.jar \
  --WORKER_INTERVAL_MINUTES=30 \
  --SMTP_HOST=smtp.example.com \
  --SMTP_PORT=587 \
  --SMTP_USERNAME=user@example.com \
  --SMTP_PASSWORD=secret \
  --SMTP_AUTH=true \
  --SMTP_STARTTLS=true \
  --EMAIL_RECIPIENT=team@example.com
```

## Switching databases

### PostgreSQL

1. Uncomment the PostgreSQL driver in `todo-shared/pom.xml`.
2. Start the app with the `postgres` profile:

```bash
mvn spring-boot:run -pl todo-api -Dspring-boot.run.profiles=postgres \
  -Dspring-boot.run.arguments="--DB_USERNAME=myuser --DB_PASSWORD=mypass"
```

### MySQL

1. Uncomment the MySQL driver in `todo-shared/pom.xml`.
2. Change `V1__create_todo_items.sql` — `BIGINT AUTO_INCREMENT` is already MySQL-compatible.
3. Add a MySQL datasource profile to `todo-api/src/main/resources/application.yml`.

## Docker

### Build the API image

```bash
# Run from src/backend/java/
docker build -f todo-api/Dockerfile -t todo-api-java .
```

### Build the worker image

```bash
# Run from src/backend/java/
docker build -f todo-worker/Dockerfile -t todo-worker-java .
```

### Run the API container

```bash
docker run -d -p 8080:8080 --name todo-api-java todo-api-java
```

The API is available at <http://localhost:8080>.  
Swagger UI: <http://localhost:8080/swagger>  
H2 Console (embedded DB): <http://localhost:8080/h2-console>

### Run the worker container

The worker shares the same H2 database file as the API. Mount the same volume so both processes access the same data:

```bash
# 1. Start the API with a named volume
docker run -d -p 8080:8080 -v todo-java-data:/app --name todo-api-java todo-api-java

# 2. Start the worker pointing at the same volume and your SMTP server
docker run -d \
  -v todo-java-data:/app \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=user@example.com \
  -e SMTP_PASSWORD=secret \
  -e SMTP_AUTH=true \
  -e SMTP_STARTTLS=true \
  -e EMAIL_RECIPIENT=team@example.com \
  -e WORKER_INTERVAL_MINUTES=60 \
  --name todo-worker-java \
  todo-worker-java
```

### Persist the H2 database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8080:8080 -v todo-java-data:/app --name todo-api-java todo-api-java
```

### Stop and remove the containers

```bash
docker stop todo-api-java todo-worker-java
docker rm todo-api-java todo-worker-java
```
