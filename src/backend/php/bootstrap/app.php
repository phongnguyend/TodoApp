<?php

use App\Console\Commands\ProcessIncompleteRemindersCommand;
use Illuminate\Console\Scheduling\Schedule;
use Illuminate\Foundation\Application;
use Illuminate\Foundation\Configuration\Exceptions;
use Illuminate\Foundation\Configuration\Middleware;

$app = Application::configure(basePath: dirname(__DIR__))
    ->withRouting(
        api: __DIR__.'/../routes/api.php',
        health: '/up',
    )
    ->withCommands([
        ProcessIncompleteRemindersCommand::class,
    ])
    ->withSchedule(function (Schedule $schedule): void {
        // Run every hour; adjust the frequency to suit your needs.
        $schedule->command('app:process-incomplete-reminders')->hourly();
    })
    ->withMiddleware(function (Middleware $middleware) {
        //
    })
    ->withExceptions(function (Exceptions $exceptions) {
        //
    })->create();

$app->singleton(
    \Illuminate\Contracts\Debug\ExceptionHandler::class,
    \App\Exceptions\Handler::class,
);

return $app;
