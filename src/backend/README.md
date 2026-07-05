# Backend

This directory contains multiple backend implementations of the same Todo API in different languages/frameworks. All implementations share the same database schema described below.

## Database Schema

Naming conventions may differ across implementations (e.g., `snake_case` vs `PascalCase` vs `camelCase`), but all projects must implement the same tables, columns, data types, and constraints.

---

### Table: `todo_items`

Stores individual todo tasks.

| Column | Data Type | Constraints | Notes |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key |
| `title` | VARCHAR(200) | NOT NULL | Short description of the task |
| `description` | TEXT | NULLABLE | Optional longer description |
| `is_completed` | BOOLEAN | NOT NULL, DEFAULT `false` | Completion flag |
| `created_at` | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()` | Set by the database on insert |
| `updated_at` | TIMESTAMP WITH TIME ZONE | NULLABLE | Set by the database on update |

---

### Table: `email_logs`

Audit trail for every outbound email attempt. Records persist even when SMTP delivery fails.

| Column | Data Type | Constraints | Notes |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key |
| `recipient` | VARCHAR(255) | NOT NULL | Destination email address |
| `subject` | VARCHAR(500) | NOT NULL | Email subject line |
| `body` | TEXT | NOT NULL | Full email body content |
| `status` | VARCHAR(50) | NOT NULL, DEFAULT `'pending'` | Allowed values: `pending`, `sent`, `failed` |
| `created_at` | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()` | Set by the database on insert |
| `sent_at` | TIMESTAMP WITH TIME ZONE | NULLABLE | Populated when delivery succeeds |
| `error_message` | TEXT | NULLABLE | Populated when delivery fails |

---

### Table: `files`

Stores metadata about uploaded files.

| Column | Data Type | Constraints | Notes |
|---|---|---|---|
| `id` | INTEGER | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key |
| `name` | VARCHAR(255) | NOT NULL | Original file name (without path) |
| `extension` | VARCHAR(20) | NOT NULL | File extension, without the leading dot (e.g. `pdf`, `png`) |
| `size` | BIGINT | NOT NULL | File size in bytes |
| `content_type` | VARCHAR(100) | NULLABLE | MIME type of the file |
| `location` | VARCHAR(500) | NOT NULL | Storage path or URL where the file content is stored |
| `created_at` | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()` | Set by the database on insert |
| `updated_at` | TIMESTAMP WITH TIME ZONE | NULLABLE | Set by the database on update |

---

## Naming Convention Mapping

| Canonical (SQL) | Python / SQLAlchemy | .NET / EF Core | Go | Java / JPA | Node.js / Prisma | PHP / Eloquent |
|---|---|---|---|---|---|---|
| `todo_items` | `todo_items` | `TodoItems` | `todo_items` | `todo_items` | `todo_items` | `todo_items` |
| `email_logs` | `email_logs` | `EmailLogs` | `email_logs` | `email_logs` | `email_logs` | `email_logs` |
| `files` | `files` | `Files` | `files` | `files` | `files` | `files` |
| `is_completed` | `is_completed` | `IsCompleted` | `is_completed` | `is_completed` | `isCompleted` | `is_completed` |
| `created_at` | `created_at` | `CreatedAt` | `created_at` | `created_at` | `createdAt` | `created_at` |
| `updated_at` | `updated_at` | `UpdatedAt` | `updated_at` | `updated_at` | `updatedAt` | `updated_at` |
| `sent_at` | `sent_at` | `SentAt` | `sent_at` | `sent_at` | `sentAt` | `sent_at` |
| `error_message` | `error_message` | `ErrorMessage` | `error_message` | `error_message` | `errorMessage` | `error_message` |
| `content_type` | `content_type` | `ContentType` | `content_type` | `content_type` | `contentType` | `content_type` |

## API Endpoints

All implementations expose the same REST endpoints under the `/api/todo-items` prefix. Path parameter syntax may differ by framework (`{id}` vs `:id`), but the routes are functionally identical.

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/api/todo-items` | List all todo items (paginated) |
| `GET` | `/api/todo-items/incomplete` | List incomplete items (paginated) |
| `GET` | `/api/todo-items/{id}` | Get a single todo item |
| `POST` | `/api/todo-items` | Create a todo item |
| `PUT` | `/api/todo-items/{id}` | Update a todo item |
| `PATCH` | `/api/todo-items/{id}/complete` | Mark a todo item as complete |
| `DELETE` | `/api/todo-items/{id}` | Delete a todo item |
| `POST` | `/api/todo-items/import/csv` | Import todo items from a CSV file (`multipart/form-data`) |
| `POST` | `/api/todo-items/import/excel` | Import todo items from an Excel file (`multipart/form-data`) |
| `GET` | `/api/todo-items/export/csv` | Export todo items as a CSV file |
| `GET` | `/api/todo-items/export/excel` | Export todo items as an Excel file |
| `GET` | `/api/files` | List all uploaded files (paginated) |
| `GET` | `/api/files/{id}` | Get a single file's metadata |
| `GET` | `/api/files/{id}/download` | Download a file's content |
| `POST` | `/api/files` | Upload a file (`multipart/form-data`) |
| `DELETE` | `/api/files/{id}` | Delete a file |

### Pagination Query Parameters

All paginated endpoints (`GET /api/todo-items`, `GET /api/todo-items/incomplete`, and `GET /api/files`) accept the following query parameters. The parameter name casing may differ per implementation but must convey the same semantics.

| Parameter (canonical) | Default | Description | Implementation note |
|---|---|---|---|
| `page` | `1` | Page number (1-based) | Same across all implementations |
| `page_size` / `pageSize` | `20` | Items per page | `page_size` in Python & PHP; `pageSize` in .NET, Go, Java, Node.js |

---

## Implementations

| Language / Framework | Path |
|---|---|
| Python (FastAPI + SQLAlchemy) | `python/` |
| .NET (ASP.NET Core + EF Core) | `dotnet/` |
| Go (net/http) | `go/` |
| Java (Spring Boot) | `java/` |
| Node.js (NestJS + Prisma) | `nodejs/` |
| PHP (Laravel) | `php/` |
