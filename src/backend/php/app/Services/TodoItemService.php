<?php

namespace App\Services;

use App\Http\Requests\CreateTodoItemRequest;
use App\Http\Requests\UpdateTodoItemRequest;
use App\Models\TodoItem;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Services\Contracts\TodoItemServiceInterface;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Business-logic layer.
 * Analogous to a scoped service registered via builder.Services.AddScoped<ITodoItemService, TodoItemService>().
 */
class TodoItemService implements TodoItemServiceInterface
{
    public function __construct(
        private readonly TodoItemRepositoryInterface $repository
    ) {}

    // ── Queries ───────────────────────────────────────────────────────────────

    public function getAll(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->repository->paginate($page, $perPage);
    }

    public function getIncomplete(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->repository->paginateIncomplete($page, $perPage);
    }

    public function getById(int $id): TodoItem
    {
        $todo = $this->repository->findById($id);

        if ($todo === null) {
            throw new ModelNotFoundException("Todo item {$id} not found.");
        }

        return $todo;
    }

    // ── Commands ──────────────────────────────────────────────────────────────

    public function create(CreateTodoItemRequest $request): TodoItem
    {
        return $this->repository->create($request->validated());
    }

    public function update(int $id, UpdateTodoItemRequest $request): TodoItem
    {
        $todo = $this->getById($id);
        return $this->repository->update($todo, $request->validated());
    }

    public function delete(int $id): void
    {
        $todo = $this->getById($id);
        $this->repository->delete($todo);
    }

    public function markComplete(int $id): TodoItem
    {
        $todo = $this->getById($id);
        return $this->repository->update($todo, ['is_completed' => true]);
    }
}
