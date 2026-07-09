<?php

namespace Tests\Feature;

use App\Models\TodoItem;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Http\UploadedFile;
use PhpOffice\PhpSpreadsheet\IOFactory;
use PhpOffice\PhpSpreadsheet\Spreadsheet;
use PhpOffice\PhpSpreadsheet\Writer\Xlsx;
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

    // ── POST /api/todo-items/import/csv ───────────────────────────────────────

    public function test_importCsv_creates_todos_and_reports_counts(): void
    {
        $csv = "title,description,is_completed\n"
             . "Buy groceries,Milk and eggs,false\n"
             . "Read a book,,true\n"
             . ",Missing title,false\n";
        $file = UploadedFile::fake()->createWithContent('todos.csv', $csv);

        $response = $this->postJson('/api/todo-items/import/csv', ['file' => $file]);

        $response->assertOk()
                 ->assertJson([
                     'imported' => 2,
                     'failed'   => 1,
                 ])
                 ->assertJsonPath('errors.0.row', 4)
                 ->assertJsonPath('errors.0.error', 'Title is required.');

        $this->assertDatabaseHas('todo_items', ['title' => 'Buy groceries', 'description' => 'Milk and eggs']);
        $this->assertDatabaseHas('todo_items', ['title' => 'Read a book', 'is_completed' => true]);
        $this->assertDatabaseCount('todo_items', 2);
    }

    public function test_importCsv_returns_422_when_file_missing(): void
    {
        $response = $this->postJson('/api/todo-items/import/csv', []);

        $response->assertUnprocessable()
                 ->assertJsonValidationErrors(['file']);
    }

    // ── GET /api/todo-items/export/csv ────────────────────────────────────────

    public function test_exportCsv_returns_csv_content_for_all_todos(): void
    {
        TodoItem::factory()->create(['title' => 'Buy groceries', 'description' => 'Milk and eggs']);
        TodoItem::factory()->completed()->create(['title' => 'Read a book']);

        $response = $this->get('/api/todo-items/export/csv');

        $response->assertOk()
                 ->assertHeader('Content-Type', 'text/csv; charset=UTF-8')
                 ->assertHeader('Content-Disposition', 'attachment; filename="todo_items.csv"');

        $content = $response->getContent();
        $this->assertStringContainsString('id,title,description,is_completed,created_at,updated_at', $content);
        $this->assertStringContainsString('Buy groceries', $content);
        $this->assertStringContainsString('Read a book', $content);
    }

    public function test_exportCsv_returns_header_only_when_no_todos(): void
    {
        $response = $this->get('/api/todo-items/export/csv');

        $response->assertOk();
        $lines = array_filter(explode("\n", (string) $response->getContent()));
        $this->assertCount(1, $lines);
    }

    // ── POST /api/todo-items/import/excel ─────────────────────────────

    public function test_importExcel_creates_todos_and_reports_counts(): void
    {
        $file = $this->makeExcelUploadFile([
            ['title', 'description', 'is_completed'],
            ['Buy groceries', 'Milk and eggs', false],
            ['Read a book', '', true],
            ['', 'Missing title', false],
        ]);

        $response = $this->postJson('/api/todo-items/import/excel', ['file' => $file]);

        $response->assertOk()
                 ->assertJson([
                     'imported' => 2,
                     'failed'   => 1,
                 ])
                 ->assertJsonPath('errors.0.row', 4)
                 ->assertJsonPath('errors.0.error', 'Title is required.');

        $this->assertDatabaseHas('todo_items', ['title' => 'Buy groceries', 'description' => 'Milk and eggs']);
        $this->assertDatabaseHas('todo_items', ['title' => 'Read a book', 'is_completed' => true]);
        $this->assertDatabaseCount('todo_items', 2);
    }

    public function test_importExcel_returns_422_when_file_missing(): void
    {
        $response = $this->postJson('/api/todo-items/import/excel', []);

        $response->assertUnprocessable()
                 ->assertJsonValidationErrors(['file']);
    }

    // ── GET /api/todo-items/export/excel ────────────────────────────

    public function test_exportExcel_returns_excel_content_for_all_todos(): void
    {
        TodoItem::factory()->create(['title' => 'Buy groceries', 'description' => 'Milk and eggs']);
        TodoItem::factory()->completed()->create(['title' => 'Read a book']);

        $response = $this->get('/api/todo-items/export/excel');

        $response->assertOk()
                 ->assertHeader('Content-Type', 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet')
                 ->assertHeader('Content-Disposition', 'attachment; filename="todo_items.xlsx"');

        $rows = $this->readExcelRows($response->getContent());
        $this->assertSame(['id', 'title', 'description', 'is_completed', 'created_at', 'updated_at'], $rows[0]);
        $titles = array_column($rows, 1);
        $this->assertContains('Buy groceries', $titles);
        $this->assertContains('Read a book', $titles);
    }

    public function test_exportExcel_returns_header_only_when_no_todos(): void
    {
        $response = $this->get('/api/todo-items/export/excel');

        $response->assertOk();
        $rows = $this->readExcelRows($response->getContent());
        $this->assertCount(1, $rows);
    }

    // ── Excel test helpers ────────────────────────────────────────────────

    private function makeExcelUploadFile(array $rows): UploadedFile
    {
        $spreadsheet = new Spreadsheet();
        $sheet = $spreadsheet->getActiveSheet();
        foreach ($rows as $rowIndex => $row) {
            $sheet->fromArray($row, null, 'A' . ($rowIndex + 1));
        }

        $stream = fopen('php://temp', 'r+');
        (new Xlsx($spreadsheet))->save($stream);
        rewind($stream);
        $content = (string) stream_get_contents($stream);
        fclose($stream);

        return UploadedFile::fake()->createWithContent('todos.xlsx', $content);
    }

    private function readExcelRows(string $content): array
    {
        $path = tempnam(sys_get_temp_dir(), 'xlsx');
        file_put_contents($path, $content);

        try {
            return IOFactory::load($path)->getActiveSheet()->toArray(null, true, false, false);
        } finally {
            unlink($path);
        }
    }
}
