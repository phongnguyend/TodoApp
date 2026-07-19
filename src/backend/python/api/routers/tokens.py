from typing import Annotated

from fastapi import APIRouter, Depends, Response
from sqlalchemy.orm import Session

from api.repositories.user_repository import UserRepository
from api.schemas.user import TokenRequest, TokenResponse
from api.services.user_service import UserService
from shared.database import get_db

router = APIRouter(prefix="/api/tokens", tags=["Tokens"])
DbDep = Annotated[Session, Depends(get_db)]


def _service(db: DbDep) -> UserService:
    return UserService(UserRepository(db))


ServiceDep = Annotated[UserService, Depends(_service)]


@router.post("", response_model=TokenResponse)
def create_token(request: TokenRequest, response: Response, service: ServiceDep) -> TokenResponse:
    token = service.create_token(request)
    response.headers["Cache-Control"] = "no-store"
    response.headers["Pragma"] = "no-cache"
    return token
