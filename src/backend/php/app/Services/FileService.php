<?php

namespace App\Services;

use App\Models\File;
use App\Repositories\Contracts\FileRepositoryInterface;
use App\Services\Contracts\FileServiceInterface;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Http\UploadedFile;
use Illuminate\Pagination\LengthAwarePaginator;

/**
 * Business-logic layer.
 * Analogous to a scoped service registered via builder.Services.AddScoped<IFileService, FileService>().
 */
class FileService implements FileServiceInterface
{
    public function __construct(
        private readonly FileRepositoryInterface $repository
    ) {}

    // ── Helpers ───────────────────────────────────────────────────────────────

    private function getOrFail(int $id): File
    {
        $file = $this->repository->findById($id);

        if ($file === null) {
            throw new ModelNotFoundException("File {$id} not found.");
        }

        return $file;
    }

    private function storageDirectory(): string
    {
        return (string) env('FILE_STORAGE_PATH', storage_path('app/uploads'));
    }

    // ── Queries ───────────────────────────────────────────────────────────────

    public function getAll(int $page = 1, int $perPage = 20): LengthAwarePaginator
    {
        return $this->repository->paginate($page, $perPage);
    }

    public function getById(int $id): File
    {
        return $this->getOrFail($id);
    }

    // ── Commands ──────────────────────────────────────────────────────────────

    public function upload(UploadedFile $uploadedFile, ?int $actorUserId = null): File
    {
        // Strip any directory components from the client-supplied name to prevent path traversal.
        $originalName = basename($uploadedFile->getClientOriginalName());
        $extension = ltrim(
            strtolower($uploadedFile->getClientOriginalExtension() ?: pathinfo($originalName, PATHINFO_EXTENSION)),
            '.'
        );

        $storageDir = $this->storageDirectory();
        if (!is_dir($storageDir)) {
            mkdir($storageDir, 0755, true);
        }

        // A random prefix avoids collisions/overwrites between uploads that share a name.
        $storedName = uniqid('', true) . '_' . $originalName;
        $movedFile = $uploadedFile->move($storageDir, $storedName);

        $data = [
            'name'         => $originalName,
            'extension'    => $extension,
            'size'         => $movedFile->getSize(),
            'content_type' => $uploadedFile->getClientMimeType(),
            'location'     => $movedFile->getPathname(),
        ];
        if ($actorUserId !== null) $data['created_by_user_id'] = $actorUserId;
        return $this->repository->create($data);
    }

    public function getDownloadTarget(int $id): array
    {
        $file = $this->getOrFail($id);

        if (!is_file($file->location)) {
            throw new ModelNotFoundException("File {$id} content not found on disk.");
        }

        return [
            'path'         => $file->location,
            'name'         => $file->name,
            'content_type' => $file->content_type ?: 'application/octet-stream',
        ];
    }

    public function delete(int $id): void
    {
        $file = $this->getOrFail($id);
        $this->repository->delete($file);

        if (is_file($file->location)) {
            unlink($file->location);
        }
    }
}
