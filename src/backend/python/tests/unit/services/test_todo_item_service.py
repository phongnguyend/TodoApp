from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException

from app.models.todo_item import TodoItem
from app.schemas.todo_item import (
    CreateTodoItemRequest,
    UpdateTodoItemRequest,
)
from app.services.todo_item_service import TodoItemService


# ── Helpers ───────────────────────────────────────────────────────────────────

def _make_todo(
    id: int = 1,
    title: str = "Test Todo",
    description: str | None = None,
    is_completed: bool = False,
) -> TodoItem:
    todo = TodoItem(title=title, description=description, is_completed=is_completed)
    todo.id = id
    todo.created_at = datetime(2024, 1, 1, 12, 0, 0, tzinfo=timezone.utc)
    todo.updated_at = None
    return todo


# ── get_all ───────────────────────────────────────────────────────────────────

class TestGetAll:
    def test_returns_items_and_total(self):
        repo = MagicMock()
        repo.get_all.return_value = ([_make_todo(1, "A"), _make_todo(2, "B")], 2)

        result = TodoItemService(repo).get_all(page=1, page_size=20)

        assert result.total == 2
        assert len(result.items) == 2
        assert result.items[0].title == "A"

    def test_calculates_skip_from_page(self):
        repo = MagicMock()
        repo.get_all.return_value = ([], 50)

        TodoItemService(repo).get_all(page=3, page_size=10)

        repo.get_all.assert_called_once_with(skip=20, limit=10)

    def test_total_pages_uses_ceiling_division(self):
        repo = MagicMock()
        repo.get_all.return_value = ([], 21)

        result = TodoItemService(repo).get_all(page=1, page_size=20)

        assert result.total_pages == 2

    def test_empty_result(self):
        repo = MagicMock()
        repo.get_all.return_value = ([], 0)

        result = TodoItemService(repo).get_all()

        assert result.total == 0
        assert result.items == []


# ── get_incomplete ────────────────────────────────────────────────────────────

class TestGetIncomplete:
    def test_returns_incomplete_items(self):
        repo = MagicMock()
        repo.get_incomplete.return_value = ([_make_todo(1, "Undone")], 1)

        result = TodoItemService(repo).get_incomplete(page=1, page_size=20)

        assert result.total == 1
        repo.get_incomplete.assert_called_once_with(skip=0, limit=20)

    def test_calculates_skip_from_page(self):
        repo = MagicMock()
        repo.get_incomplete.return_value = ([], 0)

        TodoItemService(repo).get_incomplete(page=2, page_size=5)

        repo.get_incomplete.assert_called_once_with(skip=5, limit=5)


# ── get_by_id ─────────────────────────────────────────────────────────────────

class TestGetById:
    def test_returns_item_when_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = _make_todo(1, "Found")

        result = TodoItemService(repo).get_by_id(1)

        assert result.id == 1
        assert result.title == "Found"

    def test_raises_404_when_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            TodoItemService(repo).get_by_id(99)

        assert exc_info.value.status_code == 404
        assert "99" in exc_info.value.detail


# ── create ────────────────────────────────────────────────────────────────────

class TestCreate:
    def test_creates_and_returns_item(self):
        repo = MagicMock()
        repo.add.return_value = _make_todo(1, "New Todo", "Desc")

        result = TodoItemService(repo).create(CreateTodoItemRequest(title="New Todo", description="Desc"))

        assert result.id == 1
        assert result.title == "New Todo"
        assert result.description == "Desc"
        repo.add.assert_called_once()

    def test_creates_without_description(self):
        repo = MagicMock()
        repo.add.return_value = _make_todo(1, "No Desc")

        result = TodoItemService(repo).create(CreateTodoItemRequest(title="No Desc"))

        assert result.description is None

    def test_passes_title_and_description_to_repo(self):
        repo = MagicMock()
        repo.add.return_value = _make_todo(1, "Buy milk", "Whole milk")

        TodoItemService(repo).create(CreateTodoItemRequest(title="Buy milk", description="Whole milk"))

        added: TodoItem = repo.add.call_args[0][0]
        assert added.title == "Buy milk"
        assert added.description == "Whole milk"


# ── update ────────────────────────────────────────────────────────────────────

class TestUpdate:
    def test_updates_title(self):
        repo = MagicMock()
        todo = _make_todo(1, "Old Title")
        repo.get_by_id.return_value = todo
        repo.update.return_value = _make_todo(1, "New Title")

        result = TodoItemService(repo).update(1, UpdateTodoItemRequest(title="New Title"))

        assert result.title == "New Title"
        assert todo.title == "New Title"

    def test_updates_description(self):
        repo = MagicMock()
        todo = _make_todo(1, "Title", "Old Desc")
        repo.get_by_id.return_value = todo
        repo.update.return_value = todo

        TodoItemService(repo).update(1, UpdateTodoItemRequest(description="New Desc"))

        assert todo.description == "New Desc"

    def test_updates_is_completed(self):
        repo = MagicMock()
        todo = _make_todo(1, "Title", is_completed=False)
        repo.get_by_id.return_value = todo
        repo.update.return_value = todo

        TodoItemService(repo).update(1, UpdateTodoItemRequest(is_completed=True))

        assert todo.is_completed is True

    def test_skips_none_fields(self):
        repo = MagicMock()
        todo = _make_todo(1, "Original", "Original Desc", False)
        repo.get_by_id.return_value = todo
        repo.update.return_value = todo

        TodoItemService(repo).update(1, UpdateTodoItemRequest())

        assert todo.title == "Original"
        assert todo.description == "Original Desc"
        assert todo.is_completed is False

    def test_raises_404_when_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            TodoItemService(repo).update(99, UpdateTodoItemRequest(title="X"))

        assert exc_info.value.status_code == 404


# ── delete ────────────────────────────────────────────────────────────────────

class TestDelete:
    def test_deletes_item(self):
        repo = MagicMock()
        todo = _make_todo(1)
        repo.get_by_id.return_value = todo

        TodoItemService(repo).delete(1)

        repo.delete.assert_called_once_with(todo)

    def test_raises_404_when_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            TodoItemService(repo).delete(99)

        assert exc_info.value.status_code == 404


# ── mark_complete ─────────────────────────────────────────────────────────────

class TestMarkComplete:
    def test_sets_is_completed_to_true(self):
        repo = MagicMock()
        todo = _make_todo(1, "Pending", is_completed=False)
        repo.get_by_id.return_value = todo
        repo.update.return_value = _make_todo(1, "Pending", is_completed=True)

        result = TodoItemService(repo).mark_complete(1)

        assert todo.is_completed is True
        assert result.is_completed is True
        repo.update.assert_called_once_with(todo)

    def test_raises_404_when_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            TodoItemService(repo).mark_complete(99)

        assert exc_info.value.status_code == 404
