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

## Naming Convention Mapping

| Canonical (SQL) | Python / SQLAlchemy | .NET / EF Core | Go | Java / JPA | Node.js / Prisma | PHP / Eloquent |
|---|---|---|---|---|---|---|
| `todo_items` | `todo_items` | `TodoItems` | `todo_items` | `todo_items` | `todo_items` | `todo_items` |
| `email_logs` | `email_logs` | `EmailLogs` | `email_logs` | `email_logs` | `email_logs` | `email_logs` |
| `is_completed` | `is_completed` | `IsCompleted` | `is_completed` | `is_completed` | `isCompleted` | `is_completed` |
| `created_at` | `created_at` | `CreatedAt` | `created_at` | `created_at` | `createdAt` | `created_at` |
| `updated_at` | `updated_at` | `UpdatedAt` | `updated_at` | `updated_at` | `updatedAt` | `updated_at` |
| `sent_at` | `sent_at` | `SentAt` | `sent_at` | `sent_at` | `sentAt` | `sent_at` |
| `error_message` | `error_message` | `ErrorMessage` | `error_message` | `error_message` | `errorMessage` | `error_message` |

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

### Pagination Query Parameters

All paginated endpoints (`GET /api/todo-items` and `GET /api/todo-items/incomplete`) accept the following query parameters. The parameter name casing may differ per implementation but must convey the same semantics.

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
