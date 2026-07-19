from fastapi import FastAPI, HTTPException, Request
from fastapi.exception_handlers import http_exception_handler, request_validation_exception_handler
from fastapi.exceptions import RequestValidationError
from fastapi.responses import JSONResponse
from fastapi.encoders import jsonable_encoder
from fastapi.middleware.cors import CORSMiddleware

from shared.config import settings
from api.routers import files, todo_items, tokens, users

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
app.include_router(files.router)
app.include_router(users.router)
app.include_router(tokens.router)


@app.exception_handler(RequestValidationError)
async def validation_exception_handler(request: Request, exc: RequestValidationError):
    if request.url.path == "/api/tokens":
        return JSONResponse(
            status_code=400,
            content={"error": "The request is invalid.", "errors": jsonable_encoder(exc.errors())},
        )
    return await request_validation_exception_handler(request, exc)


@app.exception_handler(HTTPException)
async def api_http_exception_handler(request: Request, exc: HTTPException):
    if request.url.path == "/api/tokens" and exc.status_code == 401:
        return JSONResponse(status_code=401, content={"error": str(exc.detail)}, headers=exc.headers)
    return await http_exception_handler(request, exc)


@app.get("/", include_in_schema=False)
def health_check():
    return {"status": "healthy", "app": settings.APP_NAME, "version": settings.APP_VERSION}
