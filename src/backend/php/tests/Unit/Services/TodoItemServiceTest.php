<?php

namespace Tests\Unit\Services;

use App\Http\Requests\CreateTodoItemRequest;
use App\Http\Requests\UpdateTodoItemRequest;
use App\Models\TodoItem;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Services\TodoItemService;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Pagination\LengthAwarePaginator;
use Mockery;
use Tests\TestCase;

class TodoItemServiceTest extends TestCase
{
    private TodoItemRepositoryInterface $repository;
    private TodoItemService $service;

    protected function setUp(): void
    {
        parent::setUp();
        $this->repository = Mockery::mock(TodoItemRepositoryInterface::class);
        $this->service    = new TodoItemService($this->repository);
    }

    protected function tearDown(): void
    {
        Mockery::close();
        parent::tearDown();
    }

    // ── getAll ────────────────────────────────────────────────────────────────

    public function test_getAll_delegates_to_repository_paginate(): void
    {
        $paginator = Mockery::mock(LengthAwarePaginator::class);
        $this->repository->shouldReceive('paginate')->with(1, 20)->once()->andReturn($paginator);

        $result = $this->service->getAll(1, 20);

        $this->assertSame($paginator, $result);
    }

    // ── getIncomplete ─────────────────────────────────────────────────────────

    public function test_getIncomplete_delegates_to_repository_paginateIncomplete(): void
    {
        $paginator = Mockery::mock(LengthAwarePaginator::class);
        $this->repository->shouldReceive('paginateIncomplete')->with(2, 10)->once()->andReturn($paginator);

        $result = $this->service->getIncomplete(2, 10);

        $this->assertSame($paginator, $result);
    }

    // ── getById ───────────────────────────────────────────────────────────────

    public function test_getById_returns_todo_when_found(): void
    {
        $todo = new TodoItem(['title' => 'Test']);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($todo);

        $result = $this->service->getById(1);

        $this->assertSame($todo, $result);
    }

    public function test_getById_throws_ModelNotFoundException_when_not_found(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->getById(99);
    }

    // ── create ────────────────────────────────────────────────────────────────

    public function test_create_passes_validated_data_to_repository(): void
    {
        $data    = ['title' => 'New todo', 'description' => null];
        $todo    = new TodoItem($data);
        $request = Mockery::mock(CreateTodoItemRequest::class);
        $request->shouldReceive('validated')->once()->andReturn($data);
        $this->repository->shouldReceive('create')->with($data)->once()->andReturn($todo);

        $result = $this->service->create($request);

        $this->assertSame($todo, $result);
    }

    // ── update ────────────────────────────────────────────────────────────────

    public function test_update_finds_todo_and_delegates_to_repository(): void
    {
        $existing = new TodoItem(['title' => 'Old title']);
        $updated  = new TodoItem(['title' => 'New title']);
        $data     = ['title' => 'New title'];
        $request  = Mockery::mock(UpdateTodoItemRequest::class);
        $request->shouldReceive('validated')->once()->andReturn($data);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($existing);
        $this->repository->shouldReceive('update')->with($existing, $data)->once()->andReturn($updated);

        $result = $this->service->update(1, $request);

        $this->assertSame($updated, $result);
    }

    public function test_update_throws_ModelNotFoundException_when_not_found(): void
    {
        $request = Mockery::mock(UpdateTodoItemRequest::class);
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->update(99, $request);
    }

    // ── delete ────────────────────────────────────────────────────────────────

    public function test_delete_finds_and_removes_todo(): void
    {
        $todo = new TodoItem(['title' => 'To delete']);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($todo);
        $this->repository->shouldReceive('delete')->with($todo)->once();

        $this->service->delete(1);

        $this->expectNotToPerformAssertions();
    }

    public function test_delete_throws_ModelNotFoundException_when_not_found(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->delete(99);
    }

    // ── markComplete ──────────────────────────────────────────────────────────

    public function test_markComplete_sets_is_completed_to_true(): void
    {
        $todo      = new TodoItem(['title' => 'Incomplete']);
        $completed = new TodoItem(['title' => 'Incomplete', 'is_completed' => true]);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($todo);
        $this->repository->shouldReceive('update')->with($todo, ['is_completed' => true])->once()->andReturn($completed);

        $result = $this->service->markComplete(1);

        $this->assertSame($completed, $result);
    }

    public function test_markComplete_throws_ModelNotFoundException_when_not_found(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->markComplete(99);
    }
}
