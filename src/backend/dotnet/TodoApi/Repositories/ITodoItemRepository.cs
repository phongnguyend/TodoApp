using TodoShared.Models;

namespace TodoApi.Repositories;

public interface ITodoItemRepository : IRepository<TodoItem>
{
    Task<TodoItem?> GetByTitleAsync(string title, CancellationToken ct = default);
    Task<(IEnumerable<TodoItem> Items, int Total)> GetIncompleteAsync(int skip, int take, CancellationToken ct = default);
}
