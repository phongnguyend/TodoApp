# Todo API - PHP / Laravel

A RESTful API for managing todo items built with **Laravel 12**, **Eloquent ORM**, and Laravel's built-in migration system - the PHP equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF                          | PHP equivalent                                                |
| ------------------------------------------ | ------------------------------------------------------------- |
| ASP.NET Core                               | **Laravel 12**                                                |
| Entity Framework Core                      | **Eloquent ORM**                                              |
| EF Migrations (`dotnet ef migrations add`) | **Laravel Migrations** (`php artisan make:migration`)         |
| `DbContext`                                | Eloquent `Model` / `DB` facade                                |
| Controllers (`ControllerBase`)             | **Laravel Controllers**                                       |
| Services (`IService` / `Service`)          | `TodoItemServiceInterface` / `TodoItemService`                |
| Repository pattern (`IRepository<T>`)      | `RepositoryInterface` / `BaseRepository`                      |
| DTOs / Request validators                  | **Form Requests** + **API Resources**                         |
| `appsettings.json` / `IConfiguration`      | `.env` + `config/`                                            |
| Dependency Injection (constructor DI)      | Laravel **Service Container** (`AppServiceProvider`)          |
| Swagger / OpenAPI                          | **L5-Swagger** (`darkaonline/l5-swagger`)                     |
| Global exception handler / ProblemDetails  | `app/Exceptions/Handler.php`                                  |
| Unit tests (xUnit / NUnit)                 | **PHPUnit 11** + **Mockery**                                  |
| Background / hosted service                | **Artisan command** + **Laravel Scheduler** (`schedule:work`) |

## Project structure

```
src/backend/php/
├── app/
│   ├── Console/
│   │   └── Commands/
│   │       └── ProcessIncompleteRemindersCommand.php  # Background job (Artisan command)
│   ├── Exceptions/
│   │   └── Handler.php                        # Global error handler (ProblemDetails equivalent)
│   ├── Http/
│   │   ├── Controllers/Api/
│   │   │   ├── TodoItemController.php          # REST controller (ControllerBase)
│   │   │   ├── TodoItemAttachmentController.php # REST controller for attachment references
│   │   │   └── FileController.php               # REST controller for uploaded files
│   │   ├── Requests/
│   │   │   ├── CreateTodoItemRequest.php        # Validated create DTO
│   │   │   ├── UpdateTodoItemRequest.php        # Validated update DTO
│   │   │   ├── SaveTodoItemAttachmentRequest.php # Validated attachment-reference DTO
│   │   │   └── UploadFileRequest.php            # Validated file-upload DTO (required/size rules)
│   │   └── Resources/
│   │       ├── TodoItemResource.php             # Response DTO (AutoMapper profile)
│   │       ├── TodoItemAttachmentResource.php   # Response DTO for attachment references
│   │       └── FileResource.php                 # Response DTO for file metadata
│   ├── Models/
│   │   ├── TodoItem.php                         # Eloquent entity
│   │   ├── TodoItemAttachment.php               # Eloquent attachment-reference entity
│   │   ├── EmailLog.php                         # Eloquent entity for email audit trail
│   │   └── File.php                             # Eloquent entity for uploaded-file metadata
│   ├── Providers/
│   │   └── AppServiceProvider.php               # IoC bindings (Program.cs AddScoped)
│   ├── Repositories/
│   │   ├── Contracts/
│   │   │   ├── RepositoryInterface.php          # IRepository<T>
│   │   │   ├── TodoItemRepositoryInterface.php  # ITodoItemRepository
│   │   │   ├── TodoItemAttachmentRepositoryInterface.php # Attachment repository contract
│   │   │   └── FileRepositoryInterface.php      # IFileRepository
│   │   ├── BaseRepository.php                   # GenericRepository<T>
│   │   ├── TodoItemRepository.php               # Concrete implementation
│   │   ├── TodoItemAttachmentRepository.php     # Attachment repository implementation
│   │   └── FileRepository.php                   # Concrete implementation
│   └── Services/
│       ├── Contracts/
│       │   ├── TodoItemServiceInterface.php     # ITodoItemService
│       │   ├── TodoItemAttachmentServiceInterface.php # Attachment service contract
│       │   └── FileServiceInterface.php         # IFileService
│       ├── TodoItemService.php                  # Business logic
│       ├── TodoItemAttachmentService.php        # Attachment-reference business logic
│       └── FileService.php                      # Business logic - upload/download/delete on disk
├── database/
│   ├── factories/
│   │   ├── TodoItemFactory.php                  # Model factory for tests
│   │   └── FileFactory.php                      # Model factory for tests
│   └── migrations/
│       ├── 2024_01_01_000000_create_todo_items_table.php
│       ├── 2024_01_02_000000_create_email_logs_table.php
│       ├── 2024_01_03_000000_create_files_table.php
│       └── 2026_07_18_000000_create_todo_item_attachments_table.php
├── routes/
│   └── api.php                                  # Route definitions
├── tests/
│   ├── TestCase.php                             # Base test case
│   ├── Unit/
│   │   ├── Requests/
│   │   │   ├── CreateTodoItemRequestTest.php    # Validation rule tests
│   │   │   ├── UpdateTodoItemRequestTest.php    # Validation rule tests
│   │   │   └── UploadFileRequestTest.php        # Validation rule tests (required/max size)
│   │   └── Services/
│   │       ├── TodoItemServiceTest.php          # Service unit tests (Mockery)
│   │       ├── TodoItemAttachmentServiceTest.php # Attachment service unit tests
│   │       └── FileServiceTest.php              # Service unit tests (Mockery + real temp files)
│   └── Feature/
│       ├── TodoItemApiTest.php                  # HTTP integration tests
│       ├── TodoItemAttachmentApiTest.php        # Attachment endpoint integration tests
│       └── FileApiTest.php                      # HTTP integration tests (upload/download/delete)
├── phpunit.xml
├── composer.json
├── build/
│   ├── api/
│   │   └── Dockerfile            # API container image
│   └── worker/
│       └── Dockerfile            # Background worker container image
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

Edit `.env` - by default it uses **SQLite**. For SQLite, create the database file:

```bash
touch database/database.sqlite
```

File uploads are written to the path set by `FILE_STORAGE_PATH` (defaults to `storage/app/uploads` when unset), and are capped by `MAX_UPLOAD_SIZE_BYTES` (default `10485760`, i.e. 10 MB).

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

See the [shared API contract](../README.md#api-endpoints) in the backend README.

## File uploads

The `/api/files` endpoints (list, get metadata, download, upload, delete) store uploaded file content on disk and persist metadata (`name`, `extension`, `size`, `content_type`, `location`, timestamps) in the `files` table.

### Configuration (`.env`)

```ini
FILE_STORAGE_PATH=/absolute/path/to/storage/app/uploads  # Directory where uploaded file content is stored (defaults to storage/app/uploads when unset)
MAX_UPLOAD_SIZE_BYTES=10485760                            # Maximum accepted upload size, in bytes (default 10 MB)
```

### Notes

- Uploaded file names are sanitized (directory components stripped via `basename()`) and stored on disk under a random-prefixed name to prevent path traversal and filename collisions.
- The internal storage `location` is never exposed in API responses; file content is retrieved via `GET /api/files/{id}/download`.
- Deleting a file removes both the database row and the file content on disk.
- Uploads exceeding `MAX_UPLOAD_SIZE_BYTES` are rejected by `UploadFileRequest`'s validation rules with `422 Unprocessable Entity` (Laravel's standard validation-failure response).

## Switching databases

Update `DB_CONNECTION` and related variables in `.env`:

- **MySQL**: `DB_CONNECTION=mysql` + install nothing (included in Laravel)
- **PostgreSQL**: `DB_CONNECTION=pgsql` + `composer require doctrine/dbal`
- **SQLite** (default, dev only): `DB_CONNECTION=sqlite`

## Running tests

The test suite uses **PHPUnit 11** with an in-memory SQLite database - no extra setup required.

> **Missing PHP extensions - Windows gotcha**  
> The in-memory SQLite test database requires the `pdo_sqlite` and `sqlite3` extensions.  
> On a default Windows PHP installation both extensions are present but **commented out** in `php.ini`.  
> Open your active `php.ini` (run `php --ini` to find it) and uncomment the two lines:
>
> ```ini
> extension=pdo_sqlite
> extension=sqlite3
> ```
>
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

| Suite   | File                        | What is tested                                                                                 |
| ------- | --------------------------- | ---------------------------------------------------------------------------------------------- |
| Unit    | `TodoItemServiceTest`       | All service methods; repository mocked via Mockery                                             |
| Unit    | `CreateTodoItemRequestTest` | `title` required/length, `description` optional/length                                         |
| Unit    | `UpdateTodoItemRequestTest` | All fields optional, type/length constraints                                                   |
| Unit    | `FileServiceTest`           | All service methods; repository mocked via Mockery, real temp files for upload/download/delete |
| Unit    | `UploadFileRequestTest`     | `file` required, rejected when exceeding `MAX_UPLOAD_SIZE_BYTES`                               |
| Feature | `TodoItemApiTest`           | All 7 todo-item endpoints - status codes, response shape, database state                       |
| Feature | `FileApiTest`               | All 5 file endpoints - upload/download/delete round-trip against a temp storage dir            |

## Docker

### Build the image

```bash
# Run from src/backend/php/
docker build -f build/api/Dockerfile -t todo-api-php .
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

