# Backend

This directory contains multiple backend implementations of the same Todo API in different languages/frameworks. All implementations share the same database schema described below.

## Database Schema

Naming conventions may differ across implementations (e.g., `snake_case` vs `PascalCase` vs `camelCase`), but all projects must implement the same tables, columns, data types, and constraints.

---

### Table: `todo_items`

Stores individual todo tasks.

| Column         | Data Type                | Constraints                           | Notes                         |
| -------------- | ------------------------ | ------------------------------------- | ----------------------------- |
| `id`           | INTEGER                  | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key                 |
| `title`        | VARCHAR(200)             | NOT NULL                              | Short description of the task |
| `description`  | TEXT                     | NULLABLE                              | Optional longer description   |
| `is_completed` | BOOLEAN                  | NOT NULL, DEFAULT `false`             | Completion flag               |
| `created_at`   | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()`             | Set by the database on insert |
| `updated_at`   | TIMESTAMP WITH TIME ZONE | NULLABLE                              | Set by the database on update |

---

### Table: `todo_item_attachments`

Stores attachment references for a todo item. The attachment itself is not a duplicate of the file content; it only references a row in the `files` table.

| Column         | Data Type                | Constraints                              | Notes                                                        |
| -------------- | ------------------------ | ---------------------------------------- | ------------------------------------------------------------ |
| `id`           | INTEGER                  | PRIMARY KEY, AUTO INCREMENT, NOT NULL    | Surrogate key                                                |
| `todo_item_id` | INTEGER                  | NOT NULL, FOREIGN KEY -> `todo_items.id` | The todo item this attachment belongs to                     |
| `file_id`      | INTEGER                  | NOT NULL, FOREIGN KEY -> `files.id`      | Reference to the uploaded file metadata in the `files` table |
| `created_at`   | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()`                | Set by the database on insert                                |
| `updated_at`   | TIMESTAMP WITH TIME ZONE | NULLABLE                                 | Set by the database on update                                |

Additional constraints:

- Unique constraint on `(todo_item_id, file_id)` to prevent duplicate attachments for the same file on the same todo item.

---

### Table: `email_logs`

Audit trail for every outbound email attempt. Records persist even when SMTP delivery fails.

| Column          | Data Type                | Constraints                           | Notes                                       |
| --------------- | ------------------------ | ------------------------------------- | ------------------------------------------- |
| `id`            | INTEGER                  | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key                               |
| `recipient`     | VARCHAR(255)             | NOT NULL                              | Destination email address                   |
| `subject`       | VARCHAR(500)             | NOT NULL                              | Email subject line                          |
| `body`          | TEXT                     | NOT NULL                              | Full email body content                     |
| `status`        | VARCHAR(50)              | NOT NULL, DEFAULT `'pending'`         | Allowed values: `pending`, `sent`, `failed` |
| `created_at`    | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()`             | Set by the database on insert               |
| `sent_at`       | TIMESTAMP WITH TIME ZONE | NULLABLE                              | Populated when delivery succeeds            |
| `error_message` | TEXT                     | NULLABLE                              | Populated when delivery fails               |

---

### Table: `files`

Stores metadata about uploaded files.

| Column         | Data Type                | Constraints                           | Notes                                                       |
| -------------- | ------------------------ | ------------------------------------- | ----------------------------------------------------------- |
| `id`           | INTEGER                  | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key                                               |
| `name`         | VARCHAR(255)             | NOT NULL                              | Original file name (without path)                           |
| `extension`    | VARCHAR(20)              | NOT NULL                              | File extension, without the leading dot (e.g. `pdf`, `png`) |
| `size`         | BIGINT                   | NOT NULL                              | File size in bytes                                          |
| `content_type` | VARCHAR(100)             | NULLABLE                              | MIME type of the file                                       |
| `location`     | VARCHAR(500)             | NOT NULL                              | Storage path or URL where the file content is stored        |
| `created_at`   | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()`             | Set by the database on insert                               |
| `updated_at`   | TIMESTAMP WITH TIME ZONE | NULLABLE                              | Set by the database on update                               |

---

### Table: `users`

Stores application users for authentication and profile-management features.

| Column          | Data Type                | Constraints                           | Notes                         |
| --------------- | ------------------------ | ------------------------------------- | ----------------------------- |
| `id`            | INTEGER                  | PRIMARY KEY, AUTO INCREMENT, NOT NULL | Surrogate key                 |
| `username`      | VARCHAR(50)              | NOT NULL, UNIQUE                      | Unique login or display name  |
| `email`         | VARCHAR(255)             | NOT NULL, UNIQUE                      | Email address                 |
| `password_hash` | VARCHAR(255)             | NOT NULL                              | Hashed password value         |
| `is_active`     | BOOLEAN                  | NOT NULL, DEFAULT `true`              | Account active status         |
| `created_at`    | TIMESTAMP WITH TIME ZONE | NOT NULL, DEFAULT `now()`             | Set by the database on insert |
| `updated_at`    | TIMESTAMP WITH TIME ZONE | NULLABLE                              | Set by the database on update |

