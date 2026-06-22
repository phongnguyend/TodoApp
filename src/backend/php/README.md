# Todo API — PHP / Laravel

A RESTful API for managing todo items built with **Laravel 12**, **Eloquent ORM**, and Laravel's built-in migration system — the PHP equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF | PHP equivalent |
|---|---|
| ASP.NET Core | **Laravel 12** |
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
| Unit tests (xUnit / NUnit) | **PHPUnit 11** + **Mockery** |

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
│   ├── factories/
│   │   └── TodoItemFactory.php                  # Model factory for tests
│   └── migrations/
│       └── 2024_01_01_000000_create_todo_items_table.php
├── routes/
│   └── api.php                                  # Route definitions
├── tests/
│   ├── TestCase.php                             # Base test case
│   ├── Unit/
│   │   ├── Requests/
│   │   │   ├── CreateTodoItemRequestTest.php    # Validation rule tests
│   │   │   └── UpdateTodoItemRequestTest.php    # Validation rule tests
│   │   └── Services/
│   │       └── TodoItemServiceTest.php          # Service unit tests (Mockery)
│   └── Feature/
│       └── TodoItemApiTest.php                  # HTTP integration tests
├── phpunit.xml
├── composer.json
└── .env.example
```

## Getting started

### Prerequisites: Install PHP & Composer

#### Windows

Using [Scoop](https://scoop.sh/):

```powershell
scoop install php composer
```

Using [Chocolatey](https://chocolatey.org/):

```powershell
choco install php composer
```

Or download manually:
- PHP: <https://windows.php.net/download/> (grab a **Thread Safe** x64 zip, extract, and add to `PATH`)
- Composer: <https://getcomposer.org/Composer-Setup.exe>

#### macOS

```bash
brew install php composer
```

#### Linux (Debian / Ubuntu)

```bash
sudo apt update
sudo apt install php php-cli php-mbstring php-xml php-sqlite3 unzip curl
curl -sS https://getcomposer.org/installer | php
sudo mv composer.phar /usr/local/bin/composer
```

#### Linux (Fedora / RHEL)

```bash
sudo dnf install php php-cli php-mbstring php-xml php-pdo
curl -sS https://getcomposer.org/installer | php
sudo mv composer.phar /usr/local/bin/composer
```

Verify the installation:

```bash
php --version
composer --version
```

---

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

## Running tests

The test suite uses **PHPUnit 11** with an in-memory SQLite database — no extra setup required.

> **Missing PHP extensions — Windows gotcha**  
> The in-memory SQLite test database requires the `pdo_sqlite` and `sqlite3` extensions.  
> On a default Windows PHP installation both extensions are present but **commented out** in `php.ini`.  
> Open your active `php.ini` (run `php --ini` to find it) and uncomment the two lines:
> ```ini
> extension=pdo_sqlite
> extension=sqlite3
> ```
> Without these, PHPUnit will fail with a `could not find driver` (PDO) error.

```bash
# All suites
./vendor/bin/phpunit

# Unit tests only
./vendor/bin/phpunit --testsuite Unit

# Feature (HTTP) tests only
./vendor/bin/phpunit --testsuite Feature
```

### Test coverage

| Suite | File | What is tested |
|---|---|---|
| Unit | `TodoItemServiceTest` | All service methods; repository mocked via Mockery |
| Unit | `CreateTodoItemRequestTest` | `title` required/length, `description` optional/length |
| Unit | `UpdateTodoItemRequestTest` | All fields optional, type/length constraints |
| Feature | `TodoItemApiTest` | All 7 endpoints — status codes, response shape, database state |

## Docker

### Build the image

```bash
# Run from src/backend/php/
docker build -t todo-api-php .
```

### Run the container

```bash
docker run -d -p 8080:8080 \
  -e APP_KEY=base64:$(openssl rand -base64 32) \
  --name todo-api-php todo-api-php
```

The API is available at <http://localhost:8080>.  
Swagger UI: <http://localhost:8080/api/documentation>

### Persist the SQLite database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 8080:8080 \
  -e APP_KEY=base64:$(openssl rand -base64 32) \
  -v todo-php-data:/var/www/html/database \
  --name todo-api-php todo-api-php
```

### Stop and remove the container

```bash
docker stop todo-api-php
docker rm todo-api-php
```
