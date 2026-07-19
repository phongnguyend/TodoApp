from abc import abstractmethod

from sqlalchemy import func
from sqlalchemy.orm import Session

from shared.models.email_log import EmailLog
from shared.models.user import User
from api.repositories.base_repository import BaseRepository, IRepository


class IUserRepository(IRepository[User]):
    @abstractmethod
    def get_by_email(self, email: str) -> User | None: ...

    @abstractmethod
    def username_exists(self, username: str, excluding_id: int | None = None) -> bool: ...

    @abstractmethod
    def email_exists(self, email: str, excluding_id: int | None = None) -> bool: ...

    @abstractmethod
    def add_email_log(self, email_log: EmailLog) -> EmailLog: ...


class UserRepository(BaseRepository[User], IUserRepository):
    def __init__(self, db: Session) -> None:
        super().__init__(db, User)
        self._db = db

    def get_by_email(self, email: str) -> User | None:
        return self._db.query(User).filter(func.lower(User.email) == email.lower()).first()

    def username_exists(self, username: str, excluding_id: int | None = None) -> bool:
        query = self._db.query(User).filter(func.lower(User.username) == username.lower())
        if excluding_id is not None:
            query = query.filter(User.id != excluding_id)
        return self._db.query(query.exists()).scalar() is True

    def email_exists(self, email: str, excluding_id: int | None = None) -> bool:
        query = self._db.query(User).filter(func.lower(User.email) == email.lower())
        if excluding_id is not None:
            query = query.filter(User.id != excluding_id)
        return self._db.query(query.exists()).scalar() is True

    def add_email_log(self, email_log: EmailLog) -> EmailLog:
        self._db.add(email_log)
        self._db.flush()
        return email_log
