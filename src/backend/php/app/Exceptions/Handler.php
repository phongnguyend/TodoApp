<?php

namespace App\Exceptions;

use Illuminate\Database\Eloquent\ModelNotFoundException;
use Illuminate\Foundation\Exceptions\Handler as ExceptionHandler;
use Illuminate\Validation\ValidationException;

/**
 * Global exception handler - analogous to ASP.NET Core's UseExceptionHandler /
 * IExceptionFilter / ProblemDetails middleware.
 */
class Handler extends ExceptionHandler
{
    protected $dontFlash = [
        'current_password',
        'password',
        'password_confirmation',
    ];

    public function register(): void
    {
        $this->renderable(function (UserConflictException $e) {
            return response()->json(['message' => $e->getMessage()], 409);
        });

        $this->renderable(function (InvalidPasswordException|InvalidPasswordResetTokenException $e) {
            return response()->json(['message' => $e->getMessage()], 400);
        });

        $this->renderable(function (ModelNotFoundException $e) {
            return response()->json([
                'message' => $e->getMessage() ?: 'Resource not found.',
            ], 404);
        });

        $this->renderable(function (ValidationException $e) {
            return response()->json([
                'message' => 'Validation failed.',
                'errors' => $e->errors(),
            ], 422);
        });
    }
}
