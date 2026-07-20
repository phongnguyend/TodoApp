# Build stage
FROM ghcr.io/astral-sh/uv:0.11.29 AS uv
FROM python:3.13-slim AS build
WORKDIR /app

COPY --from=uv /uv /uvx /bin/
COPY pyproject.toml uv.lock ./
RUN uv sync --locked --no-dev --no-install-project

# Runtime stage
FROM python:3.13-slim AS runtime
WORKDIR /app

COPY --from=build /app/.venv /app/.venv
COPY . .

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1
ENV PATH="/app/.venv/bin:$PATH"

# The worker shares the same database as the API container.
# Mount the same volume / point DATABASE_URL at the shared DB.
CMD ["sh", "-c", "alembic upgrade head && python -m worker.main"]
