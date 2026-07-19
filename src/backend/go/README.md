# Todo API - Go / Gin + GORM

A RESTful API for managing todo items built with **Gin**, **GORM**, and Go - the Go equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF                     | Go equivalent                                                 |
| ------------------------------------- | ------------------------------------------------------------- |
| ASP.NET Core                          | **Gin** (HTTP framework)                                      |
| Entity Framework Core                 | **GORM** (ORM)                                                |
| EF Migrations                         | GORM `AutoMigrate` / `golang-migrate`                         |
| `DbContext`                           | `*gorm.DB` passed via constructor DI                          |
| Controllers                           | **Handlers** (`TodoItemHandler`)                              |
| Services (`IService` / `Service`)     | `TodoItemService` interface + impl                            |
| Repository pattern                    | `TodoItemRepository` interface + impl                         |
| DTOs / Data Annotations               | Go structs with `binding` tags (Gin validator)                |
| `appsettings.json` / `IConfiguration` | `godotenv` + `config.Config` struct                           |
| Dependency Injection                  | Manual constructor injection (composition root in `main.go`)  |
| Swagger / OpenAPI                     | `swaggo/gin-swagger` at `/swagger/index.html`                 |
| `Program.cs`                          | `cmd/api/main.go`                                             |
| xUnit + Moq                           | `testing` (stdlib) + `testify` (hand-written interface mocks) |

## Project structure

```
src/backend/go/
├── cmd/
│   ├── api/
│   │   └── main.go                        # API entry point (Program.cs)
│   └── worker/
│       └── main.go                        # Background worker entry point
├── internal/
│   ├── config/
│   │   ├── config.go                      # Settings (appsettings.json)
│   │   └── config_test.go                 # Config unit tests
│   ├── database/
│   │   └── database.go                    # GORM setup (DbContext)
│   ├── models/
│   │   ├── todo_item.go                   # GORM entity - TodoItem
│   │   ├── todo_item_attachment.go        # GORM entity - todo-to-file attachment
│   │   ├── email_log.go                   # GORM entity - EmailLog
│   │   ├── file.go                        # GORM entity - File (uploaded-file metadata)
│   │   └── user.go                        # GORM entity - User account
│   ├── dto/
│   │   ├── todo_item.go                   # Request/response DTOs
│   │   ├── todo_item_attachment.go        # Attachment request/response DTOs
│   │   ├── file.go                        # File metadata response DTO
│   │   └── user.go                        # User and password request/response DTOs
│   ├── repository/
│   │   ├── repository.go                  # Repository interfaces
│   │   ├── todo_item_repository.go        # GORM implementation - TodoItem
│   │   ├── todo_item_attachment_repository.go # GORM implementation - attachments
│   │   ├── email_log_repository.go        # GORM implementation - EmailLog
│   │   ├── file_repository.go             # GORM implementation - File
│   │   ├── user_repository.go             # GORM implementation - User
│   │   └── user_repository_test.go        # User persistence tests
│   ├── service/
│   │   ├── todo_item_service.go           # Business logic
│   │   ├── todo_item_service_test.go      # Service unit tests
│   │   ├── file_service.go                # Business logic (upload/download/delete on disk)
│   │   ├── file_service_test.go           # Service unit tests
│   │   ├── todo_item_attachment_service.go # Attachment business logic
│   │   ├── todo_item_attachment_service_test.go # Attachment service tests
│   │   ├── user_service.go                # User, profile, and password business logic
│   │   └── user_service_test.go           # User service tests
│   ├── handler/
│   │   ├── todo_item_handler.go           # HTTP handlers (Controller)
│   │   ├── todo_item_handler_test.go      # Handler unit tests
│   │   ├── file_handler.go                # HTTP handlers - /api/files (Controller)
│   │   ├── file_handler_test.go           # Handler unit tests
│   │   ├── todo_item_attachment_handler.go # Nested attachment HTTP handlers
│   │   ├── todo_item_attachment_handler_test.go # Attachment handler tests
│   │   ├── user_handler.go                # HTTP handlers - /api/users
│   │   └── user_handler_test.go           # User endpoint handler tests
│   ├── security/
│   │   ├── user_security.go               # Password hashing and golang-jwt token handling
│   │   └── user_security_test.go          # Password and token tests
│   ├── router/
│   │   └── router.go                      # Route registration
│   └── worker/
│       └── job/
│           └── incomplete_todos_email.go  # Email digest job
├── build/
│   ├── api/
│   │   └── Dockerfile                     # API container image
│   └── worker/
│       └── Dockerfile                     # Background worker container image
├── go.mod
├── .env.example
└── README.md
```

## Getting started

### 1. Install Go

**Windows (winget):**

```powershell
winget install GoLang.Go
```

**macOS (Homebrew):**

```bash
brew install go
```

**Linux:**

```bash
sudo apt install golang-go        # Debian/Ubuntu
sudo dnf install golang           # Fedora/RHEL
```

Or download the installer directly: <https://go.dev/dl/>

Ensure Go 1.25+ is installed (`go version` to verify).

### 2. Install dependencies

```bash
cd src/backend/go
go mod tidy
```

