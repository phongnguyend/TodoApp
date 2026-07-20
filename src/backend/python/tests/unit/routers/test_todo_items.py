from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException
from fastapi.testclient import TestClient

from shared.database import get_db
from api.main import app
from api.routers import todo_items
from api.schemas.todo_item import ImportResult, PaginatedResponse, TodoItemResponse
from api.security import get_current_user_id


# ── Helpers ───────────────────────────────────────────────────────────────────

def _make_response(
    id: int = 1,
    title: str = "Test Todo",
    description: str | None = None,
    is_completed: bool = False,
) -> TodoItemResponse:
    return TodoItemResponse(
        id=id,
        title=title,
        description=description,
        is_completed=is_completed,
        created_at=datetime(2024, 1, 1, 12, 0, 0, tzinfo=timezone.utc),
        updated_at=None,
    )


def _paginated(
    items: list[TodoItemResponse],
    total: int | None = None,
    page: int = 1,
    page_size: int = 20,
) -> PaginatedResponse[TodoItemResponse]:
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

    app.dependency_overrides[todo_items._service] = override_service
    app.dependency_overrides[get_db] = override_db
    app.dependency_overrides[get_current_user_id] = lambda: 7
    with TestClient(app) as c:
        yield c
    app.dependency_overrides.clear()


# ── GET /api/todo-items/ ──────────────────────────────────────────────────────

class TestGetAll:
    def test_returns_200_with_items(self, client, mock_service):
        items = [_make_response(1, "A"), _make_response(2, "B")]
        mock_service.get_all.return_value = _paginated(items, 2)

        response = client.get("/api/todo-items/")

        assert response.status_code == 200
        data = response.json()
        assert data["total"] == 2
        assert len(data["items"]) == 2

    def test_passes_pagination_params_to_service(self, client, mock_service):
        mock_service.get_all.return_value = _paginated([], 0)

        client.get("/api/todo-items/?page=2&page_size=5")

        mock_service.get_all.assert_called_once_with(page=2, page_size=5)

    def test_returns_empty_list(self, client, mock_service):
        mock_service.get_all.return_value = _paginated([], 0)

        response = client.get("/api/todo-items/")

        assert response.status_code == 200
        assert response.json()["items"] == []


# ── GET /api/todo-items/incomplete ───────────────────────────────────────────

class TestGetIncomplete:
    def test_returns_200(self, client, mock_service):
        mock_service.get_incomplete.return_value = _paginated([_make_response(1)])

        response = client.get("/api/todo-items/incomplete")

        assert response.status_code == 200
        assert len(response.json()["items"]) == 1

    def test_passes_pagination_params_to_service(self, client, mock_service):
        mock_service.get_incomplete.return_value = _paginated([], 0)

        client.get("/api/todo-items/incomplete?page=2&page_size=10")

        mock_service.get_incomplete.assert_called_once_with(page=2, page_size=10)


# ── GET /api/todo-items/{todo_id} ─────────────────────────────────────────────

class TestGetById:
    def test_returns_200_with_item(self, client, mock_service):
        mock_service.get_by_id.return_value = _make_response(1, "Found")

        response = client.get("/api/todo-items/1")

        assert response.status_code == 200
        assert response.json()["id"] == 1
        assert response.json()["title"] == "Found"

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.get_by_id.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.get("/api/todo-items/99")

        assert response.status_code == 404

    def test_calls_service_with_correct_id(self, client, mock_service):
        mock_service.get_by_id.return_value = _make_response(42, "Item")

        client.get("/api/todo-items/42")

        mock_service.get_by_id.assert_called_once_with(42)


# ── POST /api/todo-items/ ─────────────────────────────────────────────────────

class TestCreate:
    def test_returns_201_on_success(self, client, mock_service):
        mock_service.create.return_value = _make_response(1, "New Todo", "A description")

        response = client.post("/api/todo-items/", json={"title": "New Todo", "description": "A description"})

        assert response.status_code == 201
        assert response.json()["title"] == "New Todo"

    def test_returns_201_without_description(self, client, mock_service):
        mock_service.create.return_value = _make_response(1, "No Desc")

        response = client.post("/api/todo-items/", json={"title": "No Desc"})

        assert response.status_code == 201

    def test_returns_422_when_title_missing(self, client, mock_service):
        response = client.post("/api/todo-items/", json={})

        assert response.status_code == 422

    def test_returns_422_when_title_empty(self, client, mock_service):
        response = client.post("/api/todo-items/", json={"title": ""})

        assert response.status_code == 422


# ── PUT /api/todo-items/{todo_id} ─────────────────────────────────────────────

class TestUpdate:
    def test_returns_200_on_success(self, client, mock_service):
        mock_service.update.return_value = _make_response(1, "Updated Title")

        response = client.put("/api/todo-items/1", json={"title": "Updated Title"})

        assert response.status_code == 200
        assert response.json()["title"] == "Updated Title"

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.update.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.put("/api/todo-items/99", json={"title": "X"})

        assert response.status_code == 404

    def test_calls_service_with_correct_id(self, client, mock_service):
        mock_service.update.return_value = _make_response(7, "Title")

        client.put("/api/todo-items/7", json={"title": "Title"})

        call_args = mock_service.update.call_args
        assert call_args[0][0] == 7


