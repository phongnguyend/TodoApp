from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.config import settings
from app.routers import todo_items

# ── Application bootstrap (analogous to Program.cs / Startup.cs) ──────────────

app = FastAPI(
    title=settings.APP_NAME,
    version=settings.APP_VERSION,
    docs_url="/swagger",        # Swagger UI at /swagger (mirrors ASP.NET Core default)
    redoc_url="/redoc",
    openapi_url="/swagger/v1/swagger.json",
)

# ── Middleware ─────────────────────────────────────────────────────────────────

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# ── Routers (analogous to MapControllers / minimal-API route groups) ──────────

app.include_router(todo_items.router)


@app.get("/", include_in_schema=False)
def health_check():
    return {"status": "healthy", "app": settings.APP_NAME, "version": settings.APP_VERSION}
