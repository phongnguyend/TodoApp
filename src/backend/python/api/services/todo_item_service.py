from abc import ABC, abstractmethod

from fastapi import HTTPException, status
from sqlalchemy.orm import Session

from shared.models.todo_item import TodoItem
from api.repositories.todo_item_repository import ITodoItemRepository, TodoItemRepository
from api.schemas.todo_item import (
    CreateTodoItemRequest,
    PaginatedResponse,
    TodoItemResponse,
    UpdateTodoItemRequest,
)


class ITodoItemService(ABC):
    """Service interface — mirrors ITodoService in C#."""

    @abstractmethod
    def get_all(self, page: int, page_size: int) -> PaginatedResponse[TodoItemResponse]: ...

    @abstractmethod
    def get_incomplete(self, page: int, page_size: int) -> PaginatedResponse[TodoItemResponse]: ...

    @abstractmethod
    def get_by_id(self, todo_id: int) -> TodoItemResponse: ...

    @abstractmethod
    def create(self, request: CreateTodoItemRequest) -> TodoItemResponse: ...

    @abstractmethod
    def update(self, todo_id: int, request: UpdateTodoItemRequest) -> TodoItemResponse: ...

    @abstractmethod
    def delete(self, todo_id: int) -> None: ...

    @abstractmethod
    def mark_complete(self, todo_id: int) -> TodoItemResponse: ...


class TodoItemService(ITodoItemService):
    """Business-logic layer (analogous to an ASP.NET Core service registered via DI)."""

    def __init__(self, repository: ITodoItemRepository) -> None:
        self._repo = repository

    # ── Helpers ───────────────────────────────────────────────────────────────

    def _get_or_404(self, todo_id: int) -> TodoItem:
        todo = self._repo.get_by_id(todo_id)
        if todo is None:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail=f"Todo item {todo_id} not found.")
        return todo

    @staticmethod
    def _to_paginated(items: list[TodoItem], total: int, page: int, page_size: int) -> PaginatedResponse[TodoItemResponse]:
        return PaginatedResponse(
            items=[TodoItemResponse.model_validate(item) for item in items],
            total=total,
            page=page,
            page_size=page_size,
            total_pages=-(-total // page_size),  # ceiling division
        )

    # ── Queries ───────────────────────────────────────────────────────────────

    def get_all(self, page: int = 1, page_size: int = 20) -> PaginatedResponse[TodoItemResponse]:
        skip = (page - 1) * page_size
        items, total = self._repo.get_all(skip=skip, limit=page_size)
        return self._to_paginated(items, total, page, page_size)

    def get_incomplete(self, page: int = 1, page_size: int = 20) -> PaginatedResponse[TodoItemResponse]:
        skip = (page - 1) * page_size
        items, total = self._repo.get_incomplete(skip=skip, limit=page_size)
        return self._to_paginated(items, total, page, page_size)

    def get_by_id(self, todo_id: int) -> TodoItemResponse:
        todo = self._get_or_404(todo_id)
        return TodoItemResponse.model_validate(todo)

    # ── Commands ──────────────────────────────────────────────────────────────

    def create(self, request: CreateTodoItemRequest) -> TodoItemResponse:
        todo = TodoItem(title=request.title, description=request.description)
        created = self._repo.add(todo)
        return TodoItemResponse.model_validate(created)

    def update(self, todo_id: int, request: UpdateTodoItemRequest) -> TodoItemResponse:
        todo = self._get_or_404(todo_id)
        if request.title is not None:
            todo.title = request.title
        if request.description is not None:
            todo.description = request.description
        if request.is_completed is not None:
            todo.is_completed = request.is_completed
        updated = self._repo.update(todo)
        return TodoItemResponse.model_validate(updated)

    def delete(self, todo_id: int) -> None:
        todo = self._get_or_404(todo_id)
        self._repo.delete(todo)

    def mark_complete(self, todo_id: int) -> TodoItemResponse:
        todo = self._get_or_404(todo_id)
        todo.is_completed = True
        updated = self._repo.update(todo)
        return TodoItemResponse.model_validate(updated)


# ── Dependency factory (used by FastAPI Depends) ──────────────────────────────

def get_todo_service(db: Session) -> ITodoItemService:
    repository = TodoItemRepository(db)
    return TodoItemService(repository)