---

## Naming Convention Mapping

| Canonical (SQL)         | Python / SQLAlchemy     | .NET / EF Core        | Go                      | Java / JPA              | Node.js / Prisma      | PHP / Eloquent          |
| ----------------------- | ----------------------- | --------------------- | ----------------------- | ----------------------- | --------------------- | ----------------------- |
| `todo_items`            | `todo_items`            | `TodoItems`           | `todo_items`            | `todo_items`            | `todo_items`          | `todo_items`            |
| `todo_item_attachments` | `todo_item_attachments` | `TodoItemAttachments` | `todo_item_attachments` | `todo_item_attachments` | `todoItemAttachments` | `todo_item_attachments` |
| `todo_item_id`          | `todo_item_id`          | `TodoItemId`          | `todo_item_id`          | `todo_item_id`          | `todoItemId`          | `todo_item_id`          |
| `file_id`               | `file_id`               | `FileId`              | `file_id`               | `file_id`               | `fileId`              | `file_id`               |
| `email_logs`            | `email_logs`            | `EmailLogs`           | `email_logs`            | `email_logs`            | `email_logs`          | `email_logs`            |
| `files`                 | `files`                 | `Files`               | `files`                 | `files`                 | `files`               | `files`                 |
| `users`                 | `users`                 | `Users`               | `users`                 | `users`                 | `users`               | `users`                 |
| `is_completed`          | `is_completed`          | `IsCompleted`         | `is_completed`          | `is_completed`          | `isCompleted`         | `is_completed`          |
| `created_at`            | `created_at`            | `CreatedAt`           | `created_at`            | `created_at`            | `createdAt`           | `created_at`            |
| `updated_at`            | `updated_at`            | `UpdatedAt`           | `updated_at`            | `updated_at`            | `updatedAt`           | `updated_at`            |
| `sent_at`               | `sent_at`               | `SentAt`              | `sent_at`               | `sent_at`               | `sentAt`              | `sent_at`               |
| `error_message`         | `error_message`         | `ErrorMessage`        | `error_message`         | `error_message`         | `errorMessage`        | `error_message`         |
| `content_type`          | `content_type`          | `ContentType`         | `content_type`          | `content_type`          | `contentType`         | `content_type`          |
| `password_hash`         | `password_hash`         | `PasswordHash`        | `password_hash`         | `password_hash`         | `passwordHash`        | `password_hash`         |
| `is_active`             | `is_active`             | `IsActive`            | `is_active`             | `is_active`             | `isActive`            | `is_active`             |

## API Endpoints

All implementations expose the same REST endpoints under the `/api/todo-items` prefix for todo operations and `/api/files` for file operations. Path parameter syntax may differ by framework (`{id}` vs `:id`), but the routes are functionally identical.

### Todo Items Endpoints

| Method   | Path                            | Description                                                  |
| -------- | ------------------------------- | ------------------------------------------------------------ |
| `GET`    | `/api/todo-items`               | List all todo items (paginated)                              |
| `GET`    | `/api/todo-items/incomplete`    | List incomplete items (paginated)                            |
| `GET`    | `/api/todo-items/{id}`          | Get a single todo item                                       |
| `POST`   | `/api/todo-items`               | Create a todo item                                           |
| `PUT`    | `/api/todo-items/{id}`          | Update a todo item                                           |
| `PATCH`  | `/api/todo-items/{id}/complete` | Mark a todo item as complete                                 |
| `DELETE` | `/api/todo-items/{id}`          | Delete a todo item                                           |
| `POST`   | `/api/todo-items/import/csv`    | Import todo items from a CSV file (`multipart/form-data`)    |
| `POST`   | `/api/todo-items/import/excel`  | Import todo items from an Excel file (`multipart/form-data`) |
| `GET`    | `/api/todo-items/export/csv`    | Export todo items as a CSV file                              |
| `GET`    | `/api/todo-items/export/excel`  | Export todo items as an Excel file                           |

### Todo Item Attachments Endpoints

These endpoints manage attachment references for a specific todo item. They do not duplicate file content; they only link the todo item to an existing file stored in the `files` table.

| Method   | Path                                              | Description                                             |
| -------- | ------------------------------------------------- | ------------------------------------------------------- |
| `GET`    | `/api/todo-items/{id}/attachments`                | List all attachment references for a specific todo item |
| `POST`   | `/api/todo-items/{id}/attachments`                | Create a new attachment reference for a todo item       |
| `GET`    | `/api/todo-items/{id}/attachments/{attachmentId}` | Get a single attachment reference                       |
| `PUT`    | `/api/todo-items/{id}/attachments/{attachmentId}` | Update an attachment reference                          |
| `DELETE` | `/api/todo-items/{id}/attachments/{attachmentId}` | Remove an attachment reference from a todo item         |

### Files Endpoints

