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


settings = Settings()
