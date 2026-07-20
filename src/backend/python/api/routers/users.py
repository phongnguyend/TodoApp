from typing import Annotated

from fastapi import APIRouter, Depends, Response, status
from sqlalchemy.orm import Session

from shared.database import get_db
from api.schemas.todo_item import PaginatedResponse
from api.schemas.user import (
    ChangePasswordRequest, ConfirmPasswordResetRequest, CreateUserRequest, ResetPasswordRequest,
    SignUpRequest, UpdateProfileRequest, UpdateUserRequest, UserResponse,
)
from api.security import CurrentUserId, get_current_user_id
from api.services.user_service import IUserService, get_user_service

router = APIRouter(prefix="/api/users", tags=["Users"])
protected_router = APIRouter(
    prefix="/api/users", tags=["Users"], dependencies=[Depends(get_current_user_id)]
)
DbDep = Annotated[Session, Depends(get_db)]


def _service(db: DbDep) -> IUserService:
    return get_user_service(db)


ServiceDep = Annotated[IUserService, Depends(_service)]


# Static routes must precede /{user_id}, otherwise FastAPI treats names such as
# "me" and "signup" as integer path parameters.
@router.post("/signup", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
def signup(request: SignUpRequest, service: ServiceDep, db: DbDep):
    user = service.signup(request)
    db.commit()
    return user


@protected_router.post("/me/password", status_code=status.HTTP_204_NO_CONTENT)
def change_password(request: ChangePasswordRequest, user_id: CurrentUserId, service: ServiceDep, db: DbDep):
    service.change_password(user_id, request)
    db.commit()
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@router.post("/password/reset", status_code=status.HTTP_202_ACCEPTED)
def reset_password(request: ResetPasswordRequest, service: ServiceDep, db: DbDep):
    service.request_password_reset(request)
    db.commit()
    return {"message": "If the account exists, a password reset email has been queued."}


@router.post("/password/confirm", status_code=status.HTTP_204_NO_CONTENT)
def confirm_password(request: ConfirmPasswordResetRequest, service: ServiceDep, db: DbDep):
    service.confirm_password_reset(request)
    db.commit()
    return Response(status_code=status.HTTP_204_NO_CONTENT)


@protected_router.get("/me/profile", response_model=UserResponse)
def get_profile(user_id: CurrentUserId, service: ServiceDep):
    return service.get_profile(user_id)


@protected_router.put("/me/profile", response_model=UserResponse)
def update_profile(request: UpdateProfileRequest, user_id: CurrentUserId, service: ServiceDep, db: DbDep):
    user = service.update_profile(user_id, request)
    db.commit()
    return user


@protected_router.get("/", response_model=PaginatedResponse[UserResponse])
def get_all(service: ServiceDep, page: int = 1, page_size: int = 20):
    return service.get_all(page=page, page_size=page_size)


@protected_router.post("/", response_model=UserResponse, status_code=status.HTTP_201_CREATED)
def create(request: CreateUserRequest, service: ServiceDep, db: DbDep,
           actor_user_id: CurrentUserId):
    user = service.create(request, actor_user_id)
    db.commit()
    return user


@protected_router.get("/{user_id}", response_model=UserResponse)
def get_by_id(user_id: int, service: ServiceDep):
    return service.get_by_id(user_id)


@protected_router.put("/{user_id}", response_model=UserResponse)
def update(user_id: int, request: UpdateUserRequest, service: ServiceDep, db: DbDep,
           actor_user_id: CurrentUserId):
    user = service.update(user_id, request, actor_user_id)
    db.commit()
    return user


@protected_router.patch("/{user_id}/activate", response_model=UserResponse)
def activate(user_id: int, service: ServiceDep, db: DbDep,
             actor_user_id: CurrentUserId):
    user = service.set_active(user_id, True, actor_user_id)
    db.commit()
    return user


@protected_router.patch("/{user_id}/deactivate", response_model=UserResponse)
def deactivate(user_id: int, service: ServiceDep, db: DbDep,
               actor_user_id: CurrentUserId):
    user = service.set_active(user_id, False, actor_user_id)
    db.commit()
    return user
