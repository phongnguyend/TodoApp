<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

class UserResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id' => $this->id,
            'username' => $this->username,
            'email' => $this->email,
            'is_active' => $this->is_active,
            'created_at' => $this->created_at?->toIso8601String(),
            'created_by_user_id' => $this->created_by_user_id,
            'updated_at' => $this->updated_at?->toIso8601String(),
            'updated_by_user_id' => $this->updated_by_user_id,
        ];
    }
}
