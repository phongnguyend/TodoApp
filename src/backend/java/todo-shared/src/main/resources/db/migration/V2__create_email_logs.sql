-- V2: email_logs table
-- Audit trail for every outbound email attempt.
-- Records persist regardless of SMTP delivery outcome.
-- Works with H2 (dev) and MySQL. For PostgreSQL use BIGSERIAL instead of BIGINT AUTO_INCREMENT.

CREATE TABLE IF NOT EXISTS email_logs
(
    id            BIGINT AUTO_INCREMENT PRIMARY KEY,
    recipient     VARCHAR(255) NOT NULL,
    subject       VARCHAR(500) NOT NULL,
    body          TEXT         NOT NULL,
    status        VARCHAR(50)  NOT NULL DEFAULT 'pending',
    created_at    TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    sent_at       TIMESTAMP,
    error_message TEXT
);
