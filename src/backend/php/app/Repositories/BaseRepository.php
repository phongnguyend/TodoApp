<?php

namespace App\Repositories;

use App\Repositories\Contracts\RepositoryInterface;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Generic base repository backed by Eloquent.
 * Analogous to GenericRepository<T> / EfRepository<T> in C#.
 */
abstract class BaseRepository implements RepositoryInterface
{
    public function __construct(protected readonly Model $model) {}

    public function findById(int $id): ?Model
    {
        return $this->model->find($id);
    }

    public function paginate(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->model->newQuery()
            ->orderByDesc('created_at')
            ->paginate($perPage, ['*'], 'page', $page);
    }

    public function create(array $data): Model
    {
        return $this->model->create($data);
    }

    public function update(Model $model, array $data): Model
    {
        $model->update($data);
        return $model->refresh();
    }

    public function delete(Model $model): void
    {
        $model->delete();
    }
}
