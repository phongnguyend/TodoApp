<?php

namespace Tests\Unit\Requests;

use App\Http\Requests\UploadFileRequest;
use Illuminate\Http\UploadedFile;
use Illuminate\Support\Facades\Validator;
use Tests\TestCase;

class UploadFileRequestTest extends TestCase
{
    private function validate(array $data): \Illuminate\Validation\Validator
    {
        return Validator::make($data, (new UploadFileRequest())->rules());
    }

    public function test_valid_file_passes(): void
    {
        $v = $this->validate(['file' => UploadedFile::fake()->create('photo.png', 10, 'image/png')]);
        $this->assertFalse($v->fails());
    }

    public function test_file_is_required(): void
    {
        $v = $this->validate([]);
        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('file', $v->errors()->toArray());
    }

    public function test_rejects_files_exceeding_default_max_size(): void
    {
        // Default MAX_UPLOAD_SIZE_BYTES is 10 MB (10240 KB); one KB over that must fail.
        $v = $this->validate(['file' => UploadedFile::fake()->create('big.bin', 10241)]);

        $this->assertTrue($v->fails());
        $this->assertArrayHasKey('file', $v->errors()->toArray());
    }

    public function test_accepts_file_at_default_max_size(): void
    {
        // Exactly at the default 10 MB (10240 KB) limit must pass ("max" is inclusive).
        $v = $this->validate(['file' => UploadedFile::fake()->create('at-limit.bin', 10240)]);

        $this->assertFalse($v->fails());
    }
}
