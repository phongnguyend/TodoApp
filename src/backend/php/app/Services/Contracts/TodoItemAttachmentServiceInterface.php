<?php

namespace App\Services\Contracts;

use App\Models\TodoItemAttachment;
use Illuminate\Support\Collection;

interface TodoItemAttachmentServiceInterface
{
    public function getAll(int $todoItemId): Collection;
    public function getById(int $todoItemId, int $attachmentId): TodoItemAttachment;
    public function create(int $todoItemId, int $fileId): TodoItemAttachment;
    public function update(int $todoItemId, int $attachmentId, int $fileId): TodoItemAttachment;
    public function delete(int $todoItemId, int $attachmentId): void;
}
