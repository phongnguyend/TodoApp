from datetime import datetime
import re

from pydantic import BaseModel, Field, field_validator


def _validate_email(value: str) -> str:
    value = value.strip().lower()
    if len(value) > 255 or not re.fullmatch(r"[^\s@]+@[^\s@]+\.[^\s@]+", value):
        raise ValueError("A valid email address is required.")
    return value


class CreateUserRequest(BaseModel):
    username: str = Field(min_length=1, max_length=50)
    email: str
    password: str = Field(min_length=8, max_length=128)
    is_active: bool = True

    _email = field_validator("email")(_validate_email)


class UpdateUserRequest(BaseModel):
    username: str | None = Field(default=None, min_length=1, max_length=50)
    email: str | None = None
    password: str | None = Field(default=None, min_length=8, max_length=128)

    @field_validator("email")
    @classmethod
    def validate_email(cls, value: str | None) -> str | None:
        return _validate_email(value) if value is not None else None


class SignUpRequest(BaseModel):
    username: str = Field(min_length=1, max_length=50)
    email: str
    password: str = Field(min_length=8, max_length=128)

    _email = field_validator("email")(_validate_email)


class ChangePasswordRequest(BaseModel):
    current_password: str = Field(min_length=1)
    new_password: str = Field(min_length=8, max_length=128)


class ResetPasswordRequest(BaseModel):
    email: str
    _email = field_validator("email")(_validate_email)


class ConfirmPasswordResetRequest(BaseModel):
    token: str = Field(min_length=1)
    new_password: str = Field(min_length=8, max_length=128)


class UpdateProfileRequest(BaseModel):
    username: str | None = Field(default=None, min_length=1, max_length=50)
    email: str | None = None

    @field_validator("email")
    @classmethod
    def validate_email(cls, value: str | None) -> str | None:
        return _validate_email(value) if value is not None else None


class UserResponse(BaseModel):
    model_config = {"from_attributes": True}

    id: int
    username: str
    email: str
    is_active: bool
    created_at: datetime
    updated_at: datetime | None
