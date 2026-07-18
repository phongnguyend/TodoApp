"""add_todo_item_attachments

Revision ID: ccddeeff3344
Revises: bbccddee2233
"""

from collections.abc import Sequence

import sqlalchemy as sa
from alembic import op

revision: str = "ccddeeff3344"
down_revision: str | None = "bbccddee2233"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    op.create_table(
        "todo_item_attachments",
        sa.Column("id", sa.Integer(), autoincrement=True, nullable=False),
        sa.Column("todo_item_id", sa.Integer(), nullable=False),
        sa.Column("file_id", sa.Integer(), nullable=False),
        sa.Column("created_at", sa.DateTime(timezone=True), server_default=sa.func.now(), nullable=False),
        sa.Column("updated_at", sa.DateTime(timezone=True), nullable=True),
        sa.ForeignKeyConstraint(["file_id"], ["files.id"], ondelete="CASCADE"),
        sa.ForeignKeyConstraint(["todo_item_id"], ["todo_items.id"], ondelete="CASCADE"),
        sa.PrimaryKeyConstraint("id"),
        sa.UniqueConstraint("todo_item_id", "file_id", name="uq_todo_item_attachments_todo_item_file"),
    )
    op.create_index(op.f("ix_todo_item_attachments_id"), "todo_item_attachments", ["id"], unique=False)
    op.create_index(op.f("ix_todo_item_attachments_file_id"), "todo_item_attachments", ["file_id"], unique=False)
    op.create_index(
        op.f("ix_todo_item_attachments_todo_item_id"), "todo_item_attachments", ["todo_item_id"], unique=False
    )


def downgrade() -> None:
    op.drop_index(op.f("ix_todo_item_attachments_todo_item_id"), table_name="todo_item_attachments")
    op.drop_index(op.f("ix_todo_item_attachments_file_id"), table_name="todo_item_attachments")
    op.drop_index(op.f("ix_todo_item_attachments_id"), table_name="todo_item_attachments")
    op.drop_table("todo_item_attachments")

