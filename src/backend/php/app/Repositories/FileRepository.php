<?php

namespace App\Repositories;

use App\Models\File;
use App\Repositories\Contracts\FileRepositoryInterface;

/**
 * Eloquent-backed file repository.
 * Analogous to EF Core FileRepository : BaseRepository<File>.
 */
class FileRepository extends BaseRepository implements FileRepositoryInterface
{
    public function __construct(File $model)
    {
        parent::__construct($model);
    }
}
