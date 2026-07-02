from datetime import datetime, timezone
from io import BytesIO
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException
from fastapi.testclient import TestClient

from shared.database import get_db
from api.main import app
from api.routers import files
from api.schemas.file import FileResponse
from api.schemas.todo_item import PaginatedResponse


# ── Helpers ───────────────────────────────────────────────────────────────────

def _make_response(
    id: int = 1,
    name: str = "test.txt",
    extension: str = "txt",
    size: int = 100,
    content_type: str | None = "text/plain",
) -> FileResponse:
    return FileResponse(
        id=id,
        name=name,
        extension=extension,
        size=size,
        content_type=content_type,
        created_at=datetime(2024, 1, 1, 12, 0, 0, tzinfo=timezone.utc),
        updated_at=None,
    )


def _paginated(
    items: list[FileResponse],
    total: int | None = None,
    page: int = 1,
    page_size: int = 20,
) -> PaginatedResponse[FileResponse]:
    if total is None:
        total = len(items)
    return PaginatedResponse(
        items=items,
        total=total,
        page=page,
        page_size=page_size,
        total_pages=max(1, -(-total // page_size)) if total else 1,
    )


# ── Fixtures ──────────────────────────────────────────────────────────────────

@pytest.fixture
def mock_service():
    return MagicMock()


@pytest.fixture
def client(mock_service):
    mock_db = MagicMock()

    def override_service():
        return mock_service

    def override_db():
        yield mock_db

    app.dependency_overrides[files._service] = override_service
    app.dependency_overrides[get_db] = override_db
    with TestClient(app) as c:
        yield c
    app.dependency_overrides.clear()


# ── GET /api/files/ ───────────────────────────────────────────────────────────

class TestGetAll:
    def test_returns_200_with_items(self, client, mock_service):
        items = [_make_response(1, "a.txt"), _make_response(2, "b.txt")]
        mock_service.get_all.return_value = _paginated(items, 2)

        response = client.get("/api/files/")

        assert response.status_code == 200
        data = response.json()
        assert data["total"] == 2
        assert len(data["items"]) == 2

    def test_passes_pagination_params_to_service(self, client, mock_service):
        mock_service.get_all.return_value = _paginated([], 0)

        client.get("/api/files/?page=2&page_size=5")

        mock_service.get_all.assert_called_once_with(page=2, page_size=5)

    def test_returns_empty_list(self, client, mock_service):
        mock_service.get_all.return_value = _paginated([], 0)

        response = client.get("/api/files/")

        assert response.status_code == 200
        assert response.json()["items"] == []


# ── GET /api/files/{file_id} ──────────────────────────────────────────────────

class TestGetById:
    def test_returns_200_with_item(self, client, mock_service):
        mock_service.get_by_id.return_value = _make_response(1, "found.txt")

        response = client.get("/api/files/1")

        assert response.status_code == 200
        assert response.json()["id"] == 1
        assert response.json()["name"] == "found.txt"

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.get_by_id.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.get("/api/files/99")

        assert response.status_code == 404

    def test_calls_service_with_correct_id(self, client, mock_service):
        mock_service.get_by_id.return_value = _make_response(42, "item.txt")

        client.get("/api/files/42")

        mock_service.get_by_id.assert_called_once_with(42)


# ── POST /api/files/ ──────────────────────────────────────────────────────────

class TestUpload:
    def test_returns_201_on_success(self, client, mock_service):
        mock_service.upload.return_value = _make_response(1, "photo.png", "png", 11, "image/png")

        response = client.post(
            "/api/files/",
            files={"upload_file": ("photo.png", BytesIO(b"binary-da"), "image/png")},
        )

        assert response.status_code == 201
        assert response.json()["name"] == "photo.png"

    def test_returns_422_when_no_file_provided(self, client, mock_service):
        response = client.post("/api/files/")

        assert response.status_code == 422


# ── GET /api/files/{file_id}/download ─────────────────────────────────────────

class TestDownload:
    def test_returns_200_with_file_content(self, client, mock_service, tmp_path):
        file_path = tmp_path / "download.txt"
        file_path.write_text("file contents")
        mock_service.get_download_target.return_value = (str(file_path), "download.txt", "text/plain")

        response = client.get("/api/files/1/download")

        assert response.status_code == 200
        assert response.content == b"file contents"

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.get_download_target.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.get("/api/files/99/download")

        assert response.status_code == 404


# ── DELETE /api/files/{file_id} ───────────────────────────────────────────────

class TestDelete:
    def test_returns_204_on_success(self, client, mock_service):
        mock_service.delete.return_value = None

        response = client.delete("/api/files/1")

        assert response.status_code == 204

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.delete.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.delete("/api/files/99")

        assert response.status_code == 404

    def test_calls_service_with_correct_id(self, client, mock_service):
        mock_service.delete.return_value = None

        client.delete("/api/files/3")

        mock_service.delete.assert_called_once_with(3)
