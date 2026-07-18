# Todo API - ASP.NET Core / Entity Framework Core

A RESTful API for managing todo items built with **ASP.NET Core**, **Entity Framework Core**, and **SQLite** (swappable to SQL Server or PostgreSQL).

## Tech stack

| Concern              | Technology                                          |
| -------------------- | --------------------------------------------------- |
| Web framework        | ASP.NET Core 10 Web API                             |
| ORM                  | Entity Framework Core 10                            |
| Database (default)   | SQLite (via `Microsoft.EntityFrameworkCore.Sqlite`) |
| Migrations           | EF Core Migrations (`dotnet ef`)                    |
| API docs             | Built-in OpenAPI + **Scalar** UI                    |
| Dependency Injection | Built-in ASP.NET Core DI                            |
| Unit testing         | xUnit + Moq                                         |

## Project structure

```
src/backend/dotnet/
├── TodoApp.slnx
├── TodoApi/
│   ├── Dockerfile                          # API container image
│   ├── Program.cs                          # App bootstrap & DI registration
│   ├── appsettings.json                    # Connection string, file storage & logging
│   ├── Controllers/
│   │   ├── TodoItemsController.cs          # REST endpoints - /api/todo-items
│   │   └── FilesController.cs              # REST endpoints - /api/files
│   ├── Data/
│   │   ├── AppDbContext.cs                 # EF Core DbContext
│   │   └── Migrations/                     # EF Core migrations
│   ├── DTOs/
│   │   ├── TodoItemDtos.cs                 # Request / response models
│   │   ├── TodoItemAttachmentDtos.cs       # Attachment request / response models
│   │   └── FileDtos.cs                     # FileResponse / FileDownloadTarget models
│   ├── Repositories/
│   │   ├── IRepository.cs                  # Generic IRepository<T>
│   │   ├── BaseRepository.cs               # Generic BaseRepository<T>
│   │   ├── ITodoItemRepository.cs
│   │   ├── TodoItemRepository.cs
│   │   ├── IEmailLogRepository.cs
│   │   ├── EmailLogRepository.cs
│   │   ├── IFileRepository.cs
│   │   ├── FileRepository.cs
│   │   ├── ITodoItemAttachmentRepository.cs
│   │   └── TodoItemAttachmentRepository.cs
│   └── Services/
│       ├── ITodoItemService.cs
│       ├── TodoItemService.cs
│       ├── IFileService.cs                 # upload/download/delete on disk
│       ├── FileService.cs
│       ├── ITodoItemAttachmentService.cs
│       ├── TodoItemAttachmentService.cs
│       └── FileTooLargeException.cs        # thrown when upload exceeds MaxUploadSizeBytes
├── TodoShared/
│   ├── Data/
│   │   └── TodoDbContext.cs                # Shared EF Core model configuration
│   └── Models/
│       ├── TodoItem.cs                     # Todo item entity
│       ├── TodoItemAttachment.cs           # Todo item attachment entity
│       ├── EmailLog.cs                     # Email audit log entity
│       └── FileEntity.cs                   # File metadata entity
├── TodoWorker/
│   ├── Dockerfile                          # Worker container image
│   ├── Program.cs                          # Worker bootstrap & DI registration
│   ├── appsettings.json                    # Connection string, SMTP & worker settings
│   ├── Data/
│   │   └── WorkerDbContext.cs              # EF Core DbContext (read-only schema, no migrations)
│   ├── Models/
│   │   ├── TodoItem.cs                     # Local POCO matching TodoItems table
│   │   └── EmailLog.cs                     # Local POCO matching EmailLogs table
│   └── Services/
│       ├── IEmailService.cs
│       ├── SmtpEmailService.cs             # SMTP delivery via System.Net.Mail
│       └── WorkerService.cs               # BackgroundService - periodic email job
└── TodoApi.Tests/
    ├── Controllers/
    │   ├── TodoItemsControllerTests.cs     # Controller unit tests
    │   └── FilesControllerTests.cs         # Controller unit tests
    └── Services/
        ├── TodoItemServiceTests.cs         # Service unit tests
        ├── TodoItemAttachmentServiceTests.cs # Attachment service tests
        └── FileServiceTests.cs             # Service unit tests
```

## Getting started

### Prerequisites

