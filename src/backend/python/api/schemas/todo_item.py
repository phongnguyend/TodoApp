from datetime import datetime

from pydantic import BaseModel, Field


# ── Create DTO (analogous to a CreateTodoItemRequest / command model) ──────────

class CreateTodoItemRequest(BaseModel):
    title: str = Field(..., min_length=1, max_length=200, examples=["Buy groceries"])
    description: str | None = Field(default=None, max_length=2000, examples=["Milk, eggs, bread"])


# ── Update DTO ─────────────────────────────────────────────────────────────────

class UpdateTodoItemRequest(BaseModel):
    title: str | None = Field(default=None, min_length=1, max_length=200, examples=["Buy groceries"])
    description: str | None = Field(default=None, max_length=2000, examples=["Milk, eggs, bread"])
    is_completed: bool | None = Field(default=None)


# ── Response DTO (analogous to a TodoItemDto / view model) ────────────────────

class TodoItemResponse(BaseModel):
    model_config = {"from_attributes": True}  # like AutoMapper / ORM-mapped DTO

    id: int
    title: str
    description: str | None
    is_completed: bool
    created_at: datetime
    created_by_user_id: int | None = None
    updated_at: datetime | None
    updated_by_user_id: int | None = None


# ── Paginated list response ────────────────────────────────────────────────────

class PaginatedResponse[T](BaseModel):
    items: list[T]
    total: int
    page: int
    page_size: int
    total_pages: int


# ── CSV import/export DTOs ─────────────────────────────────────────────────────

class ImportRowError(BaseModel):
    row: int
    error: str


class ImportResult(BaseModel):
    imported: int
    failed: int
    errors: list[ImportRowError] = Field(default_factory=list)
