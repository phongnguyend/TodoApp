<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\Controller;
use App\Http\Requests\CreateTodoItemRequest;
use App\Http\Requests\UpdateTodoItemRequest;
use App\Http\Resources\TodoItemResource;
use App\Services\Contracts\TodoItemServiceInterface;
use Illuminate\Http\JsonResponse;
use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\AnonymousResourceCollection;

/**
 * REST controller for todo items.
 * Analogous to TodoItemsController : ControllerBase in ASP.NET Core.
 * The service is resolved by Laravel's IoC container via constructor injection.
 */
class TodoItemController extends Controller
{
    public function __construct(
        private readonly TodoItemServiceInterface $service
    ) {}

    /**
     * GET /api/todo-items
     */
    public function index(Request $request): AnonymousResourceCollection
    {
        $page    = (int) $request->query('page', 1);
        $perPage = (int) $request->query('page_size', config('app.default_page_size', 20));

        return TodoItemResource::collection(
            $this->service->getAll($page, $perPage)
        );
    }

    /**
     * GET /api/todo-items/incomplete
     */
    public function incomplete(Request $request): AnonymousResourceCollection
    {
        $page    = (int) $request->query('page', 1);
        $perPage = (int) $request->query('page_size', config('app.default_page_size', 20));

        return TodoItemResource::collection(
            $this->service->getIncomplete($page, $perPage)
        );
    }

    /**
     * GET /api/todo-items/{id}
     */
    public function show(int $id): TodoItemResource
    {
        return new TodoItemResource($this->service->getById($id));
    }

    /**
     * POST /api/todo-items
     */
    public function store(CreateTodoItemRequest $request): JsonResponse
    {
        $todo = $this->service->create($request);

        return (new TodoItemResource($todo))
            ->response()
            ->setStatusCode(201);
    }

    /**
     * PUT /api/todo-items/{id}
     */
    public function update(UpdateTodoItemRequest $request, int $id): TodoItemResource
    {
        return new TodoItemResource($this->service->update($id, $request));
    }

    /**
     * PATCH /api/todo-items/{id}/complete
     */
    public function complete(int $id): TodoItemResource
    {
        return new TodoItemResource($this->service->markComplete($id));
    }

    /**
     * DELETE /api/todo-items/{id}
     */
    public function destroy(int $id): JsonResponse
    {
        $this->service->delete($id);

        return response()->json(null, 204);
    }
}