| Method   | Path                       | Description                           |
| -------- | -------------------------- | ------------------------------------- |
| `GET`    | `/api/files`               | List all uploaded files (paginated)   |
| `GET`    | `/api/files/{id}`          | Get a single file's metadata          |
| `GET`    | `/api/files/{id}/download` | Download a file's content             |
| `POST`   | `/api/files`               | Upload a file (`multipart/form-data`) |
| `DELETE` | `/api/files/{id}`          | Delete a file                         |

### Users Endpoints

#### User Management

| Method  | Path                         | Description               |
| ------- | ---------------------------- | ------------------------- |
| `GET`   | `/api/users`                 | List users (paginated)    |
| `GET`   | `/api/users/{id}`            | Get a single user by id   |
| `POST`  | `/api/users`                 | Create a new user         |
| `PUT`   | `/api/users/{id}`            | Update a user's details   |
| `PATCH` | `/api/users/{id}/activate`   | Activate a user account   |
| `PATCH` | `/api/users/{id}/deactivate` | Deactivate a user account |

#### User Self Management

| Method | Path                          | Description                             |
| ------ | ----------------------------- | --------------------------------------- |
| `POST` | `/api/users/signup`           | Register a new account                  |
| `POST` | `/api/users/password/change`  | Change the current user's password      |
| `POST` | `/api/users/password/reset`   | Request a password reset email          |
| `POST` | `/api/users/password/confirm` | Confirm a password reset with a token   |
| `GET`  | `/api/users/profile`          | Read the authenticated user's profile   |
| `PUT`  | `/api/users/profile`          | Update the authenticated user's profile |

### Tokens Endpoint

| Method | Path          | Description                                             |
| ------ | ------------- | ------------------------------------------------------- |
| `POST` | `/api/tokens` | Authenticate an active user and generate a JWT token    |

#### Generate Token

`POST /api/tokens` accepts a JSON request containing the user's email and
password:

```json
{
  "email": "alice@example.com",
  "password": "password123"
}
```

The email is trimmed and compared case-insensitively. A successful request
returns `200 OK` with a short-lived bearer token:

```json
{
  "access_token": "<jwt>",
  "token_type": "Bearer",
  "expires_in": 3600
}
```

The response must include `Cache-Control: no-store` and `Pragma: no-cache`.
The JWT uses the `HS256` algorithm and contains the following claims:

| Claim | Description                                      |
| ----- | ------------------------------------------------ |
| `sub` | User id encoded as a string                      |
| `iat` | Token issue time as a Unix timestamp             |
| `exp` | Token expiration time as a Unix timestamp        |

The token is signed with `JWT_SECRET_KEY` and its lifetime is configured by
`JWT_TOKEN_LIFETIME_MINUTES`, which defaults to `60` minutes. Tokens are only
issued when the password is correct and the user is active.

Unknown emails, incorrect passwords, and inactive accounts all return
`401 Unauthorized` with the same response so the endpoint does not disclose
account state:

```json
{
  "error": "Invalid email or password."
}
```

The `401` response must include `WWW-Authenticate: Bearer`. Malformed JSON or
an invalid request shape returns `400 Bad Request`. Issued access tokens are
not persisted in the database. Refresh tokens, revocation, logout, MFA,
scopes, and roles are outside this endpoint's initial scope.

JWT signing and bearer-token validation must use the standard maintained
library for each stack rather than application-owned JWT parsing or HMAC code:

| Implementation | JWT/authentication library                                      |
| -------------- | --------------------------------------------------------------- |
| Python         | PyJWT                                                           |
| .NET           | ASP.NET Core JwtBearer + Microsoft IdentityModel                |
| Go             | `github.com/golang-jwt/jwt/v5`                                  |
| Java           | Spring Security OAuth2 Resource Server + Nimbus JOSE JWT         |
| Node.js        | `@nestjs/jwt` + Passport JWT                                    |
| PHP            | `firebase/php-jwt`                                              |

### Pagination Query Parameters

All paginated endpoints (`GET /api/todo-items`, `GET /api/todo-items/incomplete`, and `GET /api/files`) accept the following query parameters. The parameter name casing may differ per implementation but must convey the same semantics.

| Parameter (canonical)    | Default | Description           | Implementation note                                                |
| ------------------------ | ------- | --------------------- | ------------------------------------------------------------------ |
| `page`                   | `1`     | Page number (1-based) | Same across all implementations                                    |
| `page_size` / `pageSize` | `20`    | Items per page        | `page_size` in Python & PHP; `pageSize` in .NET, Go, Java, Node.js |

---

## Implementations

| Language / Framework                                  | Path      |
| ----------------------------------------------------- | --------- |
| Python (FastAPI + SQLAlchemy + JWT/Auth)              | `python/` |
| .NET (ASP.NET Core + EF Core + ASP.NET Core Identity) | `dotnet/` |
| Go (net/http + GORM + JWT/Auth)                       | `go/`     |
| Java (Spring Boot + JPA + Spring Security)            | `java/`   |
| Node.js (NestJS + Prisma + Passport/JWT)              | `nodejs/` |
| PHP (Laravel + Eloquent + Firebase PHP-JWT)           | `php/`    |
