<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\HasMany;

/**
 * Eloquent model for a todo item.
 * Analogous to an EF Core entity class with data annotations.
 *
 * @property int         $id
 * @property string      $title
 * @property string|null $description
 * @property bool        $is_completed
 * @property \Carbon\Carbon $created_at
 * @property \Carbon\Carbon|null $updated_at
 */
class TodoItem extends Model
{
    use HasFactory;

    protected $table = 'todo_items';

    /**
     * Mass-assignable attributes (analogous to [BindProperty] or DTO mapping).
     */
    protected $fillable = [
        'title',
        'description',
        'is_completed',
        'created_by_user_id',
        'updated_by_user_id',
    ];

    /**
     * Default attribute values (mirrors the DB column defaults).
     */
    protected $attributes = [
        'is_completed' => false,
    ];

    /**
     * Attribute type casts (analogous to EF value converters / column types).
     */
    protected $casts = [
        'is_completed' => 'boolean',
        'created_at'   => 'datetime',
        'updated_at'   => 'datetime',
    ];

    public function attachments(): HasMany
    {
        return $this->hasMany(TodoItemAttachment::class);
    }
}
