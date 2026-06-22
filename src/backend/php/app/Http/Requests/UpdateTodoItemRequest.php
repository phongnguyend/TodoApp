<?php

namespace App\Http\Requests;

use Illuminate\Foundation\Http\FormRequest;

/**
 * Validates incoming update requests (all fields optional — PATCH semantics).
 * Analogous to an UpdateTodoItemRequest DTO with nullable properties in C#.
 */
class UpdateTodoItemRequest extends FormRequest
{
    public function authorize(): bool
    {
        return true;
    }

    public function rules(): array
    {
        return [
            'title'        => ['sometimes', 'filled', 'string', 'max:200'],
            'description'  => ['sometimes', 'nullable', 'string', 'max:2000'],
            'is_completed' => ['sometimes', 'boolean'],
        ];
    }
}
