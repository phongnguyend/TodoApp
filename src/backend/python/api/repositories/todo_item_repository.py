from abc import abstractmethod

from sqlalchemy.orm import Session

from shared.models.todo_item import TodoItem
from api.repositories.base_repository import BaseRepository, IRepository


class ITodoItemRepository(IRepository[TodoItem]):
    """Contract for the todo-item repository - mirrors an interface in C#."""

    @abstractmethod
    def get_by_title(self, title: str) -> TodoItem | None: ...

    @abstractmethod
    def get_incomplete(self, skip: int = 0, limit: int = 20) -> tuple[list[TodoItem], int]: ...

    @abstractmethod
    def get_all_items(self) -> list[TodoItem]: ...


class TodoItemRepository(BaseRepository[TodoItem], ITodoItemRepository):
    """SQLAlchemy-backed todo-item repository (analogous to EF Repository implementation)."""

    def __init__(self, db: Session) -> None:
        super().__init__(db, TodoItem)

    def get_by_title(self, title: str) -> TodoItem | None:
        return self._db.query(TodoItem).filter(TodoItem.title == title).first()

    def get_incomplete(self, skip: int = 0, limit: int = 20) -> tuple[list[TodoItem], int]:
        query = self._db.query(TodoItem).filter(TodoItem.is_completed.is_(False))
        total = query.count()
        items = query.offset(skip).limit(limit).all()
        return items, total

    def get_all_items(self) -> list[TodoItem]:
        return self._db.query(TodoItem).all()
