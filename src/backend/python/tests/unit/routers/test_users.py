from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException
from fastapi.testclient import TestClient

from shared.database import get_db
from api.main import app
from api.routers import users
from api.schemas.todo_item import PaginatedResponse
from api.schemas.user import UserResponse
from api.security import get_current_user_id


def _response(user_id: int = 1, username: str = "alice", is_active: bool = True) -> UserResponse:
    return UserResponse(
        id=user_id,
        username=username,
        email=f"{username}@example.com",
        is_active=is_active,
        created_at=datetime(2026, 1, 1, tzinfo=timezone.utc),
        updated_at=None,
    )


def _page(items: list[UserResponse], page: int = 1, page_size: int = 20):
    return PaginatedResponse(
        items=items,
        total=len(items),
        page=page,
        page_size=page_size,
        total_pages=1 if items else 0,
    )


@pytest.fixture
def mock_service():
    return MagicMock()


@pytest.fixture
def mock_db():
    return MagicMock()


@pytest.fixture
def client(mock_service, mock_db):
    app.dependency_overrides[users._service] = lambda: mock_service
    app.dependency_overrides[get_current_user_id] = lambda: 7

    def override_db():
        yield mock_db

    app.dependency_overrides[get_db] = override_db
    with TestClient(app) as test_client:
        yield test_client
    app.dependency_overrides.clear()


class TestGetAll:
    def test_returns_paginated_users_and_passes_query(self, client, mock_service):
        mock_service.get_all.return_value = _page([_response()], page=2, page_size=5)

        response = client.get("/api/users/?page=2&page_size=5")

        assert response.status_code == 200
        assert response.json()["items"][0]["username"] == "alice"
        mock_service.get_all.assert_called_once_with(page=2, page_size=5)


class TestGetById:
    def test_returns_user(self, client, mock_service):
        mock_service.get_by_id.return_value = _response(42)

        response = client.get("/api/users/42")

        assert response.status_code == 200
        assert response.json()["id"] == 42
        mock_service.get_by_id.assert_called_once_with(42)

    def test_propagates_not_found(self, client, mock_service):
        mock_service.get_by_id.side_effect = HTTPException(status_code=404, detail="User not found.")

        assert client.get("/api/users/99").status_code == 404


class TestCreate:
    def test_returns_201_and_commits(self, client, mock_service, mock_db):
        mock_service.create.return_value = _response()

        response = client.post(
            "/api/users/",
            json={"username": "alice", "email": "alice@example.com", "password": "password123"},
        )

        assert response.status_code == 201
        request = mock_service.create.call_args.args[0]
        assert request.username == "alice"
        assert request.is_active is True
        mock_db.commit.assert_called_once()

    def test_rejects_invalid_email_and_short_password(self, client, mock_service, mock_db):
        response = client.post(
            "/api/users/", json={"username": "alice", "email": "invalid", "password": "short"}
        )

        assert response.status_code == 422
        mock_service.create.assert_not_called()
        mock_db.commit.assert_not_called()

    def test_propagates_conflict_without_committing(self, client, mock_service, mock_db):
        mock_service.create.side_effect = HTTPException(status_code=409, detail="Username is already in use.")

        response = client.post(
            "/api/users/",
            json={"username": "alice", "email": "alice@example.com", "password": "password123"},
        )

        assert response.status_code == 409
        mock_db.commit.assert_not_called()


class TestUpdate:
    def test_returns_updated_user_and_commits(self, client, mock_service, mock_db):
        mock_service.update.return_value = _response(3, "updated")

        response = client.put("/api/users/3", json={"username": "updated"})

        assert response.status_code == 200
        assert response.json()["username"] == "updated"
        assert mock_service.update.call_args.args[0] == 3
        assert mock_service.update.call_args.args[1].username == "updated"
        mock_db.commit.assert_called_once()


class TestActivation:
    @pytest.mark.parametrize(
        ("path", "active"),
        [("activate", True), ("deactivate", False)],
    )
    def test_sets_active_state_and_commits(self, path, active, client, mock_service, mock_db):
        mock_service.set_active.return_value = _response(4, is_active=active)

        response = client.patch(f"/api/users/4/{path}")

        assert response.status_code == 200
        assert response.json()["is_active"] is active
        mock_service.set_active.assert_called_once_with(4, active, 7)
        mock_db.commit.assert_called_once()


class TestSignup:
    def test_returns_201_and_commits(self, client, mock_service, mock_db):
        mock_service.signup.return_value = _response()

        response = client.post(
            "/api/users/signup",
            json={"username": "alice", "email": "alice@example.com", "password": "password123"},
        )

        assert response.status_code == 201
        assert mock_service.signup.call_args.args[0].email == "alice@example.com"
        mock_db.commit.assert_called_once()


class TestProfile:
    def test_get_uses_authenticated_user_id(self, client, mock_service):
        mock_service.get_profile.return_value = _response(7)

        response = client.get("/api/users/profile")

        assert response.status_code == 200
        mock_service.get_profile.assert_called_once_with(7)

    def test_update_uses_authenticated_user_id_and_commits(self, client, mock_service, mock_db):
        mock_service.update_profile.return_value = _response(7, "new-name")

        response = client.put("/api/users/profile", json={"username": "new-name"})

        assert response.status_code == 200
        assert mock_service.update_profile.call_args.args[0] == 7
        assert mock_service.update_profile.call_args.args[1].username == "new-name"
        mock_db.commit.assert_called_once()

    def test_requires_authentication(self, client):
        app.dependency_overrides.pop(get_current_user_id)

        response = client.get("/api/users/profile")

        assert response.status_code == 401


class TestPasswords:
    def test_change_returns_204_and_commits(self, client, mock_service, mock_db):
        response = client.post(
            "/api/users/password/change",
            json={"current_password": "password123", "new_password": "new-password123"},
        )

        assert response.status_code == 204
        assert mock_service.change_password.call_args.args[0] == 7
        assert mock_service.change_password.call_args.args[1].new_password == "new-password123"
        mock_db.commit.assert_called_once()

    def test_reset_always_returns_202_and_commits(self, client, mock_service, mock_db):
        response = client.post("/api/users/password/reset", json={"email": "missing@example.com"})

        assert response.status_code == 202
        assert "account exists" in response.json()["message"]
        assert mock_service.request_password_reset.call_args.args[0].email == "missing@example.com"
        mock_db.commit.assert_called_once()

    def test_confirm_returns_204_and_commits(self, client, mock_service, mock_db):
        response = client.post(
            "/api/users/password/confirm",
            json={"token": "signed-token", "new_password": "new-password123"},
        )

        assert response.status_code == 204
        assert mock_service.confirm_password_reset.call_args.args[0].token == "signed-token"
        mock_db.commit.assert_called_once()

    def test_confirm_rejects_short_password(self, client, mock_service, mock_db):
        response = client.post(
            "/api/users/password/confirm", json={"token": "signed-token", "new_password": "short"}
        )

        assert response.status_code == 422
        mock_service.confirm_password_reset.assert_not_called()
        mock_db.commit.assert_not_called()
