# Build stage
FROM python:3.13-slim AS build
WORKDIR /app

RUN pip install --upgrade pip
COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

# Runtime stage
FROM python:3.13-slim AS runtime
WORKDIR /app

COPY --from=build /install /usr/local
COPY . .

ENV PYTHONDONTWRITEBYTECODE=1
ENV PYTHONUNBUFFERED=1

# The worker shares the same database as the API container.
# Mount the same volume / point DATABASE_URL at the shared DB.
CMD ["sh", "-c", "alembic upgrade head && python -m TodoWorker.main"]
