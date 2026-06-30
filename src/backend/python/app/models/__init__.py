# Import all models here so that alembic autogenerate can detect them via Base.metadata
from app.models.email_log import EmailLog  # noqa: F401
from app.models.todo_item import TodoItem  # noqa: F401

__all__ = ["TodoItem", "EmailLog"]
