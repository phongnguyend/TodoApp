<?php

namespace Tests\Feature;

use App\Models\File;
use App\Models\TodoItem;
use App\Models\TodoItemAttachment;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class TodoItemAttachmentApiTest extends TestCase
{
    use RefreshDatabase;

    public function test_attachment_crud_is_scoped_to_todo_item(): void
    {
        $todo = TodoItem::factory()->create();
        $otherTodo = TodoItem::factory()->create();
        $file = File::factory()->create();
        $replacement = File::factory()->create();

        $created = $this->postJson("/api/todo-items/{$todo->id}/attachments", ['file_id' => $file->id]);
        $created->assertCreated()
            ->assertJsonPath('data.todo_item_id', $todo->id)
            ->assertJsonPath('data.file_id', $file->id);
        $attachmentId = $created->json('data.id');

        $this->getJson("/api/todo-items/{$todo->id}/attachments")
            ->assertOk()->assertJsonCount(1, 'data');
        $this->getJson("/api/todo-items/{$todo->id}/attachments/{$attachmentId}")
            ->assertOk()->assertJsonPath('data.id', $attachmentId);
        $this->getJson("/api/todo-items/{$otherTodo->id}/attachments/{$attachmentId}")
            ->assertNotFound();

        $this->putJson("/api/todo-items/{$todo->id}/attachments/{$attachmentId}", ['file_id' => $replacement->id])
            ->assertOk()->assertJsonPath('data.file_id', $replacement->id);
        $this->deleteJson("/api/todo-items/{$todo->id}/attachments/{$attachmentId}")
            ->assertNoContent();
        $this->assertDatabaseMissing('todo_item_attachments', ['id' => $attachmentId]);
    }

    public function test_create_is_idempotent_and_validates_references(): void
    {
        $todo = TodoItem::factory()->create();
        $file = File::factory()->create();

        $firstId = $this->postJson("/api/todo-items/{$todo->id}/attachments", ['file_id' => $file->id])
            ->assertCreated()->json('data.id');
        $this->postJson("/api/todo-items/{$todo->id}/attachments", ['file_id' => $file->id])
            ->assertCreated()->assertJsonPath('data.id', $firstId);

        $this->assertDatabaseCount('todo_item_attachments', 1);
        $this->postJson("/api/todo-items/{$todo->id}/attachments", [])->assertUnprocessable();
        $this->postJson("/api/todo-items/{$todo->id}/attachments", ['file_id' => 9999])->assertNotFound();
        $this->getJson('/api/todo-items/9999/attachments')->assertNotFound();
    }

    public function test_deleting_parent_records_cascades_attachment_references(): void
    {
        $todo = TodoItem::factory()->create();
        $file = File::factory()->create();
        TodoItemAttachment::query()->create(['todo_item_id' => $todo->id, 'file_id' => $file->id]);

        $todo->delete();

        $this->assertDatabaseCount('todo_item_attachments', 0);
    }
}
