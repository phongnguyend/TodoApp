namespace TodoApi.DTOs;

public record CreateTodoItemAttachmentRequest(int FileId);

public record TodoItemAttachmentResponse(
    int Id,
    int TodoItemId,
    int FileId,
    DateTime CreatedAt,
    int? CreatedByUserId,
    DateTime? UpdatedAt,
    int? UpdatedByUserId
)
{
    public TodoItemAttachmentResponse(int id, int todoItemId, int fileId,
        DateTime createdAt, DateTime? updatedAt)
        : this(id, todoItemId, fileId, createdAt, null, updatedAt, null) { }
}
