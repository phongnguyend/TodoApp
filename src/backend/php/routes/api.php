<?php

use App\Http\Controllers\Api\FileController;
use App\Http\Controllers\Api\TodoItemAttachmentController;
use App\Http\Controllers\Api\TodoItemController;
use App\Http\Controllers\Api\UserController;
use App\Http\Middleware\AuthenticateUser;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
| All routes here are automatically prefixed with /api by bootstrap/app.php.
| This mirrors ASP.NET Core's app.MapControllers() with [Route("api/[controller]")].
|
*/

Route::middleware(AuthenticateUser::class)->prefix('todo-items')->group(function () {
    Route::get('/', [TodoItemController::class, 'index']);       // GET    /api/todo-items
    Route::get('/incomplete', [TodoItemController::class, 'incomplete']);  // GET    /api/todo-items/incomplete
    Route::post('/import/csv', [TodoItemController::class, 'importCsv']);   // POST   /api/todo-items/import/csv
    Route::post('/import/excel', [TodoItemController::class, 'importExcel']); // POST   /api/todo-items/import/excel
    Route::get('/export/csv', [TodoItemController::class, 'exportCsv']);   // GET    /api/todo-items/export/csv
    Route::get('/export/excel', [TodoItemController::class, 'exportExcel']); // GET    /api/todo-items/export/excel
    Route::get('/{id}/attachments', [TodoItemAttachmentController::class, 'index']);
    Route::post('/{id}/attachments', [TodoItemAttachmentController::class, 'store']);
    Route::get('/{id}/attachments/{attachmentId}', [TodoItemAttachmentController::class, 'show']);
    Route::put('/{id}/attachments/{attachmentId}', [TodoItemAttachmentController::class, 'update']);
    Route::delete('/{id}/attachments/{attachmentId}', [TodoItemAttachmentController::class, 'destroy']);
    Route::get('/{id}', [TodoItemController::class, 'show']);        // GET    /api/todo-items/{id}
    Route::post('/', [TodoItemController::class, 'store']);       // POST   /api/todo-items
    Route::put('/{id}', [TodoItemController::class, 'update']);      // PUT    /api/todo-items/{id}
    Route::patch('/{id}/complete', [TodoItemController::class, 'complete']); // PATCH /api/todo-items/{id}/complete
    Route::delete('/{id}', [TodoItemController::class, 'destroy']);      // DELETE /api/todo-items/{id}
});

Route::middleware(AuthenticateUser::class)->prefix('files')->group(function () {
    Route::get('/', [FileController::class, 'index']);    // GET    /api/files
    Route::get('/{id}/download', [FileController::class, 'download']); // GET    /api/files/{id}/download
    Route::get('/{id}', [FileController::class, 'show']);     // GET    /api/files/{id}
    Route::post('/', [FileController::class, 'store']);    // POST   /api/files
    Route::delete('/{id}', [FileController::class, 'destroy']);  // DELETE /api/files/{id}
});

Route::prefix('users')->group(function () {
    Route::post('/signup', [UserController::class, 'signup']);
    Route::post('/password/reset', [UserController::class, 'requestPasswordReset']);
    Route::post('/password/confirm', [UserController::class, 'confirmPasswordReset']);
});

Route::middleware(AuthenticateUser::class)->prefix('users')->group(function () {
    // Static routes must precede /{id}.
    Route::post('/password/change', [UserController::class, 'changePassword']);
    Route::get('/profile', [UserController::class, 'profile']);
    Route::put('/profile', [UserController::class, 'updateProfile']);
    Route::get('/', [UserController::class, 'index']);
    Route::post('/', [UserController::class, 'store']);
    Route::get('/{id}', [UserController::class, 'show'])->whereNumber('id');
    Route::put('/{id}', [UserController::class, 'update'])->whereNumber('id');
    Route::patch('/{id}/activate', [UserController::class, 'activate'])->whereNumber('id');
    Route::patch('/{id}/deactivate', [UserController::class, 'deactivate'])->whereNumber('id');
});

Route::post('tokens', [UserController::class, 'createToken']);