## Background worker

A separate container runs the **incomplete-todo reminder** job. It uses the Laravel scheduler (`schedule:work`) to fire the `app:process-incomplete-reminders` Artisan command **every hour**.

### What the worker does

1. Queries all incomplete todo items (`is_completed = false`).
2. Builds a plain-text email body listing each item.
3. Inserts an `email_logs` record with `status = pending`.
4. Sends the email via SMTP (`MAIL_*` env variables).
5. Updates the log to `status = sent` (and sets `sent_at`) on success, or `status = failed` (and sets `error_message`) on failure.
6. If there are no incomplete items the run is a no-op (no email sent, no log written).

### Build the worker image

```bash
# Run from src/backend/php/
docker build -f build/worker/Dockerfile -t todo-worker-php .
```

### Run the worker container

The worker must share the same database as the API. When using SQLite, mount the same volume:

```bash
docker run -d \
  -e APP_KEY=base64:$(openssl rand -base64 32) \
  -e MAIL_MAILER=smtp \
  -e MAIL_HOST=smtp.example.com \
  -e MAIL_PORT=587 \
  -e MAIL_USERNAME=user@example.com \
  -e MAIL_PASSWORD=secret \
  -e MAIL_ENCRYPTION=tls \
  -e MAIL_FROM_ADDRESS=noreply@example.com \
  -e MAIL_REMINDER_RECIPIENT=admin@example.com \
  -v todo-php-data:/var/www/html/database \
  --name todo-worker-php todo-worker-php
```

### Run the job immediately (one-shot)

```bash
php artisan app:process-incomplete-reminders
```

### View email logs

Use any SQL client (or `php artisan tinker`) to inspect the `email_logs` table:

```php
App\Models\EmailLog::latest('created_at')->get();
```
