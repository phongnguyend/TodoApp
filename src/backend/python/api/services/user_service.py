import hashlib
import time
from abc import ABC, abstractmethod
from datetime import datetime, timezone
from urllib.parse import quote

from fastapi import HTTPException, status
from sqlalchemy.orm import Session

from shared.config import settings
from shared.models.email_log import EmailLog
from shared.models.user import User
from api.repositories.user_repository import IUserRepository, UserRepository
from api.schemas.todo_item import PaginatedResponse
from api.schemas.user import (
    ChangePasswordRequest, ConfirmPasswordResetRequest, CreateUserRequest, ResetPasswordRequest,
    SignUpRequest, UpdateProfileRequest, UpdateUserRequest, UserResponse,
)
from api.security import create_signed_token, decode_signed_token, hash_password, verify_password


class IUserService(ABC):
    @abstractmethod
    def get_all(self, page: int = 1, page_size: int = 20) -> PaginatedResponse[UserResponse]: ...
    @abstractmethod
    def get_by_id(self, user_id: int) -> UserResponse: ...
    @abstractmethod
    def create(self, request: CreateUserRequest) -> UserResponse: ...
    @abstractmethod
    def update(self, user_id: int, request: UpdateUserRequest) -> UserResponse: ...
    @abstractmethod
    def set_active(self, user_id: int, is_active: bool) -> UserResponse: ...
    @abstractmethod
    def signup(self, request: SignUpRequest) -> UserResponse: ...
    @abstractmethod
    def get_profile(self, user_id: int) -> UserResponse: ...
    @abstractmethod
    def update_profile(self, user_id: int, request: UpdateProfileRequest) -> UserResponse: ...
    @abstractmethod
    def change_password(self, user_id: int, request: ChangePasswordRequest) -> None: ...
    @abstractmethod
    def request_password_reset(self, request: ResetPasswordRequest) -> None: ...
    @abstractmethod
    def confirm_password_reset(self, request: ConfirmPasswordResetRequest) -> None: ...


class UserService(IUserService):
    def __init__(self, repository: IUserRepository) -> None:
        self._repo = repository

    def _get_or_404(self, user_id: int) -> User:
        user = self._repo.get_by_id(user_id)
        if user is None:
            raise HTTPException(status_code=status.HTTP_404_NOT_FOUND, detail=f"User {user_id} not found.")
        return user

    def _ensure_unique(self, username: str, email: str, excluding_id: int | None = None) -> None:
        if self._repo.username_exists(username.strip(), excluding_id):
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Username is already in use.")
        if self._repo.email_exists(email.strip().lower(), excluding_id):
            raise HTTPException(status_code=status.HTTP_409_CONFLICT, detail="Email is already in use.")

    @staticmethod
    def _response(user: User) -> UserResponse:
        return UserResponse.model_validate(user)

    def get_all(self, page: int = 1, page_size: int = 20) -> PaginatedResponse[UserResponse]:
        page = max(1, page)
        page_size = min(settings.MAX_PAGE_SIZE, max(1, page_size))
        users, total = self._repo.get_all(skip=(page - 1) * page_size, limit=page_size)
        return PaginatedResponse(items=[self._response(user) for user in users], total=total, page=page,
                                 page_size=page_size, total_pages=-(-total // page_size))

    def get_by_id(self, user_id: int) -> UserResponse:
        return self._response(self._get_or_404(user_id))

    def create(self, request: CreateUserRequest) -> UserResponse:
        username, email = request.username.strip(), request.email.strip().lower()
        self._ensure_unique(username, email)
        user = User(username=username, email=email, password_hash=hash_password(request.password),
                    is_active=request.is_active, created_at=datetime.now(timezone.utc))
        return self._response(self._repo.add(user))

    def update(self, user_id: int, request: UpdateUserRequest) -> UserResponse:
        user = self._get_or_404(user_id)
        username = request.username.strip() if request.username is not None else user.username
        email = request.email.strip().lower() if request.email is not None else user.email
        self._ensure_unique(username, email, user_id)
        user.username, user.email = username, email
        if request.password is not None:
            user.password_hash = hash_password(request.password)
        user.updated_at = datetime.now(timezone.utc)
        return self._response(self._repo.update(user))

    def set_active(self, user_id: int, is_active: bool) -> UserResponse:
        user = self._get_or_404(user_id)
        user.is_active = is_active
        user.updated_at = datetime.now(timezone.utc)
        return self._response(self._repo.update(user))

    def signup(self, request: SignUpRequest) -> UserResponse:
        return self.create(CreateUserRequest(username=request.username, email=request.email,
                                             password=request.password, is_active=True))

    def get_profile(self, user_id: int) -> UserResponse:
        return self.get_by_id(user_id)

    def update_profile(self, user_id: int, request: UpdateProfileRequest) -> UserResponse:
        return self.update(user_id, UpdateUserRequest(username=request.username, email=request.email))

    def change_password(self, user_id: int, request: ChangePasswordRequest) -> None:
        user = self._get_or_404(user_id)
        if not user.is_active:
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="The user account is inactive.")
        if not verify_password(request.current_password, user.password_hash):
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST, detail="The current password is incorrect.")
        user.password_hash = hash_password(request.new_password)
        user.updated_at = datetime.now(timezone.utc)
        self._repo.update(user)

    def request_password_reset(self, request: ResetPasswordRequest) -> None:
        user = self._repo.get_by_email(request.email.strip().lower())
        if user is None or not user.is_active:
            return
        password_fingerprint = hashlib.sha256(user.password_hash.encode()).hexdigest()
        token = create_signed_token({"sub": str(user.id), "exp": int(time.time()) + settings.PASSWORD_RESET_TOKEN_LIFETIME_MINUTES * 60,
                                     "password": password_fingerprint}, settings.PASSWORD_RESET_SECRET_KEY)
        separator = "&" if "?" in settings.PASSWORD_RESET_CONFIRMATION_URL else "?"
        reset_url = f"{settings.PASSWORD_RESET_CONFIRMATION_URL}{separator}token={quote(token)}"
        self._repo.add_email_log(EmailLog(recipient=user.email, subject="Reset your Todo API password",
            body=f"Use this link to reset your password: {reset_url}\n\nThis link expires in {settings.PASSWORD_RESET_TOKEN_LIFETIME_MINUTES} minutes.",
            status="pending", created_at=datetime.now(timezone.utc)))

    def confirm_password_reset(self, request: ConfirmPasswordResetRequest) -> None:
        try:
            payload = decode_signed_token(request.token, settings.PASSWORD_RESET_SECRET_KEY)
            user = self._repo.get_by_id(int(payload["sub"]))
            fingerprint = hashlib.sha256(user.password_hash.encode()).hexdigest() if user else ""
            if user is None or not user.is_active or payload.get("password") != fingerprint:
                raise ValueError
        except (HTTPException, KeyError, TypeError, ValueError):
            raise HTTPException(status_code=status.HTTP_400_BAD_REQUEST,
                                detail="The password reset token is invalid or expired.")
        user.password_hash = hash_password(request.new_password)
        user.updated_at = datetime.now(timezone.utc)
        self._repo.update(user)


def get_user_service(db: Session) -> IUserService:
    return UserService(UserRepository(db))