### 3. Configure environment

```bash
copy .env.example .env
# Edit .env as needed
```

### 4. Generate Swagger docs (optional)

```bash
go install github.com/swaggo/swag/cmd/swag@latest
swag init -g cmd/api/main.go -o docs
```

### 5. Run unit tests

```bash
# Run all unit tests
go test ./...

# Verbose output
go test -v ./...

# With coverage report
go test -cover ./...
```

> **Note:** The `internal/database` package uses CGO (SQLite). Ensure `gcc` is available, or
> target only pure-Go packages: `go test ./internal/config/... ./internal/service/... ./internal/handler/...`

### 6. Run the server

```bash
go run ./cmd/api
```

The API will start on <http://localhost:8080>.  
Swagger UI → <http://localhost:8080/swagger/index.html>

### Build for production

```bash
go build -o todo-api ./cmd/api
./todo-api
```

## API endpoints

See the [shared API contract](../README.md#api-endpoints) in the backend README.

## File uploads

The `/api/files` endpoints (list, get metadata, download, upload, delete) store uploaded file content on disk and persist metadata (`name`, `extension`, `size`, `content_type`, `location`, timestamps) in the `files` table.

### Configuration (`.env`)

```ini
FILE_STORAGE_PATH=./uploads     # Directory where uploaded file content is stored
MAX_UPLOAD_SIZE_BYTES=10485760  # Maximum accepted upload size, in bytes (default 10 MB)
```

### Notes

- Uploaded file names are sanitized (directory components stripped) and stored on disk under a random-prefixed name to prevent path traversal and filename collisions.
- The internal storage `location` is never exposed in API responses; file content is retrieved via `GET /api/files/{id}/download`.
- Deleting a file removes both the database row and the file content on disk.
- Uploads exceeding `MAX_UPLOAD_SIZE_BYTES` are rejected with `413 Request Entity Too Large`.

## Switching databases

1. Change `DATABASE_DSN` in `.env`
2. Swap the driver import in `internal/database/database.go`:

**PostgreSQL:**

```go
import "gorm.io/driver/postgres"
// ...
db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
```

```bash
go get gorm.io/driver/postgres
```

**MySQL:**

```go
import "gorm.io/driver/mysql"
// ...
db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
```

```bash
go get gorm.io/driver/mysql
```

## Docker

### Build the API image

```bash
# Run from src/backend/go/
docker build -f build/api/Dockerfile -t todo-api-go .
```

### Build the worker image

```bash
docker build -f build/worker/Dockerfile -t todo-worker-go .
```

### Run the containers

```bash
docker run -d -p 8080:8080 --name todo-api-go todo-api-go
```

The API is available at <http://localhost:8080>.  
Swagger UI: <http://localhost:8080/swagger/index.html>

### Persist the SQLite database

Both the API and the worker need access to the same database file.
Mount a shared named volume so the data survives container restarts:

```bash
# API
docker run -d -p 8080:8080 -v todo-go-data:/app --name todo-api-go todo-api-go

# Worker (shares the same volume)
docker run -d \
  -v todo-go-data:/app \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_USERNAME=user@example.com \
  -e SMTP_PASSWORD=secret \
  -e EMAIL_SENDER=noreply@example.com \
  -e EMAIL_RECIPIENT=admin@example.com \
  -e WORKER_INTERVAL_MINUTES=60 \
  --name todo-worker-go todo-worker-go
```

### Stop and remove the containers

```bash
docker stop todo-api-go todo-worker-go
docker rm  todo-api-go todo-worker-go
```

## Background worker

The worker runs as a **separate process / container** (`cmd/worker/main.go`) and shares the same SQLite database as the API via a mounted volume.

### What it does

1. Queries all incomplete todo items (`is_completed = false`).
2. Builds a plain-text email digest listing every pending item.
3. Inserts an `email_logs` row with `status = 'pending'`.
4. Sends the email via SMTP (STARTTLS or plain depending on `SMTP_USE_TLS`).
5. Updates the `email_logs` row to `status = 'sent'` on success, or `status = 'failed'` with `error_message` on failure.

The job runs **immediately on startup** and then on the configured interval.

### Worker configuration

| Variable                  | Default               | Description                                       |
| ------------------------- | --------------------- | ------------------------------------------------- |
| `WORKER_INTERVAL_MINUTES` | `60`                  | How often to run the digest job                   |
| `SMTP_HOST`               | `localhost`           | SMTP server hostname                              |
| `SMTP_PORT`               | `587`                 | SMTP server port                                  |
| `SMTP_USERNAME`           | _(empty)_             | SMTP login (omit for anonymous)                   |
| `SMTP_PASSWORD`           | _(empty)_             | SMTP password                                     |
| `SMTP_USE_TLS`            | `true`                | `true` = STARTTLS (port 587); `false` = plain/SSL |
| `EMAIL_SENDER`            | `noreply@example.com` | From address                                      |
| `EMAIL_RECIPIENT`         | `admin@example.com`   | Destination address                               |

### Run the worker locally

```bash
# Uses the same .env as the API
go run ./cmd/worker
```
