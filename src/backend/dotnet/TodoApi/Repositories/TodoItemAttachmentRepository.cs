using Microsoft.EntityFrameworkCore;
using TodoApi.Data;
using TodoShared.Models;

namespace TodoApi.Repositories;

public class TodoItemAttachmentRepository(AppDbContext db) : BaseRepository<TodoItemAttachment>(db), ITodoItemAttachmentRepository
{
    public async Task<IReadOnlyList<TodoItemAttachment>> GetByTodoItemIdAsync(int todoItemId, CancellationToken ct = default)
        => await Db.TodoItemAttachments
            .Where(a => a.TodoItemId == todoItemId)
            .OrderByDescending(a => a.CreatedAt)
            .ThenByDescending(a => a.Id)
            .ToListAsync(ct);

    public async Task<TodoItemAttachment?> GetByIdForTodoItemAsync(int todoItemId, int attachmentId, CancellationToken ct = default)
        => await Db.TodoItemAttachments
            .FirstOrDefaultAsync(a => a.Id == attachmentId && a.TodoItemId == todoItemId, ct);

    public async Task<TodoItemAttachment?> GetByTodoItemAndFileAsync(int todoItemId, int fileId, CancellationToken ct = default)
        => await Db.TodoItemAttachments
            .FirstOrDefaultAsync(a => a.TodoItemId == todoItemId && a.FileId == fileId, ct);
}
