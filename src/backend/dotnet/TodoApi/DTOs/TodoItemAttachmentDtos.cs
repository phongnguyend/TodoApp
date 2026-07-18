namespace TodoApi.DTOs;

public record CreateTodoItemAttachmentRequest(int FileId);

public record TodoItemAttachmentResponse(
    int Id,
    int TodoItemId,
    int FileId,
    DateTime CreatedAt,
    DateTime? UpdatedAt
);
