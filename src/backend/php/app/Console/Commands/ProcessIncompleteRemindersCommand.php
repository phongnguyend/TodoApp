<?php

namespace App\Console\Commands;

use App\Models\EmailLog;
use App\Models\TodoItem;
use Illuminate\Console\Command;
use Illuminate\Support\Facades\Mail;

/**
 * Artisan command executed by the background worker container.
 *
 * Steps:
 *  1. Fetch all incomplete todo items.
 *  2. Build a plain-text email body.
 *  3. Persist an EmailLog record with status = 'pending'.
 *  4. Attempt SMTP delivery.
 *  5. Update the log to 'sent' (with sent_at) or 'failed' (with error_message).
 */
class ProcessIncompleteRemindersCommand extends Command
{
    protected $signature = 'app:process-incomplete-reminders';

    protected $description = 'Email a reminder listing all incomplete todo items and record the attempt in email_logs';

    public function handle(): int
    {
        $items = TodoItem::where('is_completed', false)
            ->orderByDesc('created_at')
            ->get();

        if ($items->isEmpty()) {
            $this->info('No incomplete todo items — skipping reminder email.');
            return Command::SUCCESS;
        }

        $recipient = (string) env('MAIL_REMINDER_RECIPIENT', 'admin@example.com');
        $subject   = sprintf('Todo Reminder: %d incomplete item(s)', $items->count());
        $body      = $this->buildBody($items);

        // Persist audit record before attempting delivery
        $log = EmailLog::create([
            'recipient' => $recipient,
            'subject'   => $subject,
            'body'      => $body,
            'status'    => 'pending',
        ]);

        try {
            Mail::raw($body, static function ($message) use ($recipient, $subject): void {
                $message->to($recipient)->subject($subject);
            });

            $log->update([
                'status'  => 'sent',
                'sent_at' => now(),
            ]);

            $this->info("Reminder sent to {$recipient} ({$items->count()} item(s)).");
        } catch (\Throwable $e) {
            $log->update([
                'status'        => 'failed',
                'error_message' => $e->getMessage(),
            ]);

            $this->error("Failed to send reminder: {$e->getMessage()}");
            return Command::FAILURE;
        }

        return Command::SUCCESS;
    }

    /**
     * Build a plain-text email body listing every incomplete todo item.
     *
     * @param \Illuminate\Database\Eloquent\Collection<int, TodoItem> $items
     */
    private function buildBody(iterable $items): string
    {
        $lines = [
            'The following todo items are still incomplete:',
            str_repeat('-', 44),
        ];

        foreach ($items as $item) {
            $lines[] = sprintf('• [#%d] %s', $item->id, $item->title);
            if ($item->description) {
                $lines[] = '  ' . $item->description;
            }
        }

        $lines[] = '';
        $lines[] = 'Generated at: ' . now()->toDateTimeString() . ' UTC';

        return implode("\n", $lines);
    }
}
