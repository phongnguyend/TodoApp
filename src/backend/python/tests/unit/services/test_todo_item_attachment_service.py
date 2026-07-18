from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException

from shared.models.file import File
from shared.models.todo_item import TodoItem
from shared.models.todo_item_attachment import TodoItemAttachment
from api.schemas.todo_item_attachment import CreateTodoItemAttachmentRequest
from api.services.todo_item_attachment_service import TodoItemAttachmentService


def _todo(todo_id: int = 10) -> TodoItem:
    item = TodoItem(title="Task")
    item.id = todo_id
    return item


def _file(file_id: int = 5) -> File:
    item = File(name="a.txt", extension="txt", size=1, location="/tmp/a.txt")
    item.id = file_id
    return item


def _attachment(attachment_id: int = 1, todo_id: int = 10, file_id: int = 5) -> TodoItemAttachment:
    item = TodoItemAttachment(todo_item_id=todo_id, file_id=file_id)
    item.id = attachment_id
    item.created_at = datetime(2024, 1, 1, tzinfo=timezone.utc)
    item.updated_at = None
    return item


def _service() -> tuple[TodoItemAttachmentService, MagicMock, MagicMock, MagicMock]:
    attachment_repo = MagicMock()
    todo_repo = MagicMock()
    file_repo = MagicMock()
    return TodoItemAttachmentService(attachment_repo, todo_repo, file_repo), attachment_repo, todo_repo, file_repo


class TestGetAll:
    def test_returns_attachments_for_existing_todo(self):
        service, attachments, todos, _ = _service()
        todos.get_by_id.return_value = _todo()
        attachments.get_by_todo_item_id.return_value = [_attachment()]

        result = service.get_all(10)

        assert len(result) == 1
        assert result[0].todo_item_id == 10
        assert result[0].file_id == 5
        attachments.get_by_todo_item_id.assert_called_once_with(10)

    def test_raises_404_when_todo_does_not_exist(self):
        service, attachments, todos, _ = _service()
        todos.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            service.get_all(99)

        assert exc_info.value.status_code == 404
        attachments.get_by_todo_item_id.assert_not_called()


class TestGetById:
    def test_returns_attachment_scoped_to_todo(self):
        service, attachments, todos, _ = _service()
        todos.get_by_id.return_value = _todo()
        attachments.get_by_id_for_todo_item.return_value = _attachment(3)

        result = service.get_by_id(10, 3)

        assert result.id == 3
        attachments.get_by_id_for_todo_item.assert_called_once_with(10, 3)

    def test_raises_404_when_attachment_does_not_exist_for_todo(self):
        service, attachments, todos, _ = _service()
        todos.get_by_id.return_value = _todo()
        attachments.get_by_id_for_todo_item.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            service.get_by_id(10, 99)

        assert exc_info.value.status_code == 404
        assert "99" in exc_info.value.detail


class TestCreate:
    def test_creates_attachment_when_todo_and_file_exist(self):
        service, attachments, todos, files = _service()
        todos.get_by_id.return_value = _todo()
        files.get_by_id.return_value = _file()
        attachments.get_by_todo_item_and_file.return_value = None
        attachments.add.side_effect = lambda item: _attachment(1, item.todo_item_id, item.file_id)

        result = service.create(10, CreateTodoItemAttachmentRequest(file_id=5))

        assert result.todo_item_id == 10
        assert result.file_id == 5
        attachments.add.assert_called_once()

    def test_returns_existing_attachment_for_duplicate_link(self):
        service, attachments, todos, files = _service()
        todos.get_by_id.return_value = _todo()
        files.get_by_id.return_value = _file()
        attachments.get_by_todo_item_and_file.return_value = _attachment(7)

        result = service.create(10, CreateTodoItemAttachmentRequest(file_id=5))

        assert result.id == 7
        attachments.add.assert_not_called()

    def test_raises_404_when_file_does_not_exist(self):
        service, attachments, todos, files = _service()
        todos.get_by_id.return_value = _todo()
        files.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            service.create(10, CreateTodoItemAttachmentRequest(file_id=99))

        assert exc_info.value.status_code == 404
        attachments.add.assert_not_called()


class TestUpdate:
    def test_updates_attachment_file(self):
        service, attachments, todos, files = _service()
        current = _attachment(3)
        todos.get_by_id.return_value = _todo()
        files.get_by_id.return_value = _file(6)
        attachments.get_by_id_for_todo_item.return_value = current
        attachments.get_by_todo_item_and_file.return_value = None
        attachments.update.side_effect = lambda item: item

        result = service.update(10, 3, CreateTodoItemAttachmentRequest(file_id=6))

        assert result.file_id == 6
        attachments.update.assert_called_once_with(current)

    def test_returns_existing_link_instead_of_violating_unique_constraint(self):
        service, attachments, todos, files = _service()
        todos.get_by_id.return_value = _todo()
        files.get_by_id.return_value = _file(6)
        attachments.get_by_id_for_todo_item.return_value = _attachment(3)
        attachments.get_by_todo_item_and_file.return_value = _attachment(4, file_id=6)

        result = service.update(10, 3, CreateTodoItemAttachmentRequest(file_id=6))

        assert result.id == 4
        attachments.update.assert_not_called()


class TestDelete:
    def test_deletes_scoped_attachment(self):
        service, attachments, todos, _ = _service()
        item = _attachment(3)
        todos.get_by_id.return_value = _todo()
        attachments.get_by_id_for_todo_item.return_value = item

        service.delete(10, 3)

        attachments.delete.assert_called_once_with(item)

