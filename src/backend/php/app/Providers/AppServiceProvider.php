<?php

namespace App\Providers;

use App\Repositories\Contracts\TodoItemRepositoryInterface;
use App\Repositories\TodoItemRepository;
use App\Services\Contracts\TodoItemServiceInterface;
use App\Services\TodoItemService;
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

        // Bind service interface → concrete implementation
        $this->app->bind(TodoItemServiceInterface::class, TodoItemService::class);
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
