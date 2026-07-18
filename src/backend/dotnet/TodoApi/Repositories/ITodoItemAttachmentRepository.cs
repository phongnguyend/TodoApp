using TodoShared.Models;

namespace TodoApi.Repositories;

public interface ITodoItemAttachmentRepository : IRepository<TodoItemAttachment>
{
    Task<IReadOnlyList<TodoItemAttachment>> GetByTodoItemIdAsync(int todoItemId, CancellationToken ct = default);
    Task<TodoItemAttachment?> GetByIdForTodoItemAsync(int todoItemId, int attachmentId, CancellationToken ct = default);
    Task<TodoItemAttachment?> GetByTodoItemAndFileAsync(int todoItemId, int fileId, CancellationToken ct = default);
}
