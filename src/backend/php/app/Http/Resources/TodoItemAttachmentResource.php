<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class TodoItemAttachmentResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'todo_item_id' => $this->todo_item_id,
            'file_id' => $this->file_id,
            'created_at' => $this->created_at?->toIso8601String(),
            'created_by_user_id' => $this->created_by_user_id,
            'updated_at' => $this->updated_at?->toIso8601String(),
            'updated_by_user_id' => $this->updated_by_user_id,
        ];
    }
}
