from abc import ABC, abstractmethod
from typing import Generic, TypeVar

from sqlalchemy.orm import Session

from app.database import Base

ModelT = TypeVar("ModelT", bound=Base)


class IRepository(ABC, Generic[ModelT]):
    """Generic repository interface — mirrors IRepository<T> in C#."""

    @abstractmethod
    def get_by_id(self, entity_id: int) -> ModelT | None: ...

    @abstractmethod
    def get_all(self, skip: int = 0, limit: int = 20) -> tuple[list[ModelT], int]: ...

    @abstractmethod
    def add(self, entity: ModelT) -> ModelT: ...

    @abstractmethod
    def update(self, entity: ModelT) -> ModelT: ...

    @abstractmethod
    def delete(self, entity: ModelT) -> None: ...


class BaseRepository(IRepository[ModelT], Generic[ModelT]):
    """Concrete base repository backed by SQLAlchemy (analogous to EF GenericRepository<T>)."""

    def __init__(self, db: Session, model: type[ModelT]) -> None:
        self._db = db
        self._model = model

    def get_by_id(self, entity_id: int) -> ModelT | None:
        return self._db.get(self._model, entity_id)

    def get_all(self, skip: int = 0, limit: int = 20) -> tuple[list[ModelT], int]:
        total: int = self._db.query(self._model).count()
        items: list[ModelT] = self._db.query(self._model).offset(skip).limit(limit).all()
        return items, total

    def add(self, entity: ModelT) -> ModelT:
        self._db.add(entity)
        self._db.flush()
        self._db.refresh(entity)
        return entity

    def update(self, entity: ModelT) -> ModelT:
        self._db.flush()
        self._db.refresh(entity)
        return entity

    def delete(self, entity: ModelT) -> None:
        self._db.delete(entity)
        self._db.flush()
