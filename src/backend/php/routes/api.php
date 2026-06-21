<?php

use App\Http\Controllers\Api\TodoItemController;
use Illuminate\Support\Facades\Route;

/*
|--------------------------------------------------------------------------
| API Routes
|--------------------------------------------------------------------------
| All routes here are automatically prefixed with /api by bootstrap/app.php.
| This mirrors ASP.NET Core's app.MapControllers() with [Route("api/[controller]")].
|
*/

Route::prefix('todo-items')->group(function () {
    Route::get('/',            [TodoItemController::class, 'index']);       // GET    /api/todo-items
    Route::get('/incomplete',  [TodoItemController::class, 'incomplete']);  // GET    /api/todo-items/incomplete
    Route::get('/{id}',        [TodoItemController::class, 'show']);        // GET    /api/todo-items/{id}
    Route::post('/',           [TodoItemController::class, 'store']);       // POST   /api/todo-items
    Route::put('/{id}',        [TodoItemController::class, 'update']);      // PUT    /api/todo-items/{id}
    Route::patch('/{id}/complete', [TodoItemController::class, 'complete']); // PATCH /api/todo-items/{id}/complete
    Route::delete('/{id}',    [TodoItemController::class, 'destroy']);      // DELETE /api/todo-items/{id}
});
