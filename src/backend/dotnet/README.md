# Todo API — ASP.NET Core / Entity Framework Core

A RESTful API for managing todo items built with **ASP.NET Core**, **Entity Framework Core**, and **SQLite** (swappable to SQL Server or PostgreSQL).

## Tech stack

| Concern | Technology |
|---|---|
| Web framework | ASP.NET Core 10 Web API |
| ORM | Entity Framework Core 10 |
| Database (default) | SQLite (via `Microsoft.EntityFrameworkCore.Sqlite`) |
| Migrations | EF Core Migrations (`dotnet ef`) |
| API docs | Built-in OpenAPI + **Scalar** UI |
| Dependency Injection | Built-in ASP.NET Core DI |
| Unit testing | xUnit + Moq |

## Project structure

```
src/backend/dotnet/
├── TodoApi.slnx
├── TodoApi/
│   ├── Program.cs                          # App bootstrap & DI registration
│   ├── appsettings.json                    # Connection string & logging
│   ├── Controllers/
│   │   └── TodoItemsController.cs          # REST endpoints
│   ├── Data/
│   │   ├── AppDbContext.cs                 # EF Core DbContext
│   │   └── Migrations/                     # EF Core migrations
│   ├── DTOs/
│   │   └── TodoItemDtos.cs                 # Request / response models
│   ├── Models/
│   │   └── TodoItem.cs                     # EF Core entity
│   ├── Repositories/
│   │   ├── IRepository.cs                  # Generic IRepository<T>
│   │   ├── BaseRepository.cs               # Generic BaseRepository<T>
│   │   ├── ITodoItemRepository.cs
│   │   └── TodoItemRepository.cs
│   └── Services/
│       ├── ITodoItemService.cs
│       └── TodoItemService.cs
└── TodoApi.Tests/
    ├── Controllers/
    │   └── TodoItemsControllerTests.cs     # Controller unit tests
    └── Services/
        └── TodoItemServiceTests.cs         # Service unit tests
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

### 4. Run unit tests

```bash
# Run all tests
dotnet test

# Run with detailed output
dotnet test --verbosity normal

# Run with code coverage
dotnet test --collect:"XPlat Code Coverage"
```

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
| `pageSize` | `20` | Items per page |

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

| Database | Package | Connection string |
|----------|---------|-------------------|
| SQLite (default) | `Microsoft.EntityFrameworkCore.Sqlite` | `Data Source=todo.db` |
| SQL Server | `Microsoft.EntityFrameworkCore.SqlServer` | `Server=.;Database=TodoDb;Trusted_Connection=True` |
| PostgreSQL | `Npgsql.EntityFrameworkCore.PostgreSQL` | `Host=localhost;Database=todo_db;Username=user;Password=pass` |

## Docker

### Build the image

```bash
# Run from src/backend/dotnet/
docker build -t todo-api-dotnet .
```

### Run the container

```bash
docker run -d -p 8080:8080 --name todo-api-dotnet todo-api-dotnet
```

The API is available at <http://localhost:8080>.  
Scalar API reference UI: <http://localhost:8080/scalar/v1>  
OpenAPI JSON: <http://localhost:8080/openapi/v1.json>

### Persist the SQLite database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8080:8080 -v todo-dotnet-data:/app --name todo-api-dotnet todo-api-dotnet
```

### Stop and remove the container

```bash
docker stop todo-api-dotnet
docker rm todo-api-dotnet
```
