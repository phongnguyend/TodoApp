from datetime import datetime

from pydantic import BaseModel


class CreateTodoItemAttachmentRequest(BaseModel):
    file_id: int


class TodoItemAttachmentResponse(BaseModel):
    model_config = {"from_attributes": True}

    id: int
    todo_item_id: int
    file_id: int
    created_at: datetime
    updated_at: datetime | None

