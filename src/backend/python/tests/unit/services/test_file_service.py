from datetime import datetime, timezone
from io import BytesIO
from unittest.mock import MagicMock

import pytest
from fastapi import HTTPException, UploadFile
from starlette.datastructures import Headers

from shared.config import settings
from shared.models.file import File
from api.services.file_service import FileService


# ── Helpers ───────────────────────────────────────────────────────────────────

def _make_file(
    id: int = 1,
    name: str = "test.txt",
    extension: str = "txt",
    size: int = 100,
    content_type: str | None = "text/plain",
    location: str = "/tmp/test.txt",
) -> File:
    file = File(name=name, extension=extension, size=size, content_type=content_type, location=location)
    file.id = id
    file.created_at = datetime(2024, 1, 1, 12, 0, 0, tzinfo=timezone.utc)
    file.updated_at = None
    return file


def _make_upload(
    filename: str = "photo.png", content: bytes = b"hello world", content_type: str | None = "image/png"
) -> UploadFile:
    headers = Headers({"content-type": content_type}) if content_type else None
    return UploadFile(file=BytesIO(content), filename=filename, headers=headers)


def _add_side_effect(entity: File) -> File:
    entity.id = 1
    entity.created_at = datetime(2024, 1, 1, 12, 0, 0, tzinfo=timezone.utc)
    entity.updated_at = None
    return entity


# ── get_all ───────────────────────────────────────────────────────────────────

class TestGetAll:
    def test_returns_items_and_total(self):
        repo = MagicMock()
        repo.get_all.return_value = ([_make_file(1, "a.txt"), _make_file(2, "b.txt")], 2)

        result = FileService(repo).get_all(page=1, page_size=20)

        assert result.total == 2
        assert len(result.items) == 2
        assert result.items[0].name == "a.txt"

    def test_calculates_skip_from_page(self):
        repo = MagicMock()
        repo.get_all.return_value = ([], 50)

        FileService(repo).get_all(page=3, page_size=10)

        repo.get_all.assert_called_once_with(skip=20, limit=10)

    def test_empty_result(self):
        repo = MagicMock()
        repo.get_all.return_value = ([], 0)

        result = FileService(repo).get_all()

        assert result.total == 0
        assert result.items == []


# ── get_by_id ─────────────────────────────────────────────────────────────────

class TestGetById:
    def test_returns_item_when_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = _make_file(1, "found.txt")

        result = FileService(repo).get_by_id(1)

        assert result.id == 1
        assert result.name == "found.txt"

    def test_raises_404_when_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            FileService(repo).get_by_id(99)

        assert exc_info.value.status_code == 404
        assert "99" in exc_info.value.detail


# ── upload ────────────────────────────────────────────────────────────────────

class TestUpload:
    def test_saves_file_to_storage_dir_and_returns_metadata(self, tmp_path, monkeypatch):
        monkeypatch.setattr(settings, "FILE_STORAGE_PATH", str(tmp_path))
        repo = MagicMock()
        repo.add.side_effect = _add_side_effect

        upload = _make_upload("photo.png", b"binary-data", "image/png")

        result = FileService(repo).upload(upload)

        assert result.id == 1
        assert result.name == "photo.png"
        assert result.extension == "png"
        assert result.size == len(b"binary-data")
        assert result.content_type == "image/png"
        assert len(list(tmp_path.iterdir())) == 1

    def test_rejects_files_exceeding_max_size(self, tmp_path, monkeypatch):
        monkeypatch.setattr(settings, "FILE_STORAGE_PATH", str(tmp_path))
        monkeypatch.setattr(settings, "MAX_UPLOAD_SIZE_BYTES", 5)
        repo = MagicMock()

        upload = _make_upload("big.bin", b"0123456789")

        with pytest.raises(HTTPException) as exc_info:
            FileService(repo).upload(upload)

        assert exc_info.value.status_code == 413
        repo.add.assert_not_called()

    def test_sanitizes_path_traversal_in_filename(self, tmp_path, monkeypatch):
        monkeypatch.setattr(settings, "FILE_STORAGE_PATH", str(tmp_path))
        repo = MagicMock()
        repo.add.side_effect = _add_side_effect

        upload = _make_upload("../../etc/passwd", b"data")

        result = FileService(repo).upload(upload)

        assert result.name == "passwd"
        saved_files = list(tmp_path.iterdir())
        assert len(saved_files) == 1
        assert ".." not in str(saved_files[0])


# ── get_download_target ───────────────────────────────────────────────────────

class TestGetDownloadTarget:
    def test_returns_path_name_and_content_type(self, tmp_path):
        repo = MagicMock()
        file_path = tmp_path / "data.txt"
        file_path.write_text("hello")
        repo.get_by_id.return_value = _make_file(1, "data.txt", "txt", 5, "text/plain", str(file_path))

        path, name, content_type = FileService(repo).get_download_target(1)

        assert path == str(file_path)
        assert name == "data.txt"
        assert content_type == "text/plain"

    def test_defaults_content_type_when_missing(self, tmp_path):
        repo = MagicMock()
        file_path = tmp_path / "data.bin"
        file_path.write_text("hello")
        repo.get_by_id.return_value = _make_file(1, "data.bin", "bin", 5, None, str(file_path))

        _, _, content_type = FileService(repo).get_download_target(1)

        assert content_type == "application/octet-stream"

    def test_raises_404_when_metadata_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            FileService(repo).get_download_target(99)

        assert exc_info.value.status_code == 404

    def test_raises_404_when_content_missing_from_disk(self, tmp_path):
        repo = MagicMock()
        missing_path = tmp_path / "missing.txt"
        repo.get_by_id.return_value = _make_file(1, "missing.txt", "txt", 5, "text/plain", str(missing_path))

        with pytest.raises(HTTPException) as exc_info:
            FileService(repo).get_download_target(1)

        assert exc_info.value.status_code == 404


# ── delete ────────────────────────────────────────────────────────────────────

class TestDelete:
    def test_deletes_item_and_removes_file_from_disk(self, tmp_path):
        repo = MagicMock()
        file_path = tmp_path / "to_delete.txt"
        file_path.write_text("bye")
        file = _make_file(1, "to_delete.txt", "txt", 3, "text/plain", str(file_path))
        repo.get_by_id.return_value = file

        FileService(repo).delete(1)

        repo.delete.assert_called_once_with(file)
        assert not file_path.exists()

    def test_raises_404_when_not_found(self):
        repo = MagicMock()
        repo.get_by_id.return_value = None

        with pytest.raises(HTTPException) as exc_info:
            FileService(repo).delete(99)

        assert exc_info.value.status_code == 404
