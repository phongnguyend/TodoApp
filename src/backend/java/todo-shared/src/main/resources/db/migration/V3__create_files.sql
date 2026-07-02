-- V3: files table
-- Stores metadata about uploaded files; the actual file content is stored on disk
-- at the path recorded in `location`.
-- Works with H2 (dev) and MySQL. For PostgreSQL use BIGSERIAL instead of BIGINT AUTO_INCREMENT.

CREATE TABLE IF NOT EXISTS files
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    extension    VARCHAR(20)  NOT NULL,
    size         BIGINT       NOT NULL,
    content_type VARCHAR(100),
    location     VARCHAR(500) NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP
);
