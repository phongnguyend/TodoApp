<?php

namespace App\Repositories;

use App\Models\TodoItemAttachment;
use App\Repositories\Contracts\TodoItemAttachmentRepositoryInterface;
use Illuminate\Support\Collection;

class TodoItemAttachmentRepository extends BaseRepository implements TodoItemAttachmentRepositoryInterface
{
    public function __construct(TodoItemAttachment $model)
    {
        parent::__construct($model);
    }

    public function findForTodo(int $todoItemId, int $attachmentId): ?TodoItemAttachment
    {
        return $this->model->newQuery()->where('todo_item_id', $todoItemId)->find($attachmentId);
    }

    public function findForTodoAndFile(int $todoItemId, int $fileId): ?TodoItemAttachment
    {
        return $this->model->newQuery()
            ->where('todo_item_id', $todoItemId)
            ->where('file_id', $fileId)
            ->first();
    }

    public function getForTodo(int $todoItemId): Collection
    {
        return $this->model->newQuery()
            ->where('todo_item_id', $todoItemId)
            ->orderBy('id')
            ->get();
    }
}
