"""add user audit columns

Revision ID: eeff00556677
Revises: ddeeff445566
Create Date: 2026-07-20
"""

from collections.abc import Sequence

from alembic import op
import sqlalchemy as sa


revision: str = "eeff00556677"
down_revision: str | None = "ddeeff445566"
branch_labels: str | Sequence[str] | None = None
depends_on: str | Sequence[str] | None = None


def upgrade() -> None:
    for table_name in ("todo_items", "todo_item_attachments", "email_logs", "files", "users"):
        with op.batch_alter_table(table_name) as batch_op:
            batch_op.add_column(sa.Column("created_by_user_id", sa.Integer(), nullable=True))
            batch_op.add_column(sa.Column("updated_by_user_id", sa.Integer(), nullable=True))
            batch_op.create_foreign_key(
                f"fk_{table_name}_created_by_user_id_users",
                "users",
                ["created_by_user_id"],
                ["id"],
                ondelete="SET NULL",
            )
            batch_op.create_foreign_key(
                f"fk_{table_name}_updated_by_user_id_users",
                "users",
                ["updated_by_user_id"],
                ["id"],
                ondelete="SET NULL",
            )


def downgrade() -> None:
    for table_name in reversed(("todo_items", "todo_item_attachments", "email_logs", "files", "users")):
        with op.batch_alter_table(table_name) as batch_op:
            batch_op.drop_constraint(f"fk_{table_name}_updated_by_user_id_users", type_="foreignkey")
            batch_op.drop_constraint(f"fk_{table_name}_created_by_user_id_users", type_="foreignkey")
            batch_op.drop_column("updated_by_user_id")
            batch_op.drop_column("created_by_user_id")
