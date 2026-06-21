using Microsoft.EntityFrameworkCore;
using TodoApi.Data;
using TodoApi.Models;

namespace TodoApi.Repositories;

public class TodoItemRepository(AppDbContext db) : BaseRepository<TodoItem>(db), ITodoItemRepository
{
    public async Task<TodoItem?> GetByTitleAsync(string title, CancellationToken ct = default)
        => await Db.TodoItems.FirstOrDefaultAsync(t => t.Title == title, ct);

    public async Task<(IEnumerable<TodoItem> Items, int Total)> GetIncompleteAsync(int skip, int take, CancellationToken ct = default)
    {
        var query = Db.TodoItems.Where(t => !t.IsCompleted);
        var total = await query.CountAsync(ct);
        var items = await query.Skip(skip).Take(take).ToListAsync(ct);
        return (items, total);
    }
}