# ── PATCH /api/todo-items/{todo_id}/complete ──────────────────────────────────

class TestMarkComplete:
    def test_returns_200_on_success(self, client, mock_service):
        mock_service.mark_complete.return_value = _make_response(1, "Done", is_completed=True)

        response = client.patch("/api/todo-items/1/complete")

        assert response.status_code == 200
        assert response.json()["is_completed"] is True

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.mark_complete.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.patch("/api/todo-items/99/complete")

        assert response.status_code == 404

    def test_calls_service_with_correct_id(self, client, mock_service):
        mock_service.mark_complete.return_value = _make_response(5, is_completed=True)

        client.patch("/api/todo-items/5/complete")

        mock_service.mark_complete.assert_called_once_with(5)


# ── DELETE /api/todo-items/{todo_id} ──────────────────────────────────────────

class TestDelete:
    def test_returns_204_on_success(self, client, mock_service):
        mock_service.delete.return_value = None

        response = client.delete("/api/todo-items/1")

        assert response.status_code == 204

    def test_returns_404_when_not_found(self, client, mock_service):
        mock_service.delete.side_effect = HTTPException(status_code=404, detail="Not found")

        response = client.delete("/api/todo-items/99")

        assert response.status_code == 404

    def test_calls_service_with_correct_id(self, client, mock_service):
        mock_service.delete.return_value = None

        client.delete("/api/todo-items/3")

        mock_service.delete.assert_called_once_with(3)


# ── POST /api/todo-items/import/csv ───────────────────────────────────────────

class TestImportCsv:
    def test_returns_200_with_result(self, client, mock_service):
        mock_service.import_csv.return_value = ImportResult(imported=2, failed=0, errors=[])

        csv_content = b"title,description,is_completed\nBuy milk,Whole milk,false\n"
        response = client.post(
            "/api/todo-items/import/csv",
            files={"file": ("todo_items.csv", csv_content, "text/csv")},
        )

        assert response.status_code == 200
        data = response.json()
        assert data["imported"] == 2
        assert data["failed"] == 0

    def test_returns_result_with_errors(self, client, mock_service):
        mock_service.import_csv.return_value = ImportResult(
            imported=1, failed=1, errors=[{"row": 2, "error": "Title is required."}]
        )

        csv_content = b"title\n,\n"
        response = client.post(
            "/api/todo-items/import/csv",
            files={"file": ("todo_items.csv", csv_content, "text/csv")},
        )

        assert response.status_code == 200
        data = response.json()
        assert data["failed"] == 1
        assert data["errors"][0]["row"] == 2

    def test_returns_422_when_file_missing(self, client, mock_service):
        response = client.post("/api/todo-items/import/csv")

        assert response.status_code == 422


# ── GET /api/todo-items/export/csv ────────────────────────────────────────────

class TestExportCsv:
    def test_returns_200_with_csv_content(self, client, mock_service):
        mock_service.export_csv.return_value = "id,title,description,is_completed,created_at,updated_at\n"

        response = client.get("/api/todo-items/export/csv")

        assert response.status_code == 200
        assert response.headers["content-type"].startswith("text/csv")
        assert "attachment" in response.headers["content-disposition"]
        assert response.text == "id,title,description,is_completed,created_at,updated_at\n"


# ── POST /api/todo-items/import/excel ─────────────────────────────────────────

class TestImportExcel:
    def test_returns_200_with_result(self, client, mock_service):
        mock_service.import_excel.return_value = ImportResult(imported=2, failed=0, errors=[])

        excel_content = b"fake xlsx bytes"
        response = client.post(
            "/api/todo-items/import/excel",
            files={"file": ("todo_items.xlsx", excel_content, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")},
        )

        assert response.status_code == 200
        data = response.json()
        assert data["imported"] == 2
        assert data["failed"] == 0

    def test_returns_result_with_errors(self, client, mock_service):
        mock_service.import_excel.return_value = ImportResult(
            imported=1, failed=1, errors=[{"row": 2, "error": "Title is required."}]
        )

        excel_content = b"fake xlsx bytes"
        response = client.post(
            "/api/todo-items/import/excel",
            files={"file": ("todo_items.xlsx", excel_content, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")},
        )

        assert response.status_code == 200
        data = response.json()
        assert data["failed"] == 1
        assert data["errors"][0]["row"] == 2

    def test_returns_422_when_file_missing(self, client, mock_service):
        response = client.post("/api/todo-items/import/excel")

        assert response.status_code == 422


# ── GET /api/todo-items/export/excel ──────────────────────────────────────────

class TestExportExcel:
    def test_returns_200_with_excel_content(self, client, mock_service):
        mock_service.export_excel.return_value = b"fake xlsx bytes"

        response = client.get("/api/todo-items/export/excel")

        assert response.status_code == 200
        assert response.headers["content-type"] == "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        assert "attachment" in response.headers["content-disposition"]
        assert response.content == b"fake xlsx bytes"
