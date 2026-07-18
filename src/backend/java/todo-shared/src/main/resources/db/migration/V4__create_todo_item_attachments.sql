CREATE TABLE IF NOT EXISTS todo_item_attachments
(
    id           BIGINT AUTO_INCREMENT PRIMARY KEY,
    todo_item_id BIGINT    NOT NULL,
    file_id      BIGINT    NOT NULL,
    created_at   TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP,
    CONSTRAINT fk_todo_item_attachments_todo_item
        FOREIGN KEY (todo_item_id) REFERENCES todo_items (id) ON DELETE CASCADE,
    CONSTRAINT fk_todo_item_attachments_file
        FOREIGN KEY (file_id) REFERENCES files (id) ON DELETE CASCADE,
    CONSTRAINT uq_todo_item_attachments_todo_file UNIQUE (todo_item_id, file_id)
);

CREATE INDEX IF NOT EXISTS ix_todo_item_attachments_file_id
    ON todo_item_attachments (file_id);
