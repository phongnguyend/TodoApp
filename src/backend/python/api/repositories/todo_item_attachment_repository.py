from abc import abstractmethod

from sqlalchemy.orm import Session

from shared.models.todo_item_attachment import TodoItemAttachment
from api.repositories.base_repository import BaseRepository, IRepository


class ITodoItemAttachmentRepository(IRepository[TodoItemAttachment]):
    @abstractmethod
    def get_by_todo_item_id(self, todo_item_id: int) -> list[TodoItemAttachment]: ...

    @abstractmethod
    def get_by_id_for_todo_item(self, todo_item_id: int, attachment_id: int) -> TodoItemAttachment | None: ...

    @abstractmethod
    def get_by_todo_item_and_file(self, todo_item_id: int, file_id: int) -> TodoItemAttachment | None: ...


class TodoItemAttachmentRepository(BaseRepository[TodoItemAttachment], ITodoItemAttachmentRepository):
    def __init__(self, db: Session) -> None:
        super().__init__(db, TodoItemAttachment)

    def get_by_todo_item_id(self, todo_item_id: int) -> list[TodoItemAttachment]:
        return (
            self._db.query(TodoItemAttachment)
            .filter(TodoItemAttachment.todo_item_id == todo_item_id)
            .order_by(TodoItemAttachment.created_at.desc(), TodoItemAttachment.id.desc())
            .all()
        )

    def get_by_id_for_todo_item(self, todo_item_id: int, attachment_id: int) -> TodoItemAttachment | None:
        return (
            self._db.query(TodoItemAttachment)
            .filter(
                TodoItemAttachment.id == attachment_id,
                TodoItemAttachment.todo_item_id == todo_item_id,
            )
            .first()
        )

    def get_by_todo_item_and_file(self, todo_item_id: int, file_id: int) -> TodoItemAttachment | None:
        return (
            self._db.query(TodoItemAttachment)
            .filter(
                TodoItemAttachment.todo_item_id == todo_item_id,
                TodoItemAttachment.file_id == file_id,
            )
            .first()
        )

