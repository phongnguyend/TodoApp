<?php

namespace App\Services\Contracts;

use App\Models\File;
use Illuminate\Http\UploadedFile;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Service interface - mirrors IFileService in C#.
 */
interface FileServiceInterface
{
    public function getAll(int $page, int $perPage): LengthAwarePaginator;

    public function getById(int $id): File;

    public function upload(UploadedFile $uploadedFile): File;

    /**
     * Returns ['path' => string, 'name' => string, 'content_type' => string] for streaming
     * the file's content back to the client.
     *
     * @return array{path: string, name: string, content_type: string}
     */
    public function getDownloadTarget(int $id): array;

    public function delete(int $id): void;
}
