<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

return new class extends Migration
{
    public function up(): void
    {
        foreach (['todo_items', 'todo_item_attachments', 'email_logs', 'files', 'users'] as $tableName) {
            Schema::table($tableName, function (Blueprint $table): void {
                $table->foreignId('created_by_user_id')->nullable()->constrained('users')->nullOnDelete();
                $table->foreignId('updated_by_user_id')->nullable()->constrained('users')->nullOnDelete();
            });
        }
    }

    public function down(): void
    {
        foreach (array_reverse(['todo_items', 'todo_item_attachments', 'email_logs', 'files', 'users']) as $tableName) {
            Schema::table($tableName, function (Blueprint $table): void {
                $table->dropConstrainedForeignId('updated_by_user_id');
                $table->dropConstrainedForeignId('created_by_user_id');
            });
        }
    }
};
