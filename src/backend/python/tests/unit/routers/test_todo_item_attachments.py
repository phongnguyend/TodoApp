from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException
from fastapi.testclient import TestClient

from shared.database import get_db
from api.main import app
from api.routers import todo_items
from api.schemas.todo_item_attachment import TodoItemAttachmentResponse
from api.security import get_current_user_id


def _response(attachment_id: int = 1, todo_id: int = 10, file_id: int = 5) -> TodoItemAttachmentResponse:
    return TodoItemAttachmentResponse(
        id=attachment_id,
        todo_item_id=todo_id,
        file_id=file_id,
        created_at=datetime(2024, 1, 1, tzinfo=timezone.utc),
        updated_at=None,
    )


@pytest.fixture
def mock_service():
    return MagicMock()


@pytest.fixture
def mock_db():
    return MagicMock()


@pytest.fixture
def client(mock_service, mock_db):
    app.dependency_overrides[todo_items._attachment_service] = lambda: mock_service
    app.dependency_overrides[get_current_user_id] = lambda: 7

    def override_db():
        yield mock_db

    app.dependency_overrides[get_db] = override_db
    with TestClient(app) as test_client:
        yield test_client
    app.dependency_overrides.clear()


class TestGetAttachments:
    def test_returns_attachment_list(self, client, mock_service):
        mock_service.get_all.return_value = [_response()]

        response = client.get("/api/todo-items/10/attachments")

        assert response.status_code == 200
        assert response.json()[0]["file_id"] == 5
        mock_service.get_all.assert_called_once_with(10)

    def test_returns_404_from_service(self, client, mock_service):
        mock_service.get_all.side_effect = HTTPException(status_code=404, detail="Todo item 99 not found.")

        response = client.get("/api/todo-items/99/attachments")

        assert response.status_code == 404


class TestCreateAttachment:
    def test_returns_201_and_commits(self, client, mock_service, mock_db):
        mock_service.create.return_value = _response()

        response = client.post("/api/todo-items/10/attachments", json={"file_id": 5})

        assert response.status_code == 201
        assert response.json()["id"] == 1
        assert mock_service.create.call_args.args[0] == 10
        assert mock_service.create.call_args.args[1].file_id == 5
        mock_db.commit.assert_called_once()

    def test_returns_422_without_file_id(self, client, mock_service):
        response = client.post("/api/todo-items/10/attachments", json={})

        assert response.status_code == 422
        mock_service.create.assert_not_called()


class TestGetAttachmentById:
    def test_returns_attachment(self, client, mock_service):
        mock_service.get_by_id.return_value = _response(3)

        response = client.get("/api/todo-items/10/attachments/3")

        assert response.status_code == 200
        assert response.json()["id"] == 3
        mock_service.get_by_id.assert_called_once_with(10, 3)


class TestUpdateAttachment:
    def test_returns_updated_attachment_and_commits(self, client, mock_service, mock_db):
        mock_service.update.return_value = _response(3, file_id=6)

        response = client.put("/api/todo-items/10/attachments/3", json={"file_id": 6})

        assert response.status_code == 200
        assert response.json()["file_id"] == 6
        assert mock_service.update.call_args.args[:2] == (10, 3)
        assert mock_service.update.call_args.args[2].file_id == 6
        mock_db.commit.assert_called_once()


class TestDeleteAttachment:
    def test_returns_204_and_commits(self, client, mock_service, mock_db):
        response = client.delete("/api/todo-items/10/attachments/3")

        assert response.status_code == 204
        mock_service.delete.assert_called_once_with(10, 3)
        mock_db.commit.assert_called_once()

    def test_returns_404_from_service(self, client, mock_service, mock_db):
        mock_service.delete.side_effect = HTTPException(status_code=404, detail="Attachment not found.")

        response = client.delete("/api/todo-items/10/attachments/99")

        assert response.status_code == 404
        mock_db.commit.assert_not_called()
