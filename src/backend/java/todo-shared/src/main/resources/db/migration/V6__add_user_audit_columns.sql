ALTER TABLE todo_items ADD COLUMN created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE todo_items ADD COLUMN updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE todo_item_attachments ADD COLUMN created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE todo_item_attachments ADD COLUMN updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE email_logs ADD COLUMN created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE email_logs ADD COLUMN updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE files ADD COLUMN created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE files ADD COLUMN updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;

ALTER TABLE users ADD COLUMN created_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
ALTER TABLE users ADD COLUMN updated_by_user_id BIGINT REFERENCES users(id) ON DELETE SET NULL;
