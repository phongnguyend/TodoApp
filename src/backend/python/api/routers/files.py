from typing import Annotated

from fastapi import APIRouter, Depends, File as UploadFileParam, UploadFile, status
from fastapi.responses import FileResponse as FastAPIFileResponse
from sqlalchemy.orm import Session

from shared.database import get_db
from api.schemas.file import FileResponse
from api.schemas.todo_item import PaginatedResponse
from api.services.file_service import IFileService, get_file_service
from api.security import get_optional_current_user_id

router = APIRouter(prefix="/api/files", tags=["Files"])

# ── Dependency aliases (analogous to constructor-injected services in ASP.NET Core) ──

DbDep = Annotated[Session, Depends(get_db)]


def _service(db: DbDep) -> IFileService:
    return get_file_service(db)


ServiceDep = Annotated[IFileService, Depends(_service)]


# ── Endpoints ─────────────────────────────────────────────────────────────────

@router.get("/", response_model=PaginatedResponse[FileResponse], summary="Get all uploaded files")
def get_all(service: ServiceDep, page: int = 1, page_size: int = 20):
    return service.get_all(page=page, page_size=page_size)


@router.get("/{file_id}", response_model=FileResponse, summary="Get a file's metadata by ID")
def get_by_id(file_id: int, service: ServiceDep):
    return service.get_by_id(file_id)


@router.get("/{file_id}/download", summary="Download a file's content")
def download(file_id: int, service: ServiceDep):
    path, name, content_type = service.get_download_target(file_id)
    return FastAPIFileResponse(path=path, filename=name, media_type=content_type)


@router.post("/", response_model=FileResponse, status_code=status.HTTP_201_CREATED, summary="Upload a file")
def upload(service: ServiceDep, db: DbDep, upload_file: Annotated[UploadFile, UploadFileParam(...)],
           actor_user_id: int | None = Depends(get_optional_current_user_id)):
    file = service.upload(upload_file) if actor_user_id is None else service.upload(upload_file, actor_user_id)
    db.commit()
    return file


@router.delete("/{file_id}", status_code=status.HTTP_204_NO_CONTENT, summary="Delete a file")
def delete(file_id: int, service: ServiceDep, db: DbDep):
    service.delete(file_id)
    db.commit()
