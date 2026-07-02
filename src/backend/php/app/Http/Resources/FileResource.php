<?php

namespace App\Http\Resources;

use Illuminate\Http\Request;
use Illuminate\Http\Resources\Json\JsonResource;

/**
 * Transforms a File model into an API response.
 * Analogous to a DTO / AutoMapper profile in ASP.NET Core.
 *
 * Note: the on-disk `location` is intentionally not exposed to clients; content is
 * retrieved via the dedicated download endpoint instead.
 */
class FileResource extends JsonResource
{
    public function toArray(Request $request): array
    {
        return [
            'id'           => $this->id,
            'name'         => $this->name,
            'extension'    => $this->extension,
            'size'         => $this->size,
            'content_type' => $this->content_type,
            'created_at'   => $this->created_at?->toIso8601String(),
            'updated_at'   => $this->updated_at?->toIso8601String(),
        ];
    }
}
