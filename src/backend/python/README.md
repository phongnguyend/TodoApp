# Todo API - Python / FastAPI

A RESTful API for managing todo items built with **FastAPI**, **SQLAlchemy**, and **Alembic** - the Python equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF | Python equivalent |
|---|---|
| ASP.NET Core | **FastAPI** + Uvicorn |
| Entity Framework Core | **SQLAlchemy** (ORM) |
| EF Migrations (`dotnet ef migrations`) | **Alembic** |
| `DbContext` | `SessionLocal` / `Session` |
| Controllers | FastAPI **Routers** |
| Services (`IService` / `Service`) | `ITodoItemService` / `TodoItemService` |
| Repository pattern | `IRepository<T>` / `BaseRepository` |
| DTOs / Request models | **Pydantic** schemas |
| `appsettings.json` / `IConfiguration` | `pydantic-settings` + `.env` |
| Dependency Injection | FastAPI `Depends()` |
| Swagger / OpenAPI | Built-in at `/swagger` |
| Unit tests (`xUnit` / `NUnit`) | **pytest** + `unittest.mock` |

## Project structure

```
src/backend/python/
├── api/                           # API project (analogous to dotnet/TodoApi)
│   ├── main.py                        # App entry point (Program.cs)
│   ├── config.py                      # Settings (appsettings.json)
│   ├── database.py                    # DB session + Base (DbContext)
│   ├── models/
│   │   ├── todo_item.py               # SQLAlchemy entity
│   │   ├── email_log.py               # Email audit-log entity
│   │   └── file.py                    # Uploaded-file entity
│   ├── schemas/
│   │   ├── todo_item.py               # Pydantic DTOs (request/response)
│   │   └── file.py                    # File metadata DTO (response)
│   ├── repositories/
│   │   ├── base_repository.py         # IRepository<T> + BaseRepository<T>
│   │   ├── todo_item_repository.py    # ITodoItemRepository + impl
│   │   └── file_repository.py         # IFileRepository + impl
│   ├── services/
│   │   ├── todo_item_service.py       # ITodoItemService + impl
│   │   └── file_service.py            # IFileService + impl (upload/download/delete)
│   ├── routers/
│   │   ├── todo_items.py              # TodoItemsController equivalent
│   │   └── files.py                   # FilesController equivalent
│   └── Dockerfile                     # API container
├── worker/                        # Worker project (analogous to dotnet/TodoWorker)
│   ├── main.py                        # Worker entry-point (scheduler)
│   ├── config.py                      # Settings (appsettings.json)
│   ├── database.py                    # DB session + Base (DbContext)
│   ├── models/
│   │   ├── todo_item.py               # SQLAlchemy entity
│   │   ├── email_log.py               # Email audit-log entity
│   │   └── file.py                    # Uploaded-file entity
│   ├── jobs/
│   │   └── incomplete_todos_email.py  # Digest email job
│   └── Dockerfile                     # Background worker container
├── tests/
│   └── unit/
│       ├── services/
│       │   ├── test_todo_item_service.py  # Service layer unit tests
│       │   └── test_file_service.py       # File service unit tests
│       └── routers/
│           ├── test_todo_items.py         # Router / HTTP endpoint tests
│           └── test_files.py              # File router / HTTP endpoint tests
├── alembic/
│   ├── env.py
│   ├── script.py.mako
│   └── versions/
│       ├── 20260630_0000_aabbccdd1122_initial_create.py     # todo_items + email_logs
│       └── 20260702_0000_bbccddee2233_add_files_table.py    # files
├── alembic.ini
├── pytest.ini
├── requirements.txt
└── .env.example
```

## Getting started

### 1. Create and activate a virtual environment

```bash
python -m venv .venv
# Windows
.venv\Scripts\activate
# macOS / Linux
source .venv/bin/activate
```

### 2. Install dependencies

```bash
pip install -r requirements.txt
```

### 3. Configure environment

```bash
cp .env.example .env
# Edit .env as needed
```

### 4. Run database migrations (like `dotnet ef database update`)

```bash
# Apply the bundled migrations (creates todo_items, email_logs, and files)
alembic upgrade head

# To generate a new migration after model changes (like `dotnet ef migrations add`)
alembic revision --autogenerate -m "DescribeChange"
```

### 5. Run unit tests

