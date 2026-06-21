using TodoApi.DTOs;

namespace TodoApi.Services;

public interface ITodoItemService
{
    Task<PaginatedResponse<TodoItemResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default);
    Task<PaginatedResponse<TodoItemResponse>> GetIncompleteAsync(int page, int pageSize, CancellationToken ct = default);
    Task<TodoItemResponse> GetByIdAsync(int id, CancellationToken ct = default);
    Task<TodoItemResponse> CreateAsync(CreateTodoItemRequest request, CancellationToken ct = default);
    Task<TodoItemResponse> UpdateAsync(int id, UpdateTodoItemRequest request, CancellationToken ct = default);
    Task DeleteAsync(int id, CancellationToken ct = default);
    Task<TodoItemResponse> MarkCompleteAsync(int id, CancellationToken ct = default);
}
