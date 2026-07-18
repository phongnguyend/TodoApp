using TodoApi.DTOs;

namespace TodoApi.Services;

public interface ITodoItemAttachmentService
{
    Task<IReadOnlyList<TodoItemAttachmentResponse>> GetAllAsync(int todoItemId, CancellationToken ct = default);
    Task<TodoItemAttachmentResponse> GetByIdAsync(int todoItemId, int attachmentId, CancellationToken ct = default);
    Task<TodoItemAttachmentResponse> CreateAsync(int todoItemId, CreateTodoItemAttachmentRequest request, CancellationToken ct = default);
    Task<TodoItemAttachmentResponse> UpdateAsync(int todoItemId, int attachmentId, CreateTodoItemAttachmentRequest request, CancellationToken ct = default);
    Task DeleteAsync(int todoItemId, int attachmentId, CancellationToken ct = default);
}