- [.NET 10 SDK](https://dotnet.microsoft.com/download)
- EF Core CLI tools: `dotnet tool install --global dotnet-ef`

### 1. Restore dependencies

```bash
cd src/backend/dotnet
dotnet restore
```

### 2. Apply database migrations

```bash
cd TodoApi
dotnet ef database update
```

### 3. Run the API

```bash
dotnet run
```

The API starts on `https://localhost:7xxx` / `http://localhost:5xxx`.  
Scalar API reference UI: `https://localhost:7xxx/scalar/v1`  
OpenAPI JSON: `https://localhost:7xxx/openapi/v1.json`

### 4. Run the background worker

```bash
cd ../TodoWorker
dotnet run
```

The worker connects to the same `todo.db` database. Configure SMTP and recipient settings in `appsettings.json` before running.

### 5. Run unit tests

```bash
# Run all tests
dotnet test

# Run with detailed output
dotnet test --verbosity normal

# Run with code coverage
dotnet test --collect:"XPlat Code Coverage"
```

## API endpoints

See the [shared API contract](../README.md#api-endpoints) in the backend README.

## File uploads

The `/api/files` endpoints (list, get metadata, download, upload, delete) store uploaded file content on disk and persist metadata (`Name`, `Extension`, `Size`, `ContentType`, `Location`, timestamps) in the `Files` table.

### Configuration (`appsettings.json`)

```json
{
  "FileStorage": {
    "Path": "uploads",
    "MaxUploadSizeBytes": 10485760
  }
}
```

`Path` (relative or absolute) is the directory where uploaded file content is stored; it is created automatically on first upload. `MaxUploadSizeBytes` is the maximum accepted upload size, in bytes (default 10 MB). Both values can be overridden with environment variables (e.g. `FileStorage__Path`, `FileStorage__MaxUploadSizeBytes`).

### Notes

- Uploaded file names are sanitized (directory components stripped via `Path.GetFileName`) and stored on disk under a random-prefixed name (GUID) to prevent path traversal and filename collisions.
- The internal storage `Location` is never exposed in API responses; file content is retrieved via `GET /api/files/{id}/download`.
- Deleting a file removes both the database row and the file content on disk.
- Uploads exceeding `MaxUploadSizeBytes` are rejected with `413 Payload Too Large`.

## Background worker

`TodoWorker` is a separate .NET Worker Service that runs as an independent process (and container). It shares the same SQLite database as the API.

### What it does

On startup and then on a configurable interval (default: every 5 minutes), the worker:

1. Queries the `TodoItems` table for all incomplete items (`IsCompleted = false`).
2. If none are found, it skips and waits for the next tick.
3. Otherwise it composes a summary email and inserts an `email_logs` row with `status = pending`.
4. Sends the email via SMTP.
5. Updates the row to `status = sent` (with `sent_at`) on success, or `status = failed` (with `error_message`) on failure.

### Configuration (`appsettings.json`)

```json
{
  "ConnectionStrings": {
    "DefaultConnection": "Data Source=todo.db"
  },
  "Smtp": {
    "Host": "localhost",
    "Port": 25,
    "EnableSsl": false,
    "From": "noreply@todo.app",
    "Username": "",
    "Password": ""
  },
  "Worker": {
    "IntervalMinutes": 5,
    "RecipientEmail": "admin@todo.app"
  }
}
```

All values can be overridden with environment variables (e.g. `Smtp__Host`, `Worker__RecipientEmail`).

## EF Core migration commands

```bash
# Add a new migration
dotnet ef migrations add <MigrationName> --output-dir Data/Migrations

# Apply migrations
dotnet ef database update

# Revert last migration
dotnet ef migrations remove
```

## Switching databases

Update `ConnectionStrings:DefaultConnection` in `appsettings.json` and swap the EF provider package:

| Database         | Package                                   | Connection string                                             |
| ---------------- | ----------------------------------------- | ------------------------------------------------------------- |
| SQLite (default) | `Microsoft.EntityFrameworkCore.Sqlite`    | `Data Source=todo.db`                                         |
| SQL Server       | `Microsoft.EntityFrameworkCore.SqlServer` | `Server=.;Database=TodoDb;Trusted_Connection=True`            |
| PostgreSQL       | `Npgsql.EntityFrameworkCore.PostgreSQL`   | `Host=localhost;Database=todo_db;Username=user;Password=pass` |

## Docker

### Build the API image

```bash
# Run from src/backend/dotnet/TodoApi/
docker build -t todo-api-dotnet .
```

### Build the Worker image

```bash
# Run from src/backend/dotnet/TodoWorker/
docker build -t todo-worker-dotnet .
```

### Run the API container

```bash
docker run -d -p 8080:8080 --name todo-api-dotnet todo-api-dotnet
```

The API is available at <http://localhost:8080>.  
Scalar API reference UI: <http://localhost:8080/scalar/v1>  
OpenAPI JSON: <http://localhost:8080/openapi/v1.json>

### Run the Worker container

The worker must share the same database file as the API. Use a named volume and set SMTP settings via environment variables:

```bash
# Create a shared volume (once)
docker volume create todo-dotnet-data

# Start the API with the shared volume
docker run -d -p 8080:8080 -v todo-dotnet-data:/app --name todo-api-dotnet todo-api-dotnet

# Start the worker with the same volume and SMTP config
docker run -d \
  -v todo-dotnet-data:/app \
  -e Smtp__Host=mailhog \
  -e Smtp__Port=1025 \
  -e Worker__RecipientEmail=you@example.com \
  --name todo-worker-dotnet \
  todo-worker-dotnet
```

### Persist the SQLite database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8080:8080 -v todo-dotnet-data:/app --name todo-api-dotnet todo-api-dotnet
```

### Stop and remove the containers

```bash
docker stop todo-api-dotnet todo-worker-dotnet
docker rm todo-api-dotnet todo-worker-dotnet
```
