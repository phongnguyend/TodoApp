<?php

namespace App\Providers;

use App\Repositories\Contracts\FileRepositoryInterface;
use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Repositories\Contracts\TodoItemAttachmentRepositoryInterface;
use App\Repositories\FileRepository;
use App\Repositories\TodoItemRepository;
use App\Repositories\TodoItemAttachmentRepository;
use App\Services\Contracts\FileServiceInterface;
use App\Services\Contracts\TodoItemServiceInterface;
use App\Services\Contracts\TodoItemAttachmentServiceInterface;
use App\Services\FileService;
use App\Services\TodoItemService;
use App\Services\TodoItemAttachmentService;
use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\Response;
use Illuminate\Support\ServiceProvider;
use Symfony\Component\HttpFoundation\Response as HttpResponse;

class AppServiceProvider extends ServiceProvider
{
    /**
     * Register application services into the IoC container.
     * Analogous to builder.Services.AddScoped<IService, Service>() in Program.cs.
     */
    public function register(): void
    {
        // Bind repository interface → concrete implementation (scoped per request by default in Laravel)
        $this->app->bind(TodoItemRepositoryInterface::class, TodoItemRepository::class);
        $this->app->bind(FileRepositoryInterface::class, FileRepository::class);
        $this->app->bind(TodoItemAttachmentRepositoryInterface::class, TodoItemAttachmentRepository::class);

        // Bind service interface → concrete implementation
        $this->app->bind(TodoItemServiceInterface::class, TodoItemService::class);
        $this->app->bind(FileServiceInterface::class, FileService::class);
        $this->app->bind(TodoItemAttachmentServiceInterface::class, TodoItemAttachmentService::class);
    }

    /**
     * Bootstrap application services.
     * Analogous to app.UseExceptionHandler() / problem-details middleware in ASP.NET Core.
     */
    public function boot(): void
    {
        // Return consistent JSON 404 responses for ModelNotFoundException (like ProblemDetails in ASP.NET Core)
        $this->app['router']->bind('id', fn ($value) => $value);

        $this->app->singleton(
            \Illuminate\Contracts\Debug\ExceptionHandler::class,
            \App\Exceptions\Handler::class
        );
    }
}
