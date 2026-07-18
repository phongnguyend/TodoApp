<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        Schema::create('todo_item_attachments', function (Blueprint $table) {
            $table->id();
            $table->foreignId('todo_item_id')->constrained('todo_items')->cascadeOnDelete();
            $table->foreignId('file_id')->constrained('files')->cascadeOnDelete();
            $table->timestamps();
            $table->unique(['todo_item_id', 'file_id']);
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('todo_item_attachments');
    }
};
