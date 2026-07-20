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
from api.schemas.todo_item_attachment import CreateTodoItemAttachmentRequest, TodoItemAttachmentResponse
from api.services.todo_item_attachment_service import (
    ITodoItemAttachmentService,
    get_todo_item_attachment_service,
)
from api.security import get_optional_current_user_id

router = APIRouter(prefix="/api/todo-items", tags=["Todo Items"])

# ── Dependency aliases (analogous to constructor-injected services in ASP.NET Core) ──

DbDep = Annotated[Session, Depends(get_db)]


def _service(db: DbDep) -> ITodoItemService:
    return get_todo_service(db)


ServiceDep = Annotated[ITodoItemService, Depends(_service)]


def _attachment_service(db: DbDep) -> ITodoItemAttachmentService:
    return get_todo_item_attachment_service(db)


AttachmentServiceDep = Annotated[ITodoItemAttachmentService, Depends(_attachment_service)]


# ── Endpoints ─────────────────────────────────────────────────────────────────

@router.get("/", response_model=PaginatedResponse[TodoItemResponse], summary="Get all todo items")
def get_all(service: ServiceDep, page: int = 1, page_size: int = 20):
    return service.get_all(page=page, page_size=page_size)


@router.get("/incomplete", response_model=PaginatedResponse[TodoItemResponse], summary="Get incomplete todo items")
def get_incomplete(service: ServiceDep, page: int = 1, page_size: int = 20):
    return service.get_incomplete(page=page, page_size=page_size)


@router.post("/import/csv", response_model=ImportResult, summary="Import todo items from a CSV file")
def import_csv(service: ServiceDep, db: DbDep, file: Annotated[UploadFile, UploadFileParam(...)],
               actor_user_id: int | None = Depends(get_optional_current_user_id)):
    result = service.import_csv(file) if actor_user_id is None else service.import_csv(file, actor_user_id)
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


@router.post("/import/excel", response_model=ImportResult, summary="Import todo items from an Excel file")
def import_excel(service: ServiceDep, db: DbDep, file: Annotated[UploadFile, UploadFileParam(...)],
                 actor_user_id: int | None = Depends(get_optional_current_user_id)):
    result = service.import_excel(file) if actor_user_id is None else service.import_excel(file, actor_user_id)
    db.commit()
    return result


@router.get("/export/excel", summary="Export todo items as an Excel file")
def export_excel(service: ServiceDep):
    content = service.export_excel()
    return StreamingResponse(
        iter([content]),
        media_type="application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        headers={"Content-Disposition": "attachment; filename=todo_items.xlsx"},
    )


@router.get("/{todo_id}", response_model=TodoItemResponse, summary="Get a todo item by ID")
def get_by_id(todo_id: int, service: ServiceDep):
    return service.get_by_id(todo_id)


@router.post("/", response_model=TodoItemResponse, status_code=status.HTTP_201_CREATED, summary="Create a todo item")
def create(request: CreateTodoItemRequest, service: ServiceDep, db: DbDep,
           actor_user_id: int | None = Depends(get_optional_current_user_id)):
    todo = service.create(request) if actor_user_id is None else service.create(request, actor_user_id)
    db.commit()
    return todo


@router.put("/{todo_id}", response_model=TodoItemResponse, summary="Update a todo item")
def update(todo_id: int, request: UpdateTodoItemRequest, service: ServiceDep, db: DbDep,
           actor_user_id: int | None = Depends(get_optional_current_user_id)):
    todo = service.update(todo_id, request) if actor_user_id is None else service.update(todo_id, request, actor_user_id)
    db.commit()
    return todo


@router.patch("/{todo_id}/complete", response_model=TodoItemResponse, summary="Mark a todo item as complete")
def mark_complete(todo_id: int, service: ServiceDep, db: DbDep,
                  actor_user_id: int | None = Depends(get_optional_current_user_id)):
    todo = service.mark_complete(todo_id) if actor_user_id is None else service.mark_complete(todo_id, actor_user_id)
    db.commit()
    return todo


@router.get(
    "/{todo_id}/attachments",
    response_model=list[TodoItemAttachmentResponse],
    summary="Get a todo item's attachments",
)
def get_attachments(todo_id: int, service: AttachmentServiceDep):
    return service.get_all(todo_id)


@router.post(
    "/{todo_id}/attachments",
    response_model=TodoItemAttachmentResponse,
    status_code=status.HTTP_201_CREATED,
    summary="Attach a file to a todo item",
)
def create_attachment(
    todo_id: int, request: CreateTodoItemAttachmentRequest, service: AttachmentServiceDep, db: DbDep,
    actor_user_id: int | None = Depends(get_optional_current_user_id),
):
    attachment = (service.create(todo_id, request) if actor_user_id is None
                  else service.create(todo_id, request, actor_user_id))
    db.commit()
    return attachment


@router.get(
    "/{todo_id}/attachments/{attachment_id}",
    response_model=TodoItemAttachmentResponse,
    summary="Get a todo item attachment by ID",
)
def get_attachment_by_id(todo_id: int, attachment_id: int, service: AttachmentServiceDep):
    return service.get_by_id(todo_id, attachment_id)


@router.put(
    "/{todo_id}/attachments/{attachment_id}",
    response_model=TodoItemAttachmentResponse,
    summary="Update a todo item attachment",
)
def update_attachment(
    todo_id: int,
    attachment_id: int,
    request: CreateTodoItemAttachmentRequest,
    service: AttachmentServiceDep,
    db: DbDep,
    actor_user_id: int | None = Depends(get_optional_current_user_id),
):
    attachment = (service.update(todo_id, attachment_id, request) if actor_user_id is None
                  else service.update(todo_id, attachment_id, request, actor_user_id))
    db.commit()
    return attachment


@router.delete(
    "/{todo_id}/attachments/{attachment_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    summary="Delete a todo item attachment",
)
def delete_attachment(todo_id: int, attachment_id: int, service: AttachmentServiceDep, db: DbDep):
    service.delete(todo_id, attachment_id)
    db.commit()


@router.delete("/{todo_id}", status_code=status.HTTP_204_NO_CONTENT, summary="Delete a todo item")
def delete(todo_id: int, service: ServiceDep, db: DbDep):
    service.delete(todo_id)
    db.commit()
