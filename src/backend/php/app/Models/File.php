<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

/**
 * Eloquent model for an uploaded file.
 * Analogous to an EF Core entity class with data annotations.
 *
 * @property int         $id
 * @property string      $name
 * @property string      $extension
 * @property int         $size
 * @property string|null $content_type
 * @property string      $location
 * @property \Carbon\Carbon $created_at
 * @property \Carbon\Carbon|null $updated_at
 */
class File extends Model
{
    use HasFactory;

    protected $table = 'files';

    /**
     * Mass-assignable attributes (analogous to [BindProperty] or DTO mapping).
     */
    protected $fillable = [
        'name',
        'extension',
        'size',
        'content_type',
        'location',
    ];

    /**
     * Attribute type casts (analogous to EF value converters / column types).
     */
    protected $casts = [
        'size'       => 'integer',
        'created_at' => 'datetime',
        'updated_at' => 'datetime',
    ];

    public function todoItemAttachments(): HasMany
    {
        return $this->hasMany(TodoItemAttachment::class);
    }
}
