# Build stage
FROM golang:1.23-alpine AS build
WORKDIR /src

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/worker ./cmd/worker

# Runtime stage
FROM alpine:3.21
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=build /app/worker .

# The worker shares the same database volume as the API container.
# Mount the same named volume at /app so todo.db is shared.
ENTRYPOINT ["./worker"]
