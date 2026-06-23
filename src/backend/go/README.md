# Todo API — Go / Gin + GORM

A RESTful API for managing todo items built with **Gin**, **GORM**, and Go — the Go equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF | Go equivalent |
|---|---|
| ASP.NET Core | **Gin** (HTTP framework) |
| Entity Framework Core | **GORM** (ORM) |
| EF Migrations | GORM `AutoMigrate` / `golang-migrate` |
| `DbContext` | `*gorm.DB` passed via constructor DI |
| Controllers | **Handlers** (`TodoItemHandler`) |
| Services (`IService` / `Service`) | `TodoItemService` interface + impl |
| Repository pattern | `TodoItemRepository` interface + impl |
| DTOs / Data Annotations | Go structs with `binding` tags (Gin validator) |
| `appsettings.json` / `IConfiguration` | `godotenv` + `config.Config` struct |
| Dependency Injection | Manual constructor injection (composition root in `main.go`) |
| Swagger / OpenAPI | `swaggo/gin-swagger` at `/swagger/index.html` |
| `Program.cs` | `cmd/api/main.go` |
| xUnit + Moq | `testing` (stdlib) + `testify` (hand-written interface mocks) |

## Project structure

```
src/backend/go/
├── cmd/
│   └── api/
│       └── main.go                        # Entry point (Program.cs)
├── internal/
│   ├── config/
│   │   ├── config.go                      # Settings (appsettings.json)
│   │   └── config_test.go                 # Config unit tests
│   ├── database/
│   │   └── database.go                    # GORM setup (DbContext)
│   ├── models/
│   │   └── todo_item.go                   # GORM entity
│   ├── dto/
│   │   └── todo_item.go                   # Request/response DTOs
│   ├── repository/
│   │   ├── repository.go                  # Repository interface
│   │   └── todo_item_repository.go        # GORM implementation
│   ├── service/
│   │   ├── todo_item_service.go           # Business logic
│   │   └── todo_item_service_test.go      # Service unit tests
│   ├── handler/
│   │   ├── todo_item_handler.go           # HTTP handlers (Controller)
│   │   └── todo_item_handler_test.go      # Handler unit tests
│   └── router/
│       └── router.go                      # Route registration
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

Ensure Go 1.23+ is installed (`go version` to verify).

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

| Method | URL | Description |
|--------|-----|-------------|
| `GET` | `/api/todo-items` | List all todo items (paginated) |
| `GET` | `/api/todo-items/incomplete` | List incomplete items (paginated) |
| `GET` | `/api/todo-items/:id` | Get a single todo item |
| `POST` | `/api/todo-items` | Create a todo item |
| `PUT` | `/api/todo-items/:id` | Update a todo item |
| `PATCH` | `/api/todo-items/:id/complete` | Mark a todo item as complete |
| `DELETE` | `/api/todo-items/:id` | Delete a todo item |

### Pagination query parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `page` | `1` | Page number (1-based) |
| `pageSize` | `20` | Items per page |

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

### Build the image

```bash
# Run from src/backend/go/
docker build -t todo-api-go .
```

### Run the container

```bash
docker run -d -p 8080:8080 --name todo-api-go todo-api-go
```

The API is available at <http://localhost:8080>.  
Swagger UI: <http://localhost:8080/swagger/index.html>

### Persist the SQLite database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8080:8080 -v todo-go-data:/app --name todo-api-go todo-api-go
```

### Stop and remove the container

```bash
docker stop todo-api-go
docker rm todo-api-go
```
