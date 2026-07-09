import csv
import io
from abc import ABC, abstractmethod

from fastapi import HTTPException, UploadFile, status
from sqlalchemy.orm import Session

from shared.models.todo_item import TodoItem
from api.repositories.todo_item_repository import ITodoItemRepository, TodoItemRepository
from api.schemas.todo_item import (
    CreateTodoItemRequest,
    ImportResult,
    ImportRowError,
    PaginatedResponse,
    TodoItemResponse,
    UpdateTodoItemRequest,
)

# ── CSV import/export helpers ───────────────────────────────────────────────────

_CSV_FIELDNAMES = ["id", "title", "description", "is_completed", "created_at", "updated_at"]
_TRUE_VALUES = {"1", "true", "yes", "y"}


def _parse_bool(value: str | None) -> bool:
    return (value or "").strip().lower() in _TRUE_VALUES


class ITodoItemService(ABC):
    """Service interface - mirrors ITodoService in C#."""

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

    @abstractmethod
    def import_csv(self, file: UploadFile) -> ImportResult: ...

    @abstractmethod
    def export_csv(self) -> str: ...


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

    # ── CSV import/export ─────────────────────────────────────────────────────

    def import_csv(self, file: UploadFile) -> ImportResult:
        raw = file.file.read()
        text = raw.decode("utf-8-sig")
        reader = csv.DictReader(io.StringIO(text))

        imported = 0
        errors: list[ImportRowError] = []
        for row_number, row in enumerate(reader, start=2):  # header is row 1
            title = (row.get("title") or "").strip()
            if not title:
                errors.append(ImportRowError(row=row_number, error="Title is required."))
                continue

            description = (row.get("description") or "").strip() or None
            is_completed = _parse_bool(row.get("is_completed"))

            todo = TodoItem(title=title, description=description, is_completed=is_completed)
            self._repo.add(todo)
            imported += 1

        return ImportResult(imported=imported, failed=len(errors), errors=errors)

    def export_csv(self) -> str:
        items = self._repo.get_all_items()

        buffer = io.StringIO()
        writer = csv.writer(buffer)
        writer.writerow(_CSV_FIELDNAMES)
        for item in items:
            writer.writerow([
                item.id,
                item.title,
                item.description or "",
                item.is_completed,
                item.created_at.isoformat() if item.created_at else "",
                item.updated_at.isoformat() if item.updated_at else "",
            ])
        return buffer.getvalue()


# ── Dependency factory (used by FastAPI Depends) ──────────────────────────────

def get_todo_service(db: Session) -> ITodoItemService:
    repository = TodoItemRepository(db)
    return TodoItemService(repository)
