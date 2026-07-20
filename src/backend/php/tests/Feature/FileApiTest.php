<?php

namespace Tests\Feature;

use App\Models\File;
use Illuminate\Foundation\Testing\RefreshDatabase;
use Illuminate\Http\UploadedFile;
use Tests\TestCase;

class FileApiTest extends TestCase
{
    use RefreshDatabase;

    private string $tempDir;

    protected function setUp(): void
    {
        parent::setUp();
        $this->authenticateRequests();

        $this->tempDir = sys_get_temp_dir() . '/todo-php-file-api-test-' . uniqid('', true);
        mkdir($this->tempDir, 0755, true);
        putenv("FILE_STORAGE_PATH={$this->tempDir}");
        $_ENV['FILE_STORAGE_PATH'] = $this->tempDir;
    }

    protected function tearDown(): void
    {
        putenv('FILE_STORAGE_PATH');
        unset($_ENV['FILE_STORAGE_PATH']);
        foreach (glob($this->tempDir . '/*') ?: [] as $leftover) {
            @unlink($leftover);
        }
        @rmdir($this->tempDir);

        parent::tearDown();
    }

    // ── GET /api/files ────────────────────────────────────────────────────────

    public function test_index_returns_paginated_list(): void
    {
        File::factory()->count(3)->create();

        $response = $this->getJson('/api/files');

        $response->assertOk()
                 ->assertJsonStructure(['data', 'meta', 'links'])
                 ->assertJsonCount(3, 'data');
    }

    public function test_index_returns_empty_list_when_no_files(): void
    {
        $response = $this->getJson('/api/files');

        $response->assertOk()
                 ->assertJsonCount(0, 'data');
    }

    public function test_index_respects_page_size_parameter(): void
    {
        File::factory()->count(5)->create();

        $response = $this->getJson('/api/files?page_size=2');

        $response->assertOk()
                 ->assertJsonCount(2, 'data');
    }

    // ── GET /api/files/{id} ───────────────────────────────────────────────────

    public function test_show_returns_file_when_found(): void
    {
        $file = File::factory()->create(['name' => 'report.pdf']);

        $response = $this->getJson("/api/files/{$file->id}");

        $response->assertOk()
                 ->assertJsonPath('data.id', $file->id)
                 ->assertJsonPath('data.name', 'report.pdf');
    }

    public function test_show_returns_404_when_not_found(): void
    {
        $response = $this->getJson('/api/files/9999');

        $response->assertNotFound();
    }

    public function test_show_does_not_expose_location(): void
    {
        $file = File::factory()->create();

        $response = $this->getJson("/api/files/{$file->id}");

        $response->assertJsonMissingPath('data.location');
    }

    // ── POST /api/files ───────────────────────────────────────────────────────

    public function test_store_creates_and_returns_file(): void
    {
        $upload = UploadedFile::fake()->create('photo.png', 11, 'image/png');

        $response = $this->post('/api/files', ['file' => $upload]);

        $response->assertCreated()
                 ->assertJsonPath('data.name', 'photo.png')
                 ->assertJsonPath('data.extension', 'png')
                 ->assertJsonPath('data.content_type', 'image/png');

        $this->assertDatabaseHas('files', ['name' => 'photo.png']);
        $this->assertCount(1, glob($this->tempDir . '/*') ?: []);
    }

    public function test_store_returns_422_when_no_file_provided(): void
    {
        $response = $this->postJson('/api/files', []);

        $response->assertUnprocessable()
                 ->assertJsonValidationErrors(['file']);
    }

    // ── GET /api/files/{id}/download ──────────────────────────────────────────

    public function test_download_returns_file_content(): void
    {
        $upload = UploadedFile::fake()->create('notes.txt', 1, 'text/plain');
        $uploadResponse = $this->post('/api/files', ['file' => $upload]);
        $fileId = $uploadResponse->json('data.id');

        $response = $this->get("/api/files/{$fileId}/download");

        $response->assertOk();
        $this->assertStringStartsWith('text/plain', $response->headers->get('content-type'));
    }

    public function test_download_returns_404_when_not_found(): void
    {
        $response = $this->getJson('/api/files/9999/download');

        $response->assertNotFound();
    }

    // ── DELETE /api/files/{id} ────────────────────────────────────────────────

    public function test_destroy_deletes_file_and_returns_204(): void
    {
        $upload = UploadedFile::fake()->create('to-delete.txt', 1);
        $uploadResponse = $this->post('/api/files', ['file' => $upload]);
        $fileId = $uploadResponse->json('data.id');

        $response = $this->deleteJson("/api/files/{$fileId}");

        $response->assertNoContent();
        $this->assertDatabaseMissing('files', ['id' => $fileId]);
        $this->assertCount(0, glob($this->tempDir . '/*') ?: []);
    }

    public function test_destroy_returns_404_when_not_found(): void
    {
        $response = $this->deleteJson('/api/files/9999');

        $response->assertNotFound();
    }
}
