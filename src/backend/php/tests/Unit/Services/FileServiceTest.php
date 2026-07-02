<?php

namespace Tests\Unit\Services;

use App\Models\File;
use App\Repositories\Contracts\FileRepositoryInterface;
use App\Services\FileService;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Http\UploadedFile;
use Illuminate\Pagination\LengthAwarePaginator;
use Mockery;
use Tests\TestCase;

class FileServiceTest extends TestCase
{
    private FileRepositoryInterface $repository;
    private FileService $service;
    private string $tempDir;

    protected function setUp(): void
    {
        parent::setUp();

        $this->tempDir = sys_get_temp_dir() . '/todo-php-file-service-test-' . uniqid('', true);
        mkdir($this->tempDir, 0755, true);
        putenv("FILE_STORAGE_PATH={$this->tempDir}");
        $_ENV['FILE_STORAGE_PATH'] = $this->tempDir;

        $this->repository = Mockery::mock(FileRepositoryInterface::class);
        $this->service    = new FileService($this->repository);
    }

    protected function tearDown(): void
    {
        putenv('FILE_STORAGE_PATH');
        unset($_ENV['FILE_STORAGE_PATH']);
        foreach (glob($this->tempDir . '/*') ?: [] as $leftover) {
            @unlink($leftover);
        }
        @rmdir($this->tempDir);

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

    // ── getById ───────────────────────────────────────────────────────────────

    public function test_getById_returns_file_when_found(): void
    {
        $file = new File(['name' => 'test.txt']);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($file);

        $result = $this->service->getById(1);

        $this->assertSame($file, $result);
    }

    public function test_getById_throws_ModelNotFoundException_when_not_found(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->getById(99);
    }

    // ── upload ────────────────────────────────────────────────────────────────

    public function test_upload_saves_file_to_storage_dir_and_returns_metadata(): void
    {
        $this->repository->shouldReceive('create')->once()
            ->with(Mockery::on(function (array $data) {
                return $data['name'] === 'photo.png'
                    && $data['extension'] === 'png'
                    && $data['content_type'] === 'image/png'
                    && is_file($data['location']);
            }))
            ->andReturnUsing(fn (array $data) => new File($data));

        $upload = UploadedFile::fake()->create('photo.png', 11, 'image/png');

        $result = $this->service->upload($upload);

        $this->assertSame('photo.png', $result->name);
        $this->assertSame('png', $result->extension);
        $this->assertCount(1, glob($this->tempDir . '/*') ?: []);
    }

    public function test_upload_sanitizes_path_traversal_in_filename(): void
    {
        $this->repository->shouldReceive('create')->once()
            ->andReturnUsing(fn (array $data) => new File($data));

        $upload = UploadedFile::fake()->create('../../etc/passwd', 1);

        $result = $this->service->upload($upload);

        $this->assertSame('passwd', $result->name);
        $savedFiles = glob($this->tempDir . '/*') ?: [];
        $this->assertCount(1, $savedFiles);
        $this->assertStringNotContainsString('..', $savedFiles[0]);
    }

    // ── getDownloadTarget ─────────────────────────────────────────────────────

    public function test_getDownloadTarget_returns_path_name_and_content_type(): void
    {
        $location = $this->tempDir . '/stored_download.txt';
        file_put_contents($location, 'file contents');

        $file = new File([
            'name'         => 'download.txt',
            'content_type' => 'text/plain',
            'location'     => $location,
        ]);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($file);

        $target = $this->service->getDownloadTarget(1);

        $this->assertSame($location, $target['path']);
        $this->assertSame('download.txt', $target['name']);
        $this->assertSame('text/plain', $target['content_type']);
    }

    public function test_getDownloadTarget_throws_when_content_missing_from_disk(): void
    {
        $file = new File([
            'name'     => 'missing.txt',
            'location' => $this->tempDir . '/does-not-exist.txt',
        ]);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($file);

        $this->expectException(ModelNotFoundException::class);
        $this->service->getDownloadTarget(1);
    }

    public function test_getDownloadTarget_throws_when_not_found(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->getDownloadTarget(99);
    }

    // ── delete ────────────────────────────────────────────────────────────────

    public function test_delete_removes_record_and_file_from_disk(): void
    {
        $location = $this->tempDir . '/to-delete.txt';
        file_put_contents($location, 'bye');

        $file = new File(['name' => 'to-delete.txt', 'location' => $location]);
        $this->repository->shouldReceive('findById')->with(1)->once()->andReturn($file);
        $this->repository->shouldReceive('delete')->with($file)->once();

        $this->service->delete(1);

        $this->assertFileDoesNotExist($location);
    }

    public function test_delete_throws_when_not_found(): void
    {
        $this->repository->shouldReceive('findById')->with(99)->once()->andReturn(null);

        $this->expectException(ModelNotFoundException::class);
        $this->service->delete(99);
    }
}
