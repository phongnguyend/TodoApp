<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\SaveTodoItemAttachmentRequest;
use App\Http\Resources\TodoItemAttachmentResource;
use App\Services\Contracts\TodoItemAttachmentServiceInterface;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

class TodoItemAttachmentController extends Controller
{
    public function __construct(private readonly TodoItemAttachmentServiceInterface $service) {}

    public function index(int $id): AnonymousResourceCollection
    {
        return TodoItemAttachmentResource::collection($this->service->getAll($id));
    }

    public function store(SaveTodoItemAttachmentRequest $request, int $id): JsonResponse
    {
        return (new TodoItemAttachmentResource($this->service->create(
            $id, (int) $request->validated('file_id'), $this->actorUserId()
        )))
            ->response()->setStatusCode(201);
    }

    public function show(int $id, int $attachmentId): TodoItemAttachmentResource
    {
        return new TodoItemAttachmentResource($this->service->getById($id, $attachmentId));
    }

    public function update(SaveTodoItemAttachmentRequest $request, int $id, int $attachmentId): TodoItemAttachmentResource
    {
        return new TodoItemAttachmentResource(
            $this->service->update($id, $attachmentId, (int) $request->validated('file_id'), $this->actorUserId())
        );
    }

    public function destroy(int $id, int $attachmentId): JsonResponse
    {
        $this->service->delete($id, $attachmentId);
        return response()->json(null, 204);
    }

    private function actorUserId(): ?int
    {
        $value = request()->attributes->get('authenticated_user_id');
        return is_int($value) && $value > 0 ? $value : null;
    }
}
