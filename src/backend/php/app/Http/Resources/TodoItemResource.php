<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * Transforms a TodoItem model into an API response.
 * Analogous to a DTO / AutoMapper profile in ASP.NET Core.
 */
class TodoItemResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id'           => $this->id,
            'title'        => $this->title,
            'description'  => $this->description,
            'is_completed' => $this->is_completed,
            'created_at'   => $this->created_at?->toIso8601String(),
            'created_by_user_id' => $this->created_by_user_id,
            'updated_at'   => $this->updated_at?->toIso8601String(),
            'updated_by_user_id' => $this->updated_by_user_id,
        ];
    }
}
