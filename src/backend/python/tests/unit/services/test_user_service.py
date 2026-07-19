from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException

from shared.config import settings
from shared.models.user import User
from api.schemas.user import ChangePasswordRequest, CreateUserRequest, UpdateUserRequest
from api.security import hash_password, verify_password
from api.services.user_service import UserService


@pytest.fixture(autouse=True)
def fast_password_hashing(monkeypatch):
    monkeypatch.setattr(settings, "PASSWORD_HASH_ITERATIONS", 1_000)


def _user(user_id: int = 1, password: str = "password123") -> User:
    user = User(username="alice", email="alice@example.com", password_hash=hash_password(password))
    user.id = user_id
    user.is_active = True
    user.created_at = datetime(2026, 1, 1, tzinfo=timezone.utc)
    user.updated_at = None
    return user


def test_create_normalizes_email_and_hashes_password():
    repository = MagicMock()
    repository.username_exists.return_value = False
    repository.email_exists.return_value = False

    def add(user):
        user.id = 1
        return user

    repository.add.side_effect = add
    result = UserService(repository).create(
        CreateUserRequest(username=" Alice ", email="ALICE@Example.com", password="password123")
    )

    created = repository.add.call_args.args[0]
    assert result.username == "Alice"
    assert result.email == "alice@example.com"
    assert created.password_hash != "password123"
    assert verify_password("password123", created.password_hash)


def test_create_rejects_duplicate_username():
    repository = MagicMock()
    repository.username_exists.return_value = True

    with pytest.raises(HTTPException) as error:
        UserService(repository).create(
            CreateUserRequest(username="alice", email="alice@example.com", password="password123")
        )

    assert error.value.status_code == 409
    repository.add.assert_not_called()


def test_update_preserves_fields_that_are_not_supplied():
    repository = MagicMock()
    user = _user()
    repository.get_by_id.return_value = user
    repository.username_exists.return_value = False
    repository.email_exists.return_value = False
    repository.update.side_effect = lambda value: value

    result = UserService(repository).update(1, UpdateUserRequest(username="new-name"))

    assert result.username == "new-name"
    assert result.email == "alice@example.com"
    assert result.updated_at is not None


def test_change_password_checks_current_password():
    repository = MagicMock()
    user = _user()
    repository.get_by_id.return_value = user
    repository.update.side_effect = lambda value: value
    service = UserService(repository)

    with pytest.raises(HTTPException) as error:
        service.change_password(1, ChangePasswordRequest(current_password="wrong", new_password="new-password123"))
    assert error.value.status_code == 400

    service.change_password(
        1, ChangePasswordRequest(current_password="password123", new_password="new-password123")
    )
    assert verify_password("new-password123", user.password_hash)
