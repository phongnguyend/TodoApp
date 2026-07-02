<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * Create files table.
 * Stores metadata about uploaded files; content is stored on disk and referenced via `location`.
 */
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('files', function (Blueprint $table) {
            $table->id();                                          // bigint PK, auto-increment
            $table->string('name', 255);
            $table->string('extension', 20);
            $table->unsignedBigInteger('size');
            $table->string('content_type', 100)->nullable();
            $table->string('location', 500);
            $table->timestamps();                                  // created_at + updated_at
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('files');
    }
};
