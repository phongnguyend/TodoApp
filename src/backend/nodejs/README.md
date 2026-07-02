# Todo API - Node.js / NestJS + Prisma

A RESTful API for managing todo items built with **NestJS**, **Prisma**, and **TypeScript** - the Node.js equivalent of an ASP.NET Core + Entity Framework project.

## Tech-stack mapping

| ASP.NET Core + EF | Node.js equivalent |
|---|---|
| ASP.NET Core | **NestJS** + Express |
| Entity Framework Core | **Prisma** (ORM) |
| EF Migrations (`dotnet ef migrations`) | `prisma migrate dev` |
| `DbContext` | `PrismaService` (wraps `PrismaClient`) |
| Controllers | NestJS **Controllers** |
| Services (`IService` / `Service`) | NestJS **Services** |
| Repository pattern | `TodoItemRepository` provider |
| DTOs / Data Annotations | **class-validator** + **class-transformer** |
| `appsettings.json` / `IConfiguration` | `@nestjs/config` + `.env` |
| Dependency Injection | NestJS DI container |
| Swagger / OpenAPI | `@nestjs/swagger` at `/swagger` |
| Startup.cs / Program.cs | `main.ts` + `AppModule` |
| xUnit / Moq | **Jest** + `@nestjs/testing` |

## Project structure

```
src/backend/nodejs/
├── prisma/
│   └── schema.prisma                  # Prisma schema (EF model + DbContext)
├── src/
│   ├── api/
│   │   ├── Dockerfile                     # API container image
│   │   ├── main.ts                        # API bootstrap (Program.cs)
│   │   ├── app.module.ts                  # Root module (Startup.cs)
│   │   ├── todo-items/
│   │   │   ├── todo-items.module.ts           # Feature module
│   │   │   ├── todo-items.controller.ts       # Controller (route handlers)
│   │   │   ├── todo-items.controller.spec.ts  # Controller unit tests
│   │   │   ├── todo-items.service.ts          # Business logic
│   │   │   ├── todo-items.service.spec.ts     # Service unit tests
│   │   │   ├── todo-items.repository.ts       # Data access (Prisma queries)
│   │   │   ├── todo-items.repository.spec.ts  # Repository unit tests
│   │   │   └── dto/
│   │   │       ├── create-todo-item.dto.ts
│   │   │       ├── update-todo-item.dto.ts
│   │   │       └── todo-item-response.dto.ts
│   │   └── files/
│   │       ├── files.module.ts            # Feature module
│   │       ├── files.controller.ts        # Controller (list/get/download/upload/delete)
│   │       ├── files.controller.spec.ts   # Controller unit tests
│   │       ├── files.service.ts           # Business logic (storage on disk + metadata)
│   │       ├── files.service.spec.ts      # Service unit tests
│   │       ├── files.repository.ts        # Data access (Prisma queries)
│   │       ├── files.repository.spec.ts   # Repository unit tests
│   │       └── dto/
│   │           └── file-response.dto.ts
│   ├── shared/
│   │   ├── prisma/
│   │   │   ├── prisma.service.ts          # PrismaClient wrapper (DbContext)
│   │   │   └── prisma.module.ts           # Global Prisma module
│   │   └── common/
│   │       └── dto/
│   │           └── paginated-response.dto.ts
│   └── worker/
│       ├── Dockerfile                     # Background worker container image
│       ├── main.ts                        # Worker entry-point (plain Node.js process)
│       └── jobs/
│           └── incomplete-todos-email.job.ts  # Email digest job
├── package.json
├── tsconfig.json
├── nest-cli.json
└── .env.example
```

## Getting started

### 1. Install dependencies

```bash
cd src/backend/nodejs
npm install
```

### 2. Configure environment

```bash
copy .env.example .env
# Edit .env as needed
```

### 3. Generate Prisma client & run migrations (like `dotnet ef database update`)

```bash
# Generate a migration (like `dotnet ef migrations add InitialCreate`)
npx prisma migrate dev --name InitialCreate

# Prisma client is generated automatically after migrate dev
# To generate manually:
npx prisma generate
```

### 4. Run unit tests

```bash
# Run all unit tests
npm test

# Run tests in watch mode
npm run test:watch

# Run tests with coverage report
npm run test:cov
```

### 5. Start the development server

```bash
npm run start:dev
```

The Swagger UI is available at <http://localhost:3000/swagger>.

### 6. Run the background worker locally

