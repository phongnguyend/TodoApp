import os
import uuid
from abc import ABC, abstractmethod
from pathlib import Path

from fastapi import HTTPException, UploadFile, status
from sqlalchemy.orm import Session

from shared.config import settings
from shared.models.file import File as FileModel
from api.repositories.file_repository import FileRepository, IFileRepository
from api.schemas.file import FileResponse
from api.schemas.todo_item import PaginatedResponse


class IFileService(ABC):
    """Service interface - mirrors IFileService in C#."""

    @abstractmethod
    def get_all(self, page: int, page_size: int) -> PaginatedResponse[FileResponse]: ...

    @abstractmethod
    def get_by_id(self, file_id: int) -> FileResponse: ...

    @abstractmethod
    def upload(self, upload_file: UploadFile, actor_user_id: int | None = None) -> FileResponse: ...

    @abstractmethod
    def get_download_target(self, file_id: int) -> tuple[str, str, str]: ...

    @abstractmethod
    def delete(self, file_id: int) -> None: ...


class FileService(IFileService):
    """Business-logic layer (analogous to an ASP.NET Core service registered via DI)."""

    def __init__(self, repository: IFileRepository) -> None:
        self._repo = repository

    # ── Helpers ───────────────────────────────────────────────────────────────

    def _get_or_404(self, file_id: int) -> FileModel:
        file = self._repo.get_by_id(file_id)
        if file is None:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail=f"File {file_id} not found.")
        return file

    @staticmethod
    def _to_paginated(items: list[FileModel], total: int, page: int, page_size: int) -> PaginatedResponse[FileResponse]:
        return PaginatedResponse(
            items=[FileResponse.model_validate(item) for item in items],
            total=total,
            page=page,
            page_size=page_size,
            total_pages=-(-total // page_size),  # ceiling division
        )

    # ── Queries ───────────────────────────────────────────────────────────────

    def get_all(self, page: int = 1, page_size: int = 20) -> PaginatedResponse[FileResponse]:
        skip = (page - 1) * page_size
        items, total = self._repo.get_all(skip=skip, limit=page_size)
        return self._to_paginated(items, total, page, page_size)

    def get_by_id(self, file_id: int) -> FileResponse:
        file = self._get_or_404(file_id)
        return FileResponse.model_validate(file)

    # ── Commands ──────────────────────────────────────────────────────────────

    def upload(self, upload_file: UploadFile, actor_user_id: int | None = None) -> FileResponse:
        # Strip any directory components from the client-supplied name to prevent path traversal.
        original_name = Path(upload_file.filename or "unnamed").name
        extension = Path(original_name).suffix.lstrip(".")

        content = upload_file.file.read()
        if len(content) > settings.MAX_UPLOAD_SIZE_BYTES:
            raise HTTPException(
                status_code=status.HTTP_413_REQUEST_ENTITY_TOO_LARGE,
                detail=f"File exceeds the maximum allowed size of {settings.MAX_UPLOAD_SIZE_BYTES} bytes.",
            )

        storage_dir = Path(settings.FILE_STORAGE_PATH)
        storage_dir.mkdir(parents=True, exist_ok=True)

        # A random prefix avoids collisions/overwrites between uploads that share a name.
        stored_name = f"{uuid.uuid4().hex}_{original_name}"
        location = storage_dir / stored_name
        with open(location, "wb") as f:
            f.write(content)

        file = FileModel(
            name=original_name,
            extension=extension,
            size=len(content),
            content_type=upload_file.content_type,
            location=str(location),
            created_by_user_id=actor_user_id,
        )
        created = self._repo.add(file)
        return FileResponse.model_validate(created)

    def get_download_target(self, file_id: int) -> tuple[str, str, str]:
        """Returns (absolute_path, download_name, content_type) for streaming the file back to the client."""
        file = self._get_or_404(file_id)
        if not os.path.isfile(file.location):
            raise HTTPException(
                status_code=status.HTTP_404_NOT_FOUND, detail=f"File {file_id} content not found on disk."
            )
        return file.location, file.name, file.content_type or "application/octet-stream"

    def delete(self, file_id: int) -> None:
        file = self._get_or_404(file_id)
        self._repo.delete(file)
        if os.path.isfile(file.location):
            os.remove(file.location)


# ── Dependency factory (used by FastAPI Depends) ──────────────────────────────

def get_file_service(db: Session) -> IFileService:
    repository = FileRepository(db)
    return FileService(repository)
