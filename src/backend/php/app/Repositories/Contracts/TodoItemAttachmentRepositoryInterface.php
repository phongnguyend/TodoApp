<?php

namespace App\Repositories\Contracts;

use App\Models\TodoItemAttachment;
use Illuminate\Support\Collection;

interface TodoItemAttachmentRepositoryInterface extends RepositoryInterface
{
    public function findForTodo(int $todoItemId, int $attachmentId): ?TodoItemAttachment;
    public function findForTodoAndFile(int $todoItemId, int $fileId): ?TodoItemAttachment;
    public function getForTodo(int $todoItemId): Collection;
}
