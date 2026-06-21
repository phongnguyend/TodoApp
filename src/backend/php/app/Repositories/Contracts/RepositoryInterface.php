<?php

namespace App\Repositories\Contracts;

use Illuminate\Database\Eloquent\Model;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Generic repository contract — mirrors IRepository<T> in C#.
 */
interface RepositoryInterface
{
    public function findById(int $id): ?Model;

    public function paginate(int $page = 1, int $perPage = 20): LengthAwarePaginator;

    public function create(array $data): Model;

    public function update(Model $model, array $data): Model;

    public function delete(Model $model): void;
}
