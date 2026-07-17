<?php

namespace App\Repositories\Contracts;

use App\Models\TodoItem;
use Illuminate\Database\Eloquent\Collection;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Todo-item specific repository contract.
 * Analogous to ITodoItemRepository : IRepository<TodoItem> in C#.
 */
interface TodoItemRepositoryInterface extends RepositoryInterface
{
    public function paginateIncomplete(int $page = 1, int $perPage = 20): LengthAwarePaginator;

    public function findByTitle(string $title): ?TodoItem;

    /**
     * Returns every todo item, ordered like {@see paginate()}. Used for CSV export.
     */
    public function getAllOrdered(): Collection;
}
