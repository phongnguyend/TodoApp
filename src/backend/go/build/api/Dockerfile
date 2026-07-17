# Build stage
FROM golang:1.25-alpine AS build
WORKDIR /src

RUN apk add --no-cache gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/server ./cmd/api

# Runtime stage
FROM alpine:3.21
WORKDIR /app

RUN apk add --no-cache libc6-compat

COPY --from=build /app/server .

EXPOSE 8080

ENTRYPOINT ["./server"]
