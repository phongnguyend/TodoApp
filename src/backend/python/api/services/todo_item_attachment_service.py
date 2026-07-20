from abc import ABC, abstractmethod

from fastapi import HTTPException, status
from sqlalchemy.orm import Session

from shared.models.todo_item_attachment import TodoItemAttachment
from api.repositories.file_repository import FileRepository, IFileRepository
from api.repositories.todo_item_attachment_repository import (
    ITodoItemAttachmentRepository,
    TodoItemAttachmentRepository,
)
from api.repositories.todo_item_repository import ITodoItemRepository, TodoItemRepository
from api.schemas.todo_item_attachment import CreateTodoItemAttachmentRequest, TodoItemAttachmentResponse


class ITodoItemAttachmentService(ABC):
    @abstractmethod
    def get_all(self, todo_item_id: int) -> list[TodoItemAttachmentResponse]: ...

    @abstractmethod
    def get_by_id(self, todo_item_id: int, attachment_id: int) -> TodoItemAttachmentResponse: ...

    @abstractmethod
    def create(self, todo_item_id: int, request: CreateTodoItemAttachmentRequest,
               actor_user_id: int | None = None) -> TodoItemAttachmentResponse: ...

    @abstractmethod
    def update(
        self, todo_item_id: int, attachment_id: int, request: CreateTodoItemAttachmentRequest,
        actor_user_id: int | None = None
    ) -> TodoItemAttachmentResponse: ...

    @abstractmethod
    def delete(self, todo_item_id: int, attachment_id: int) -> None: ...


class TodoItemAttachmentService(ITodoItemAttachmentService):
    def __init__(
        self,
        attachment_repository: ITodoItemAttachmentRepository,
        todo_item_repository: ITodoItemRepository,
        file_repository: IFileRepository,
    ) -> None:
        self._attachments = attachment_repository
        self._todos = todo_item_repository
        self._files = file_repository

    def _require_todo(self, todo_item_id: int) -> None:
        if self._todos.get_by_id(todo_item_id) is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND, detail=f"Todo item {todo_item_id} not found."
            )

    def _require_file(self, file_id: int) -> None:
        if self._files.get_by_id(file_id) is None:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail=f"File {file_id} not found.")

    def _get_or_404(self, todo_item_id: int, attachment_id: int) -> TodoItemAttachment:
        attachment = self._attachments.get_by_id_for_todo_item(todo_item_id, attachment_id)
        if attachment is None:
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND,
                detail=f"Attachment {attachment_id} not found for todo item {todo_item_id}.",
            )
        return attachment

    def get_all(self, todo_item_id: int) -> list[TodoItemAttachmentResponse]:
        self._require_todo(todo_item_id)
        return [
            TodoItemAttachmentResponse.model_validate(item)
            for item in self._attachments.get_by_todo_item_id(todo_item_id)
        ]

    def get_by_id(self, todo_item_id: int, attachment_id: int) -> TodoItemAttachmentResponse:
        self._require_todo(todo_item_id)
        return TodoItemAttachmentResponse.model_validate(self._get_or_404(todo_item_id, attachment_id))

    def create(self, todo_item_id: int, request: CreateTodoItemAttachmentRequest,
               actor_user_id: int | None = None) -> TodoItemAttachmentResponse:
        self._require_todo(todo_item_id)
        self._require_file(request.file_id)
        existing = self._attachments.get_by_todo_item_and_file(todo_item_id, request.file_id)
        if existing is not None:
            return TodoItemAttachmentResponse.model_validate(existing)
        created = self._attachments.add(TodoItemAttachment(
            todo_item_id=todo_item_id, file_id=request.file_id, created_by_user_id=actor_user_id
        ))
        return TodoItemAttachmentResponse.model_validate(created)

    def update(
        self, todo_item_id: int, attachment_id: int, request: CreateTodoItemAttachmentRequest,
        actor_user_id: int | None = None
    ) -> TodoItemAttachmentResponse:
        self._require_todo(todo_item_id)
        self._require_file(request.file_id)
        attachment = self._get_or_404(todo_item_id, attachment_id)
        existing = self._attachments.get_by_todo_item_and_file(todo_item_id, request.file_id)
        if existing is not None and existing.id != attachment.id:
            return TodoItemAttachmentResponse.model_validate(existing)
        attachment.file_id = request.file_id
        attachment.updated_by_user_id = actor_user_id
        return TodoItemAttachmentResponse.model_validate(self._attachments.update(attachment))

    def delete(self, todo_item_id: int, attachment_id: int) -> None:
        self._require_todo(todo_item_id)
        self._attachments.delete(self._get_or_404(todo_item_id, attachment_id))


def get_todo_item_attachment_service(db: Session) -> ITodoItemAttachmentService:
    return TodoItemAttachmentService(
        TodoItemAttachmentRepository(db), TodoItemRepository(db), FileRepository(db)
    )