```bash
# Build first (compiles worker/main.ts → dist/worker/main.js)
npm run build

# Then start the worker process
npm run start:worker
```

The worker reads SMTP settings and `WORKER_INTERVAL_MINUTES` from `.env`.

The Swagger UI is available at <http://localhost:3000/swagger>.

## API endpoints

See the [shared API contract](../README.md#api-endpoints) in the backend README.

### Files

The `files` feature (`src/api/files/`) stores uploaded file content on disk and its metadata in the `files` table (via Prisma).

Uploads are handled with `@nestjs/platform-express`'s `FileInterceptor`, buffered in memory, then written to `FILE_STORAGE_PATH` under a randomly-prefixed file name (to avoid collisions/overwrites between uploads that share a name). The client-supplied file name is stripped of any directory components before being stored, to prevent path traversal.

| Variable | Default | Description |
|---|---|---|
| `FILE_STORAGE_PATH` | `./uploads` | Directory where uploaded file content is stored |
| `MAX_UPLOAD_SIZE_BYTES` | `10485760` (10 MB) | Maximum accepted upload size; larger files are rejected with `413 Payload Too Large` |

## Background worker

The worker runs as a **separate process / container** (`src/worker/Dockerfile`). It is intentionally kept outside the NestJS application context and connects to Prisma directly.

### What it does

On startup and then every `WORKER_INTERVAL_MINUTES` minutes (default: 60):

1. Queries all incomplete todo items (`isCompleted = false`).
2. Builds a plain-text + HTML email digest.
3. Inserts an `email_logs` row with `status = "pending"`.
4. Sends the email via SMTP (nodemailer - supports STARTTLS and SSL).
5. Updates the `email_logs` row to `status = "sent"` or `"failed"` (with `errorMessage`).

If there are no incomplete todos, the job is skipped and no email is sent.

### Worker environment variables

| Variable | Default | Description |
|---|---|---|
| `DATABASE_URL` | *(required)* | Same value as the API - both containers share the same database |
| `SMTP_HOST` | `localhost` | SMTP server hostname |
| `SMTP_PORT` | `587` | SMTP server port |
| `SMTP_SECURE` | `false` | `true` → implicit TLS (port 465); `false` → STARTTLS (port 587) |
| `SMTP_USERNAME` | *(empty)* | SMTP auth username |
| `SMTP_PASSWORD` | *(empty)* | SMTP auth password |
| `EMAIL_SENDER` | `noreply@example.com` | From address |
| `EMAIL_RECIPIENT` | `admin@example.com` | Destination address for digests |
| `WORKER_INTERVAL_MINUTES` | `60` | How often the job runs |

## Switching databases

Update `DATABASE_URL` in `.env` and change `provider` in `prisma/schema.prisma`:

- **PostgreSQL**: `"postgresql://user:password@localhost:5432/todo_db?schema=public"`
- **MySQL**: `"mysql://user:password@localhost:3306/todo_db"`

Then re-run:

```bash
npx prisma migrate dev
```

## Docker

### Build and run the API

```bash
# Build the API image
docker build -f src/api/Dockerfile -t todo-api-nodejs .

# Run the API container
docker run -d -p 3000:3000 \
  -e DATABASE_URL="file:/app/data/todo.db" \
  -v todo-nodejs-data:/app/data \
  --name todo-api-nodejs todo-api-nodejs
```

The API is available at <http://localhost:3000>.  
Swagger UI: <http://localhost:3000/swagger>

### Build and run the background worker

```bash
# Build the worker image
docker build -f src/worker/Dockerfile -t todo-worker-nodejs .

# Run the worker container (shares the same database volume as the API)
docker run -d \
  -e DATABASE_URL="file:/app/data/todo.db" \
  -e SMTP_HOST=smtp.example.com \
  -e SMTP_PORT=587 \
  -e SMTP_SECURE=false \
  -e SMTP_USERNAME=user@example.com \
  -e SMTP_PASSWORD=secret \
  -e EMAIL_SENDER=noreply@example.com \
  -e EMAIL_RECIPIENT=admin@example.com \
  -e WORKER_INTERVAL_MINUTES=60 \
  -v todo-nodejs-data:/app/data \
  --name todo-worker-nodejs todo-worker-nodejs
```

### Stop and remove containers

```bash
docker stop todo-api-nodejs todo-worker-nodejs
docker rm  todo-api-nodejs todo-worker-nodejs
```
