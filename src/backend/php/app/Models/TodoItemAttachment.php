<?php

namespace App\Models;

use Illuminate\Database\Eloquent\Factories\HasFactory;
use Illuminate\Database\Eloquent\Model;
use Illuminate\Database\Eloquent\Relations\BelongsTo;

class TodoItemAttachment extends Model
{
    use HasFactory;

    protected $fillable = ['todo_item_id', 'file_id'];

    public function todoItem(): BelongsTo
    {
        return $this->belongsTo(TodoItem::class);
    }

    public function file(): BelongsTo
    {
        return $this->belongsTo(File::class);
    }
}
