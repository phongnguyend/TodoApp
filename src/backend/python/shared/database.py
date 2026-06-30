from collections.abc import Generator

from sqlalchemy import create_engine
from sqlalchemy.orm import DeclarativeBase, Session, sessionmaker

from shared.config import settings

engine = create_engine(
    settings.DATABASE_URL,
    # Only needed for SQLite to allow multi-threaded access
    connect_args={"check_same_thread": False} if settings.DATABASE_URL.startswith("sqlite") else {},
)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


class Base(DeclarativeBase):
    pass


def get_db() -> Generator[Session, None, None]:
    """Dependency that provides a database session per request (analogous to EF DbContext scoped lifetime)."""
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
