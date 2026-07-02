<?php

namespace App\Services\Contracts;

use App\Http\Requests\CreateTodoItemRequest;
use App\Http\Requests\UpdateTodoItemRequest;
use App\Models\TodoItem;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Service interface - mirrors ITodoItemService in C#.
 */
interface TodoItemServiceInterface
{
    public function getAll(int $page, int $perPage): LengthAwarePaginator;

    public function getIncomplete(int $page, int $perPage): LengthAwarePaginator;

    public function getById(int $id): TodoItem;

    public function create(CreateTodoItemRequest $request): TodoItem;

    public function update(int $id, UpdateTodoItemRequest $request): TodoItem;

    public function delete(int $id): void;

    public function markComplete(int $id): TodoItem;
}
