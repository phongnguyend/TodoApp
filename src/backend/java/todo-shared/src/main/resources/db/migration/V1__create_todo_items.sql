-- V1: Initial schema
-- Analogous to the first EF Core migration (InitialCreate)
-- Works with H2 (dev) and MySQL. For PostgreSQL use BIGSERIAL instead of BIGINT AUTO_INCREMENT.

CREATE TABLE IF NOT EXISTS todo_items
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    title        VARCHAR(200) NOT NULL,
    description  TEXT,
    is_completed BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP
);
