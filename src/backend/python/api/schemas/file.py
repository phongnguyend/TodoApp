from datetime import datetime

from pydantic import BaseModel


# ── Response DTO (analogous to a FileDto / view model) ─────────────────────────
# Note: the on-disk `location` is intentionally not exposed to clients; content is
# retrieved via the dedicated download endpoint instead.

class FileResponse(BaseModel):
    model_config = {"from_attributes": True}

    id: int
    name: str
    extension: str
    size: int
    content_type: str | None
    created_at: datetime
    created_by_user_id: int | None = None
    updated_at: datetime | None
    updated_by_user_id: int | None = None
