using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoShared.Models;

namespace TodoApi.Services;

public class TodoItemAttachmentService(
    ITodoItemAttachmentRepository attachmentRepository,
    ITodoItemRepository todoItemRepository,
    IFileRepository fileRepository,
    AppDbContext db,
    IHttpContextAccessor? httpContextAccessor = null) : ITodoItemAttachmentService
{
    private int? ActorUserId => AuditActor.GetUserId(httpContextAccessor);

    private static TodoItemAttachmentResponse ToResponse(TodoItemAttachment attachment) =>
        new(attachment.Id, attachment.TodoItemId, attachment.FileId, attachment.CreatedAt,
            attachment.CreatedByUserId, attachment.UpdatedAt, attachment.UpdatedByUserId);

    private async Task<TodoItem> GetTodoItemOrThrowAsync(int todoItemId, CancellationToken ct)
    {
        var todoItem = await todoItemRepository.GetByIdAsync(todoItemId, ct)
            ?? throw new KeyNotFoundException($"Todo item {todoItemId} not found.");
        return todoItem;
    }

    private async Task<FileEntity> GetFileOrThrowAsync(int fileId, CancellationToken ct)
    {
        var file = await fileRepository.GetByIdAsync(fileId, ct)
            ?? throw new KeyNotFoundException($"File {fileId} not found.");
        return file;
    }

    private async Task<TodoItemAttachment> GetOrThrowAsync(int todoItemId, int attachmentId, CancellationToken ct)
    {
        var attachment = await attachmentRepository.GetByIdForTodoItemAsync(todoItemId, attachmentId, ct)
            ?? throw new KeyNotFoundException($"Attachment {attachmentId} not found for todo item {todoItemId}.");
        return attachment;
    }

    public async Task<IReadOnlyList<TodoItemAttachmentResponse>> GetAllAsync(int todoItemId, CancellationToken ct = default)
    {
        await GetTodoItemOrThrowAsync(todoItemId, ct);
        var attachments = await attachmentRepository.GetByTodoItemIdAsync(todoItemId, ct);
        return attachments.Select(ToResponse).ToList();
    }

    public async Task<TodoItemAttachmentResponse> GetByIdAsync(int todoItemId, int attachmentId, CancellationToken ct = default)
    {
        await GetTodoItemOrThrowAsync(todoItemId, ct);
        var attachment = await GetOrThrowAsync(todoItemId, attachmentId, ct);
        return ToResponse(attachment);
    }

    public async Task<TodoItemAttachmentResponse> CreateAsync(int todoItemId, CreateTodoItemAttachmentRequest request, CancellationToken ct = default)
    {
        await GetTodoItemOrThrowAsync(todoItemId, ct);
        await GetFileOrThrowAsync(request.FileId, ct);

        var existing = await attachmentRepository.GetByTodoItemAndFileAsync(todoItemId, request.FileId, ct);
        if (existing is not null)
        {
            return ToResponse(existing);
        }

        var attachment = new TodoItemAttachment
        {
            TodoItemId = todoItemId,
            FileId = request.FileId,
            CreatedAt = DateTime.UtcNow,
            CreatedByUserId = ActorUserId,
        };

        await attachmentRepository.AddAsync(attachment, ct);
        await db.SaveChangesAsync(ct);
        return ToResponse(attachment);
    }

    public async Task<TodoItemAttachmentResponse> UpdateAsync(int todoItemId, int attachmentId, CreateTodoItemAttachmentRequest request, CancellationToken ct = default)
    {
        await GetTodoItemOrThrowAsync(todoItemId, ct);
        await GetFileOrThrowAsync(request.FileId, ct);
        var attachment = await GetOrThrowAsync(todoItemId, attachmentId, ct);

        var existing = await attachmentRepository.GetByTodoItemAndFileAsync(todoItemId, request.FileId, ct);
        if (existing is not null && existing.Id != attachment.Id)
        {
            return ToResponse(existing);
        }

        attachment.FileId = request.FileId;
        attachment.UpdatedAt = DateTime.UtcNow;
        attachment.UpdatedByUserId = ActorUserId;
        attachmentRepository.Update(attachment);
        await db.SaveChangesAsync(ct);
        return ToResponse(attachment);
    }

    public async Task DeleteAsync(int todoItemId, int attachmentId, CancellationToken ct = default)
    {
        await GetTodoItemOrThrowAsync(todoItemId, ct);
        var attachment = await GetOrThrowAsync(todoItemId, attachmentId, ct);
        attachmentRepository.Delete(attachment);
        await db.SaveChangesAsync(ct);
    }
}
