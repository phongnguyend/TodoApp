from typing import Annotated

from fastapi import APIRouter, Depends, File as UploadFileParam, UploadFile, status
from fastapi.responses import StreamingResponse
from sqlalchemy.orm import Session

from shared.database import get_db
from api.schemas.todo_item import (
    CreateTodoItemRequest,
    ImportResult,
    PaginatedResponse,
    TodoItemResponse,
    UpdateTodoItemRequest,
)
from api.services.todo_item_service import ITodoItemService, get_todo_service

router = APIRouter(prefix="/api/todo-items", tags=["Todo Items"])

# ── Dependency aliases (analogous to constructor-injected services in ASP.NET Core) ──

DbDep = Annotated[Session, Depends(get_db)]


def _service(db: DbDep) -> ITodoItemService:
    return get_todo_service(db)


ServiceDep = Annotated[ITodoItemService, Depends(_service)]


# ── Endpoints ─────────────────────────────────────────────────────────────────

@router.get("/", response_model=PaginatedResponse[TodoItemResponse], summary="Get all todo items")
def get_all(service: ServiceDep, page: int = 1, page_size: int = 20):
    return service.get_all(page=page, page_size=page_size)


@router.get("/incomplete", response_model=PaginatedResponse[TodoItemResponse], summary="Get incomplete todo items")
def get_incomplete(service: ServiceDep, page: int = 1, page_size: int = 20):
    return service.get_incomplete(page=page, page_size=page_size)


@router.post("/import/csv", response_model=ImportResult, summary="Import todo items from a CSV file")
def import_csv(service: ServiceDep, db: DbDep, file: Annotated[UploadFile, UploadFileParam(...)]):
    result = service.import_csv(file)
    db.commit()
    return result


@router.get("/export/csv", summary="Export todo items as a CSV file")
def export_csv(service: ServiceDep):
    content = service.export_csv()
    return StreamingResponse(
        iter([content]),
        media_type="text/csv",
        headers={"Content-Disposition": "attachment; filename=todo_items.csv"},
    )


@router.get("/{todo_id}", response_model=TodoItemResponse, summary="Get a todo item by ID")
def get_by_id(todo_id: int, service: ServiceDep):
    return service.get_by_id(todo_id)


@router.post("/", response_model=TodoItemResponse, status_code=status.HTTP_201_CREATED, summary="Create a todo item")
def create(request: CreateTodoItemRequest, service: ServiceDep, db: DbDep):
    todo = service.create(request)
    db.commit()
    return todo


@router.put("/{todo_id}", response_model=TodoItemResponse, summary="Update a todo item")
def update(todo_id: int, request: UpdateTodoItemRequest, service: ServiceDep, db: DbDep):
    todo = service.update(todo_id, request)
    db.commit()
    return todo


@router.patch("/{todo_id}/complete", response_model=TodoItemResponse, summary="Mark a todo item as complete")
def mark_complete(todo_id: int, service: ServiceDep, db: DbDep):
    todo = service.mark_complete(todo_id)
    db.commit()
    return todo


@router.delete("/{todo_id}", status_code=status.HTTP_204_NO_CONTENT, summary="Delete a todo item")
def delete(todo_id: int, service: ServiceDep, db: DbDep):
    service.delete(todo_id)
    db.commit()
