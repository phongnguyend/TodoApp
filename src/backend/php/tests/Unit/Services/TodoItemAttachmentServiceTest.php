<?php

namespace Tests\Unit\Services;

use App\Models\File;
use App\Models\TodoItem;
use App\Models\TodoItemAttachment;
use App\Repositories\Contracts\FileRepositoryInterface;
use App\Repositories\Contracts\TodoItemAttachmentRepositoryInterface;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Services\TodoItemAttachmentService;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Support\Collection;
use Mockery;
use Tests\TestCase;

class TodoItemAttachmentServiceTest extends TestCase
{
    private TodoItemAttachmentRepositoryInterface $attachments;
    private TodoItemRepositoryInterface $todos;
    private FileRepositoryInterface $files;
    private TodoItemAttachmentService $service;

    protected function setUp(): void
    {
        parent::setUp();

        $this->attachments = Mockery::mock(TodoItemAttachmentRepositoryInterface::class);
        $this->todos = Mockery::mock(TodoItemRepositoryInterface::class);
        $this->files = Mockery::mock(FileRepositoryInterface::class);
        $this->service = new TodoItemAttachmentService($this->attachments, $this->todos, $this->files);
    }

    protected function tearDown(): void
    {
        Mockery::close();
        parent::tearDown();
    }

    public function test_getAll_returns_attachments_for_existing_todo(): void
    {
        $items = new Collection([$this->attachment(3)]);
        $this->todos->shouldReceive('findById')->with(10)->once()->andReturn(new TodoItem());
        $this->attachments->shouldReceive('getForTodo')->with(10)->once()->andReturn($items);

        $this->assertSame($items, $this->service->getAll(10));
    }

    public function test_getAll_throws_when_todo_does_not_exist(): void
    {
        $this->todos->shouldReceive('findById')->with(10)->once()->andReturn(null);
        $this->attachments->shouldNotReceive('getForTodo');

        $this->expectException(ModelNotFoundException::class);
        $this->service->getAll(10);
    }

    public function test_getById_rejects_attachment_from_another_todo(): void
    {
        $this->todos->shouldReceive('findById')->with(10)->once()->andReturn(new TodoItem());
        $this->attachments->shouldReceive('findForTodo')->with(10, 3)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->getById(10, 3);
    }

    public function test_create_adds_attachment_when_references_exist(): void
    {
        $created = $this->attachment(1, 10, 5);
        $this->expectExistingTodoAndFile(10, 5);
        $this->attachments->shouldReceive('findForTodoAndFile')->with(10, 5)->once()->andReturn(null);
        $this->attachments->shouldReceive('create')
            ->with(['todo_item_id' => 10, 'file_id' => 5])->once()->andReturn($created);

        $this->assertSame($created, $this->service->create(10, 5));
    }

    public function test_create_returns_existing_duplicate_without_creating(): void
    {
        $existing = $this->attachment(7, 10, 5);
        $this->expectExistingTodoAndFile(10, 5);
        $this->attachments->shouldReceive('findForTodoAndFile')->with(10, 5)->once()->andReturn($existing);
        $this->attachments->shouldNotReceive('create');

        $this->assertSame($existing, $this->service->create(10, 5));
    }

    public function test_create_throws_when_file_does_not_exist(): void
    {
        $this->todos->shouldReceive('findById')->with(10)->once()->andReturn(new TodoItem());
        $this->files->shouldReceive('findById')->with(99)->once()->andReturn(null);
        $this->attachments->shouldNotReceive('create');

        $this->expectException(ModelNotFoundException::class);
        $this->service->create(10, 99);
    }

    public function test_update_changes_the_attachment_file(): void
    {
        $current = $this->attachment(3, 10, 5);
        $updated = $this->attachment(3, 10, 6);
        $this->expectExistingTodoAndFile(10, 6);
        $this->attachments->shouldReceive('findForTodo')->with(10, 3)->once()->andReturn($current);
        $this->attachments->shouldReceive('findForTodoAndFile')->with(10, 6)->once()->andReturn(null);
        $this->attachments->shouldReceive('update')->with($current, ['file_id' => 6])->once()->andReturn($updated);

        $this->assertSame($updated, $this->service->update(10, 3, 6));
    }

    public function test_update_returns_existing_duplicate_without_updating(): void
    {
        $current = $this->attachment(3, 10, 5);
        $existing = $this->attachment(4, 10, 6);
        $this->expectExistingTodoAndFile(10, 6);
        $this->attachments->shouldReceive('findForTodo')->with(10, 3)->once()->andReturn($current);
        $this->attachments->shouldReceive('findForTodoAndFile')->with(10, 6)->once()->andReturn($existing);
        $this->attachments->shouldNotReceive('update');

        $this->assertSame($existing, $this->service->update(10, 3, 6));
    }

    public function test_delete_removes_scoped_attachment(): void
    {
        $attachment = $this->attachment(3, 10, 5);
        $this->todos->shouldReceive('findById')->with(10)->once()->andReturn(new TodoItem());
        $this->attachments->shouldReceive('findForTodo')->with(10, 3)->once()->andReturn($attachment);
        $this->attachments->shouldReceive('delete')->with($attachment)->once();

        $this->service->delete(10, 3);
        $this->addToAssertionCount(1);
    }

    private function expectExistingTodoAndFile(int $todoId, int $fileId): void
    {
        $this->todos->shouldReceive('findById')->with($todoId)->once()->andReturn(new TodoItem());
        $this->files->shouldReceive('findById')->with($fileId)->once()->andReturn(new File());
    }

    private function attachment(int $id, int $todoId = 10, int $fileId = 5): TodoItemAttachment
    {
        $attachment = new TodoItemAttachment(['todo_item_id' => $todoId, 'file_id' => $fileId]);
        $attachment->id = $id;

        return $attachment;
    }
}
