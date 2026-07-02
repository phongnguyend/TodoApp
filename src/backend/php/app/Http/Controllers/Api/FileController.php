<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\UploadFileRequest;
use App\Http\Resources\FileResource;
use App\Services\Contracts\FileServiceInterface;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;
use Symfony\Component\HttpFoundation\BinaryFileResponse;

/**
 * REST controller for uploaded files.
 * Analogous to FilesController : ControllerBase in ASP.NET Core.
 * The service is resolved by Laravel's IoC container via constructor injection.
 */
class FileController extends Controller
{
    public function __construct(
        private readonly FileServiceInterface $service
    ) {}

    /**
     * GET /api/files
     */
    public function index(Request $request): AnonymousResourceCollection
    {
        $page    = (int) $request->query('page', 1);
        $perPage = (int) $request->query('page_size', config('app.default_page_size', 20));

        return FileResource::collection(
            $this->service->getAll($page, $perPage)
        );
    }

    /**
     * GET /api/files/{id}
     */
    public function show(int $id): FileResource
    {
        return new FileResource($this->service->getById($id));
    }

    /**
     * GET /api/files/{id}/download
     */
    public function download(int $id): BinaryFileResponse
    {
        $target = $this->service->getDownloadTarget($id);

        return response()->download($target['path'], $target['name'], [
            'Content-Type' => $target['content_type'],
        ]);
    }

    /**
     * POST /api/files
     */
    public function store(UploadFileRequest $request): JsonResponse
    {
        $file = $this->service->upload($request->file('file'));

        return (new FileResource($file))
            ->response()
            ->setStatusCode(201);
    }

    /**
     * DELETE /api/files/{id}
     */
    public function destroy(int $id): JsonResponse
    {
        $this->service->delete($id);

        return response()->json(null, 204);
    }
}
