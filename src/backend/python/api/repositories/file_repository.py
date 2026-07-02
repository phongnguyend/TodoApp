from sqlalchemy.orm import Session

from shared.models.file import File
from api.repositories.base_repository import BaseRepository, IRepository


class IFileRepository(IRepository[File]):
    """Contract for the file repository - mirrors an interface in C#."""


class FileRepository(BaseRepository[File], IFileRepository):
    """SQLAlchemy-backed file repository (analogous to EF Repository implementation)."""

    def __init__(self, db: Session) -> None:
        super().__init__(db, File)
