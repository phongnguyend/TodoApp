<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Validates incoming file-upload requests.
 * Analogous to a [FromForm] IFormFile parameter with size/required validation in C#.
 */
class UploadFileRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        $maxKilobytes = intdiv((int) env('MAX_UPLOAD_SIZE_BYTES', 10 * 1024 * 1024), 1024);

        return [
            'file' => ['required', 'file', "max:{$maxKilobytes}"],
        ];
    }
}
