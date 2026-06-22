<?php

namespace Tests\Feature;

use App\Models\TodoItem;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Tests\TestCase;

class TodoItemApiTest extends TestCase
{
    use RefreshDatabase;

    // ── GET /api/todo-items ───────────────────────────────────────────────────

    public function test_index_returns_paginated_list(): void
    {
        TodoItem::factory()->count(3)->create();

        $response = $this->getJson('/api/todo-items');

        $response->assertOk()
                 ->assertJsonStructure(['data', 'meta', 'links'])
                 ->assertJsonCount(3, 'data');
    }

    public function test_index_returns_empty_list_when_no_items(): void
    {
        $response = $this->getJson('/api/todo-items');

        $response->assertOk()
                 ->assertJsonCount(0, 'data');
    }

    public function test_index_respects_page_size_parameter(): void
    {
        TodoItem::factory()->count(5)->create();

        $response = $this->getJson('/api/todo-items?page_size=2');

        $response->assertOk()
                 ->assertJsonCount(2, 'data');
    }

    // ── GET /api/todo-items/incomplete ────────────────────────────────────────

    public function test_incomplete_returns_only_incomplete_items(): void
    {
        TodoItem::factory()->count(2)->create(['is_completed' => false]);
        TodoItem::factory()->count(3)->completed()->create();

        $response = $this->getJson('/api/todo-items/incomplete');

        $response->assertOk()
                 ->assertJsonCount(2, 'data');

        foreach ($response->json('data') as $item) {
            $this->assertFalse($item['is_completed']);
        }
    }

    // ── GET /api/todo-items/{id} ──────────────────────────────────────────────

    public function test_show_returns_todo_when_found(): void
    {
        $todo = TodoItem::factory()->create(['title' => 'Read a book']);

        $response = $this->getJson("/api/todo-items/{$todo->id}");

        $response->assertOk()
                 ->assertJsonPath('data.id', $todo->id)
                 ->assertJsonPath('data.title', 'Read a book');
    }

    public function test_show_returns_404_when_not_found(): void
    {
        $response = $this->getJson('/api/todo-items/9999');

        $response->assertNotFound();
    }

    // ── POST /api/todo-items ──────────────────────────────────────────────────

    public function test_store_creates_and_returns_todo(): void
    {
        $payload = ['title' => 'Buy groceries', 'description' => 'Milk and eggs'];

        $response = $this->postJson('/api/todo-items', $payload);

        $response->assertCreated()
                 ->assertJsonPath('data.title', 'Buy groceries')
                 ->assertJsonPath('data.description', 'Milk and eggs')
                 ->assertJsonPath('data.is_completed', false);

        $this->assertDatabaseHas('todo_items', ['title' => 'Buy groceries']);
    }

    public function test_store_returns_422_when_title_missing(): void
    {
        $response = $this->postJson('/api/todo-items', ['description' => 'No title']);

        $response->assertUnprocessable()
                 ->assertJsonValidationErrors(['title']);
    }

    public function test_store_returns_422_when_title_too_long(): void
    {
        $response = $this->postJson('/api/todo-items', ['title' => str_repeat('a', 201)]);

        $response->assertUnprocessable()
                 ->assertJsonValidationErrors(['title']);
    }

    // ── PUT /api/todo-items/{id} ──────────────────────────────────────────────

    public function test_update_modifies_existing_todo(): void
    {
        $todo = TodoItem::factory()->create(['title' => 'Old title']);

        $response = $this->putJson("/api/todo-items/{$todo->id}", ['title' => 'New title']);

        $response->assertOk()
                 ->assertJsonPath('data.title', 'New title');

        $this->assertDatabaseHas('todo_items', ['id' => $todo->id, 'title' => 'New title']);
    }

    public function test_update_returns_404_when_not_found(): void
    {
        $response = $this->putJson('/api/todo-items/9999', ['title' => 'Ghost']);

        $response->assertNotFound();
    }

    // ── PATCH /api/todo-items/{id}/complete ───────────────────────────────────

    public function test_complete_marks_todo_as_completed(): void
    {
        $todo = TodoItem::factory()->create(['is_completed' => false]);

        $response = $this->patchJson("/api/todo-items/{$todo->id}/complete");

        $response->assertOk()
                 ->assertJsonPath('data.is_completed', true);

        $this->assertDatabaseHas('todo_items', ['id' => $todo->id, 'is_completed' => true]);
    }

    public function test_complete_returns_404_when_not_found(): void
    {
        $response = $this->patchJson('/api/todo-items/9999/complete');

        $response->assertNotFound();
    }

    // ── DELETE /api/todo-items/{id} ───────────────────────────────────────────

    public function test_destroy_deletes_todo_and_returns_204(): void
    {
        $todo = TodoItem::factory()->create();

        $response = $this->deleteJson("/api/todo-items/{$todo->id}");

        $response->assertNoContent();
        $this->assertDatabaseMissing('todo_items', ['id' => $todo->id]);
    }

    public function test_destroy_returns_404_when_not_found(): void
    {
        $response = $this->deleteJson('/api/todo-items/9999');

        $response->assertNotFound();
    }
}
