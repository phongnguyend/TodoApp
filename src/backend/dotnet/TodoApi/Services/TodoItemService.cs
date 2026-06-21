using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Models;
using TodoApi.Repositories;

namespace TodoApi.Services;

public class TodoItemService(ITodoItemRepository repository, AppDbContext db) : ITodoItemService
{
    // ── Mapping ───────────────────────────────────────────────────────────────

    private static TodoItemResponse ToResponse(TodoItem item) =>
        new(item.Id, item.Title, item.Description, item.IsCompleted, item.CreatedAt, item.UpdatedAt);

    private static PaginatedResponse<TodoItemResponse> ToPaginated(
        IEnumerable<TodoItem> items, int total, int page, int pageSize) =>
        new(
            items.Select(ToResponse),
            total,
            page,
            pageSize,
            (int)Math.Ceiling(total / (double)pageSize)
        );

    private async Task<TodoItem> GetOrThrowAsync(int id, CancellationToken ct)
    {
        var item = await repository.GetByIdAsync(id, ct)
            ?? throw new KeyNotFoundException($"Todo item {id} not found.");
        return item;
    }

    // ── Queries ───────────────────────────────────────────────────────────────

    public async Task<PaginatedResponse<TodoItemResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var (items, total) = await repository.GetAllAsync((page - 1) * pageSize, pageSize, ct);
        return ToPaginated(items, total, page, pageSize);
    }

    public async Task<PaginatedResponse<TodoItemResponse>> GetIncompleteAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var (items, total) = await repository.GetIncompleteAsync((page - 1) * pageSize, pageSize, ct);
        return ToPaginated(items, total, page, pageSize);
    }

    public async Task<TodoItemResponse> GetByIdAsync(int id, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        return ToResponse(item);
    }

    // ── Commands ──────────────────────────────────────────────────────────────

    public async Task<TodoItemResponse> CreateAsync(CreateTodoItemRequest request, CancellationToken ct = default)
    {
        var item = new TodoItem
        {
            Title = request.Title,
            Description = request.Description,
            CreatedAt = DateTime.UtcNow,
        };
        var created = await repository.AddAsync(item, ct);
        await db.SaveChangesAsync(ct);
        return ToResponse(created);
    }

    public async Task<TodoItemResponse> UpdateAsync(int id, UpdateTodoItemRequest request, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        if (request.Title is not null) item.Title = request.Title;
        if (request.Description is not null) item.Description = request.Description;
        if (request.IsCompleted is not null) item.IsCompleted = request.IsCompleted.Value;
        item.UpdatedAt = DateTime.UtcNow;
        repository.Update(item);
        await db.SaveChangesAsync(ct);
        return ToResponse(item);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        repository.Delete(item);
        await db.SaveChangesAsync(ct);
    }

    public async Task<TodoItemResponse> MarkCompleteAsync(int id, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        item.IsCompleted = true;
        item.UpdatedAt = DateTime.UtcNow;
        repository.Update(item);
        await db.SaveChangesAsync(ct);
        return ToResponse(item);
    }
}
