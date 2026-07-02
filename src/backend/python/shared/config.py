from pydantic_settings import BaseSettings, SettingsConfigDict


class Settings(BaseSettings):
    model_config = SettingsConfigDict(env_file=".env", env_file_encoding="utf-8")

    APP_NAME: str = "Todo API"
    APP_VERSION: str = "1.0.0"
    DEBUG: bool = False

    DATABASE_URL: str = "sqlite:///./todo.db"

    # Pagination defaults
    DEFAULT_PAGE_SIZE: int = 20
    MAX_PAGE_SIZE: int = 100

    # ── SMTP (used by the background worker) ──────────────────────────────────
    SMTP_HOST: str = "localhost"
    SMTP_PORT: int = 587
    # True  → plain SMTP + STARTTLS (port 587)
    # False → SMTP_SSL           (port 465)
    SMTP_USE_TLS: bool = True
    SMTP_USERNAME: str = ""
    SMTP_PASSWORD: str = ""
    EMAIL_SENDER: str = "noreply@example.com"
    EMAIL_RECIPIENT: str = "admin@example.com"

    # ── Background worker ─────────────────────────────────────────────────────
    WORKER_INTERVAL_MINUTES: int = 60

    # ── File uploads ───────────────────────────────────────────────────────────
    FILE_STORAGE_PATH: str = "./uploads"
    MAX_UPLOAD_SIZE_BYTES: int = 10 * 1024 * 1024  # 10 MB


settings = Settings()