```bash
# Run all unit tests
python -m pytest tests/ -v

# Run with summary (no verbose output)
python -m pytest tests/

# Run a specific layer
python -m pytest tests/unit/services/ -v
python -m pytest tests/unit/routers/ -v
```

### 6. Start the API server

```bash
uvicorn api.main:app --reload --host 0.0.0.0 --port 8000
```

The Swagger UI is available at <http://localhost:8000/swagger>.

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

## Switching databases

Update `DATABASE_URL` in `.env`:

- **PostgreSQL**: `postgresql://user:password@localhost:5432/todo_db`  
  Install: `pip install psycopg2-binary`
- **MySQL**: `mysql+pymysql://user:password@localhost:3306/todo_db`  
  Install: `pip install pymysql`

## Background worker

A separate process (and container) runs on a configurable interval and sends an email digest of all incomplete todo items.

### What it does

1. Queries every `TodoItem` where `is_completed = false`.
2. Builds a plain-text + HTML email listing the items.
3. Persists an `EmailLog` row (`status = pending`) to the database.
4. Delivers the email via SMTP (STARTTLS or SSL).
5. Updates the `EmailLog` row to `status = sent` (or `failed` with the error message).

### Email-log table

| Column | Type | Description |
|---|---|---|
| `id` | `INTEGER` | Primary key |
| `recipient` | `VARCHAR(255)` | Destination address |
| `subject` | `VARCHAR(500)` | Email subject |
| `body` | `TEXT` | Plain-text body |
| `status` | `VARCHAR(50)` | `pending` / `sent` / `failed` |
| `created_at` | `DATETIME` | Row creation time |
| `sent_at` | `DATETIME` | Delivery time (nullable) |
| `error_message` | `TEXT` | SMTP error details (nullable) |

### Configuration (`.env`)

```ini
SMTP_HOST=smtp.example.com
SMTP_PORT=587
SMTP_USE_TLS=true          # false → SMTP_SSL (port 465)
SMTP_USERNAME=user@example.com
SMTP_PASSWORD=changeme
EMAIL_SENDER=noreply@example.com
EMAIL_RECIPIENT=admin@example.com
WORKER_INTERVAL_MINUTES=60 # how often the digest is sent
```

### Run locally

```bash
# Make sure the .env file is configured and migrations are applied first
python -m worker.main
```

The worker runs the job immediately on startup, then repeats every `WORKER_INTERVAL_MINUTES` minutes.

## Docker

### Build the API image

```bash
# Run from src/backend/python/
docker build -f api/Dockerfile -t todo-api-python .
```

### Build the worker image

```bash
docker build -f worker/Dockerfile -t todo-worker-python .
```

### Run the API container

```bash
docker run -d -p 8000:8000 --name todo-api-python todo-api-python
```

The API is available at <http://localhost:8000>.  
Swagger UI: <http://localhost:8000/swagger>

### Run the worker container

The worker must share the same database as the API. Pass SMTP credentials via environment variables:

```bash
docker run -d \
  -e DATABASE_URL="sqlite:////app/data/todo.db" \
  -e SMTP_HOST="smtp.example.com" \
  -e SMTP_PORT="587" \
  -e SMTP_USERNAME="user@example.com" \
  -e SMTP_PASSWORD="changeme" \
  -e EMAIL_SENDER="noreply@example.com" \
  -e EMAIL_RECIPIENT="admin@example.com" \
  -e WORKER_INTERVAL_MINUTES="60" \
  -v todo-python-data:/app/data \
  --name todo-worker-python todo-worker-python
```

### Persist the SQLite database (API + worker sharing the same volume)

```bash
# API
docker run -d -p 8000:8000 \
  -e DATABASE_URL="sqlite:////app/data/todo.db" \
  -v todo-python-data:/app/data \
  --name todo-api-python todo-api-python

# Worker (mount the same named volume)
docker run -d \
  -e DATABASE_URL="sqlite:////app/data/todo.db" \
  -e SMTP_HOST="smtp.example.com" \
  -e SMTP_USERNAME="user@example.com" \
  -e SMTP_PASSWORD="changeme" \
  -e EMAIL_SENDER="noreply@example.com" \
  -e EMAIL_RECIPIENT="admin@example.com" \
  -v todo-python-data:/app/data \
  --name todo-worker-python todo-worker-python
```

### Stop and remove the containers

```bash
docker stop todo-api-python todo-worker-python
docker rm  todo-api-python todo-worker-python
```
