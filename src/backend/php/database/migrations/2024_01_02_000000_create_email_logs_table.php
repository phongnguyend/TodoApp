<?php

use Illuminate\Database\Migrations\Migration;
use Illuminate\Database\Schema\Blueprint;
use Illuminate\Support\Facades\Schema;

/**
 * Create email_logs table.
 * Audit trail for every outbound email attempt - records persist even when SMTP delivery fails.
 */
return new class extends Migration
{
    public function up(): void
    {
        Schema::create('email_logs', function (Blueprint $table) {
            $table->id();
            $table->string('recipient', 255);
            $table->string('subject', 500);
            $table->text('body');
            $table->string('status', 50)->default('pending'); // pending | sent | failed
            $table->timestamp('created_at')->useCurrent();
            $table->timestamp('sent_at')->nullable();
            $table->text('error_message')->nullable();
        });
    }

    public function down(): void
    {
        Schema::dropIfExists('email_logs');
    }
};
