import base64
import hashlib
import hmac
import json
import secrets
import time
from typing import Annotated

from fastapi import Depends, HTTPException, status
from fastapi.security import HTTPAuthorizationCredentials, HTTPBearer
import jwt

from shared.config import settings


def hash_password(password: str) -> str:
    salt = secrets.token_bytes(16)
    digest = hashlib.pbkdf2_hmac("sha256", password.encode(), salt, settings.PASSWORD_HASH_ITERATIONS)
    return f"pbkdf2_sha256${settings.PASSWORD_HASH_ITERATIONS}${salt.hex()}${digest.hex()}"


def verify_password(password: str, encoded: str) -> bool:
    try:
        algorithm, iterations, salt, expected = encoded.split("$", 3)
        if algorithm != "pbkdf2_sha256":
            return False
        actual = hashlib.pbkdf2_hmac("sha256", password.encode(), bytes.fromhex(salt), int(iterations))
        return hmac.compare_digest(actual.hex(), expected)
    except (ValueError, TypeError):
        return False


def _b64encode(data: bytes) -> str:
    return base64.urlsafe_b64encode(data).rstrip(b"=").decode()


def _b64decode(data: str) -> bytes:
    return base64.urlsafe_b64decode(data + "=" * (-len(data) % 4))


def create_signed_token(payload: dict, secret: str) -> str:
    body = _b64encode(json.dumps(payload, separators=(",", ":"), sort_keys=True).encode())
    signature = _b64encode(hmac.new(secret.encode(), body.encode(), hashlib.sha256).digest())
    return f"{body}.{signature}"


def create_jwt(payload: dict, secret: str) -> str:
    return jwt.encode(payload, secret, algorithm="HS256", headers={"typ": "JWT"})


def decode_signed_token(token: str, secret: str) -> dict:
    try:
        body, supplied_signature = token.split(".", 1)
        expected_signature = _b64encode(hmac.new(secret.encode(), body.encode(), hashlib.sha256).digest())
        if not hmac.compare_digest(supplied_signature, expected_signature):
            raise ValueError
        payload = json.loads(_b64decode(body))
        if payload.get("exp", 0) < int(time.time()):
            raise ValueError
        return payload
    except (ValueError, TypeError, json.JSONDecodeError):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid or expired token.")


bearer = HTTPBearer(auto_error=False)


def _decode_user_id(credentials: HTTPAuthorizationCredentials) -> int:
    try:
        payload = jwt.decode(
            credentials.credentials,
            settings.JWT_SECRET_KEY,
            algorithms=["HS256"],
            options={"require": ["sub", "exp"]},
        )
    except jwt.PyJWTError:
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid or expired token.")
    try:
        user_id = int(payload["sub"])
        if user_id < 1:
            raise ValueError
        return user_id
    except (KeyError, TypeError, ValueError):
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid authentication token.")


def get_current_user_id(
    credentials: Annotated[HTTPAuthorizationCredentials | None, Depends(bearer)],
) -> int:
    if credentials is None or credentials.scheme.lower() != "bearer":
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Authentication required.")
    return _decode_user_id(credentials)


def get_optional_current_user_id(
    credentials: Annotated[HTTPAuthorizationCredentials | None, Depends(bearer)],
) -> int | None:
    if credentials is None:
        return None
    if credentials.scheme.lower() != "bearer":
        raise HTTPException(status_code=status.HTTP_401_UNAUTHORIZED, detail="Invalid authentication token.")
    return _decode_user_id(credentials)


CurrentUserId = Annotated[int, Depends(get_current_user_id)]
