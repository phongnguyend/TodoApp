from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException
from fastapi.testclient import TestClient

from api.main import app
from api.routers import tokens
from api.schemas.user import TokenResponse


@pytest.fixture
def service():
    return MagicMock()


@pytest.fixture
def client(service):
    app.dependency_overrides[tokens._service] = lambda: service
    with TestClient(app) as test_client:
        yield test_client
    app.dependency_overrides.clear()


def test_create_token_returns_jwt_and_no_store_headers(client, service):
    service.create_token.return_value = TokenResponse(access_token="header.payload.signature", expires_in=3600)

    response = client.post("/api/tokens", json={"email": "Alice@Example.com", "password": "password123"})

    assert response.status_code == 200
    assert response.json() == {
        "access_token": "header.payload.signature", "token_type": "Bearer", "expires_in": 3600,
    }
    assert response.headers["cache-control"] == "no-store"
    assert response.headers["pragma"] == "no-cache"


def test_create_token_returns_non_disclosing_401(client, service):
    service.create_token.side_effect = HTTPException(
        status_code=401, detail="Invalid email or password.", headers={"WWW-Authenticate": "Bearer"})

    response = client.post("/api/tokens", json={"email": "missing@example.com", "password": "wrong"})

    assert response.status_code == 401
    assert response.json() == {"error": "Invalid email or password."}
    assert response.headers["www-authenticate"] == "Bearer"


def test_create_token_rejects_invalid_request_with_400(client):
    assert client.post("/api/tokens", json={"email": "invalid"}).status_code == 400
