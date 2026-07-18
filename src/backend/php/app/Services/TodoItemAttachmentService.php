<?php

namespace App\Services;

use App\Models\TodoItemAttachment;
use App\Repositories\Contracts\FileRepositoryInterface;
use App\Repositories\Contracts\TodoItemAttachmentRepositoryInterface;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Services\Contracts\TodoItemAttachmentServiceInterface;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Collection;

class TodoItemAttachmentService implements TodoItemAttachmentServiceInterface
{
    public function __construct(
        private readonly TodoItemAttachmentRepositoryInterface $attachments,
        private readonly TodoItemRepositoryInterface $todos,
        private readonly FileRepositoryInterface $files,
    ) {}

    private function requireTodo(int $id): void
    {
        if ($this->todos->findById($id) === null) {
            throw new ModelNotFoundException("Todo item {$id} not found.");
        }
    }

    private function requireFile(int $id): void
    {
        if ($this->files->findById($id) === null) {
            throw new ModelNotFoundException("File {$id} not found.");
        }
    }

    private function requireAttachment(int $todoItemId, int $attachmentId): TodoItemAttachment
    {
        $attachment = $this->attachments->findForTodo($todoItemId, $attachmentId);
        if ($attachment === null) {
            throw new ModelNotFoundException("Attachment {$attachmentId} not found for todo item {$todoItemId}.");
        }
        return $attachment;
    }

    public function getAll(int $todoItemId): Collection
    {
        $this->requireTodo($todoItemId);
        return $this->attachments->getForTodo($todoItemId);
    }

    public function getById(int $todoItemId, int $attachmentId): TodoItemAttachment
    {
        $this->requireTodo($todoItemId);
        return $this->requireAttachment($todoItemId, $attachmentId);
    }

    public function create(int $todoItemId, int $fileId): TodoItemAttachment
    {
        $this->requireTodo($todoItemId);
        $this->requireFile($fileId);
        return $this->attachments->findForTodoAndFile($todoItemId, $fileId)
            ?? $this->attachments->create(['todo_item_id' => $todoItemId, 'file_id' => $fileId]);
    }

    public function update(int $todoItemId, int $attachmentId, int $fileId): TodoItemAttachment
    {
        $this->requireTodo($todoItemId);
        $this->requireFile($fileId);
        $attachment = $this->requireAttachment($todoItemId, $attachmentId);
        $existing = $this->attachments->findForTodoAndFile($todoItemId, $fileId);
        if ($existing !== null && $existing->id !== $attachment->id) {
            return $existing;
        }
        return $this->attachments->update($attachment, ['file_id' => $fileId]);
    }

    public function delete(int $todoItemId, int $attachmentId): void
    {
        $this->requireTodo($todoItemId);
        $this->attachments->delete($this->requireAttachment($todoItemId, $attachmentId));
    }
}
