<?php

namespace App\Repositories;

use App\Models\TodoItem;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Eloquent-backed todo-item repository.
 * Analogous to EF Core TodoItemRepository : BaseRepository<TodoItem>.
 */
class TodoItemRepository extends BaseRepository implements TodoItemRepositoryInterface
{
    public function __construct(TodoItem $model)
    {
        parent::__construct($model);
    }

    public function paginateIncomplete(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return TodoItem::where('is_completed', false)
            ->orderByDesc('created_at')
            ->paginate($perPage, ['*'], 'page', $page);
    }

    public function findByTitle(string $title): ?TodoItem
    {
        return TodoItem::where('title', $title)->first();
    }
}
