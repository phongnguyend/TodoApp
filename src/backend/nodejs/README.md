# Todo API — Node.js / NestJS + Prisma

A RESTful API for managing todo items built with **NestJS**, **Prisma**, and **TypeScript** — the Node.js equivalent of an ASP.NET Core + Entity Framework project.

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

## Project structure

```
src/backend/nodejs/
├── prisma/
│   └── schema.prisma                  # Prisma schema (EF model + DbContext)
├── src/
│   ├── main.ts                        # Bootstrap (Program.cs)
│   ├── app.module.ts                  # Root module (Startup.cs)
│   ├── prisma/
│   │   ├── prisma.service.ts          # PrismaClient wrapper (DbContext)
│   │   └── prisma.module.ts           # Global Prisma module
│   ├── common/
│   │   └── dto/
│   │       └── paginated-response.dto.ts
│   └── todo-items/
│       ├── todo-items.module.ts       # Feature module
│       ├── todo-items.controller.ts   # Controller (route handlers)
│       ├── todo-items.service.ts      # Business logic
│       ├── todo-items.repository.ts   # Data access (Prisma queries)
│       └── dto/
│           ├── create-todo-item.dto.ts
│           ├── update-todo-item.dto.ts
│           └── todo-item-response.dto.ts
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

### 4. Start the development server

```bash
npm run start:dev
```

The Swagger UI is available at <http://localhost:3000/swagger>.

## API endpoints

| Method | URL | Description |
|--------|-----|-------------|
| `GET` | `/api/todo-items` | List all todo items (paginated) |
| `GET` | `/api/todo-items/incomplete` | List incomplete items (paginated) |
| `GET` | `/api/todo-items/:id` | Get a single todo item |
| `POST` | `/api/todo-items` | Create a todo item |
| `PUT` | `/api/todo-items/:id` | Update a todo item |
| `PATCH` | `/api/todo-items/:id/complete` | Mark a todo item as complete |
| `DELETE` | `/api/todo-items/:id` | Delete a todo item |

### Pagination query parameters

| Parameter | Default | Description |
|-----------|---------|-------------|
| `page` | `1` | Page number (1-based) |
| `pageSize` | `20` | Items per page |

## Switching databases

Update `DATABASE_URL` in `.env` and change `provider` in `prisma/schema.prisma`:

- **PostgreSQL**: `"postgresql://user:password@localhost:5432/todo_db?schema=public"`
- **MySQL**: `"mysql://user:password@localhost:3306/todo_db"`

Then re-run:

```bash
npx prisma migrate dev
```

## Docker

### Build the image

```bash
# Run from src/backend/nodejs/
docker build -t todo-api-nodejs .
```

### Run the container

Provide the `DATABASE_URL` environment variable (SQLite path inside the container):

```bash
docker run -d -p 3000:3000 \
  -e DATABASE_URL="file:/app/data/todo.db" \
  --name todo-api-nodejs todo-api-nodejs
```

The API is available at <http://localhost:3000>.  
Swagger UI: <http://localhost:3000/swagger>

### Persist the SQLite database

Mount a volume so the database survives container restarts:

```bash
docker run -d -p 3000:3000 \
  -e DATABASE_URL="file:/app/data/todo.db" \
  -v todo-nodejs-data:/app/data \
  --name todo-api-nodejs todo-api-nodejs
```

### Stop and remove the container

```bash
docker stop todo-api-nodejs
docker rm todo-api-nodejs
```
