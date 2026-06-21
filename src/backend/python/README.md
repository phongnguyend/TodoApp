# Todo API — Python / FastAPI

A RESTful API for managing todo items built with **FastAPI**, **SQLAlchemy**, and **Alembic** — the Python equivalent of an ASP.NET Core + Entity Framework project.

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

## Project structure

```
src/backend/python/
├── app/
│   ├── main.py                        # App entry point (Program.cs)
│   ├── config.py                      # Settings (appsettings.json)
│   ├── database.py                    # DB session + Base (DbContext)
│   ├── models/
│   │   └── todo_item.py               # SQLAlchemy entity
│   ├── schemas/
│   │   └── todo_item.py               # Pydantic DTOs (request/response)
│   ├── repositories/
│   │   ├── base_repository.py         # IRepository<T> + BaseRepository<T>
│   │   └── todo_item_repository.py    # ITodoItemRepository + impl
│   ├── services/
│   │   └── todo_item_service.py       # ITodoItemService + impl
│   └── routers/
│       └── todo_items.py              # TodoItemsController equivalent
├── alembic/
│   ├── env.py
│   ├── script.py.mako
│   └── versions/
├── alembic.ini
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
# Generate a new migration (like `dotnet ef migrations add InitialCreate`)
alembic revision --autogenerate -m "InitialCreate"

# Apply migrations
alembic upgrade head
```

### 5. Start the API server

```bash
uvicorn app.main:app --reload --host 0.0.0.0 --port 8000
```

The Swagger UI is available at <http://localhost:8000/swagger>.

## API endpoints

| Method | URL | Description |
|--------|-----|-------------|
| `GET` | `/api/todo-items/` | List all todo items (paginated) |
| `GET` | `/api/todo-items/incomplete` | List incomplete items (paginated) |
| `GET` | `/api/todo-items/{id}` | Get a single todo item |
| `POST` | `/api/todo-items/` | Create a todo item |
| `PUT` | `/api/todo-items/{id}` | Update a todo item |
| `PATCH` | `/api/todo-items/{id}/complete` | Mark a todo item as complete |
| `DELETE` | `/api/todo-items/{id}` | Delete a todo item |

### Pagination query parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `page` | `1` | Page number (1-based) |
| `page_size` | `20` | Items per page |

## Switching databases

Update `DATABASE_URL` in `.env`:

- **PostgreSQL**: `postgresql://user:password@localhost:5432/todo_db`  
  Install: `pip install psycopg2-binary`
- **MySQL**: `mysql+pymysql://user:password@localhost:3306/todo_db`  
  Install: `pip install pymysql`

## Docker

### Build the image

```bash
# Run from src/backend/python/
docker build -t todo-api-python .
```

### Run the container

```bash
docker run -d -p 8000:8000 --name todo-api-python todo-api-python
```

The API is available at <http://localhost:8000>.  
Swagger UI: <http://localhost:8000/swagger>

### Persist the SQLite database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8000:8000 \
  -e DATABASE_URL="sqlite:////app/data/todo.db" \
  -v todo-python-data:/app/data \
  --name todo-api-python todo-api-python
```

### Stop and remove the container

```bash
docker stop todo-api-python
docker rm todo-api-python
```
