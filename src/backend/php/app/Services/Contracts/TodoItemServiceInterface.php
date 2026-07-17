<?php

namespace App\Services\Contracts;

use App\Http\Requests\CreateTodoItemRequest;
use App\Http\Requests\UpdateTodoItemRequest;
use App\Models\TodoItem;
use Illuminate\Http\UploadedFile;
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

    /**
     * Parses the given CSV file and creates a todo item for each valid row.
     *
     * @return array{imported: int, failed: int, errors: array<int, array{row: int, error: string}>}
     */
    public function importCsv(UploadedFile $file): array;

    /**
     * Renders every todo item as CSV text (header row + one row per item).
     */
    public function exportCsv(): string;

    /**
     * Parses the given Excel (.xlsx/.xls) file and creates a todo item for each valid row.
     *
     * @return array{imported: int, failed: int, errors: array<int, array{row: int, error: string}>}
     */
    public function importExcel(UploadedFile $file): array;

    /**
     * Renders every todo item as an Excel (.xlsx) workbook (header row + one row per item).
     */
    public function exportExcel(): string;
}
