# Todo API — PHP / Laravel

A RESTful API for managing todo items built with **Laravel 11**, **Eloquent ORM**, and Laravel's built-in migration system — the PHP equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF | PHP equivalent |
|---|---|
| ASP.NET Core | **Laravel 11** |
| Entity Framework Core | **Eloquent ORM** |
| EF Migrations (`dotnet ef migrations add`) | **Laravel Migrations** (`php artisan make:migration`) |
| `DbContext` | Eloquent `Model` / `DB` facade |
| Controllers (`ControllerBase`) | **Laravel Controllers** |
| Services (`IService` / `Service`) | `TodoItemServiceInterface` / `TodoItemService` |
| Repository pattern (`IRepository<T>`) | `RepositoryInterface` / `BaseRepository` |
| DTOs / Request validators | **Form Requests** + **API Resources** |
| `appsettings.json` / `IConfiguration` | `.env` + `config/` |
| Dependency Injection (constructor DI) | Laravel **Service Container** (`AppServiceProvider`) |
| Swagger / OpenAPI | **L5-Swagger** (`darkaonline/l5-swagger`) |
| Global exception handler / ProblemDetails | `app/Exceptions/Handler.php` |

## Project structure

```
src/backend/php/
├── app/
│   ├── Exceptions/
│   │   └── Handler.php                        # Global error handler (ProblemDetails equivalent)
│   ├── Http/
│   │   ├── Controllers/Api/
│   │   │   └── TodoItemController.php          # REST controller (ControllerBase)
│   │   ├── Requests/
│   │   │   ├── CreateTodoItemRequest.php        # Validated create DTO
│   │   │   └── UpdateTodoItemRequest.php        # Validated update DTO
│   │   └── Resources/
│   │       └── TodoItemResource.php             # Response DTO (AutoMapper profile)
│   ├── Models/
│   │   └── TodoItem.php                         # Eloquent entity
│   ├── Providers/
│   │   └── AppServiceProvider.php               # IoC bindings (Program.cs AddScoped)
│   ├── Repositories/
│   │   ├── Contracts/
│   │   │   ├── RepositoryInterface.php          # IRepository<T>
│   │   │   └── TodoItemRepositoryInterface.php  # ITodoItemRepository
│   │   ├── BaseRepository.php                   # GenericRepository<T>
│   │   └── TodoItemRepository.php               # Concrete implementation
│   └── Services/
│       ├── Contracts/
│       │   └── TodoItemServiceInterface.php     # ITodoItemService
│       └── TodoItemService.php                  # Business logic
├── database/
│   └── migrations/
│       └── 2024_01_01_000000_create_todo_items_table.php
├── routes/
│   └── api.php                                  # Route definitions
├── composer.json
└── .env.example
```

## Getting started

### 1. Install dependencies

```bash
cd src/backend/php
composer install
```

### 2. Configure environment

```bash
cp .env.example .env
php artisan key:generate
```

Edit `.env` — by default it uses **SQLite**. For SQLite, create the database file:

```bash
touch database/database.sqlite
```

### 3. Run migrations (like `dotnet ef database update`)

```bash
php artisan migrate
```

To create a new migration (like `dotnet ef migrations add`):

```bash
php artisan make:migration add_priority_to_todo_items --table=todo_items
```

### 4. Start the development server

```bash
php artisan serve --port=8000
```

### 5. Browse the Swagger UI

Install L5-Swagger assets and generate docs:

```bash
php artisan vendor:publish --provider="L5Swagger\L5SwaggerServiceProvider"
php artisan l5-swagger:generate
```

Swagger UI → <http://localhost:8000/api/documentation>

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
| `page_size` | `20` | Items per page |

## Switching databases

Update `DB_CONNECTION` and related variables in `.env`:

- **MySQL**: `DB_CONNECTION=mysql` + install nothing (included in Laravel)
- **PostgreSQL**: `DB_CONNECTION=pgsql` + `composer require doctrine/dbal`
- **SQLite** (default, dev only): `DB_CONNECTION=sqlite`
