using System.ComponentModel.DataAnnotations;

namespace TodoApi.DTOs;

public record CreateTodoItemRequest(
    [Required, StringLength(200, MinimumLength = 1)] string Title,
    [StringLength(2000)] string? Description
);

public record UpdateTodoItemRequest(
    [StringLength(200, MinimumLength = 1)] string? Title,
    [StringLength(2000)] string? Description,
    bool? IsCompleted
);

public record TodoItemResponse(
    int Id,
    string Title,
    string? Description,
    bool IsCompleted,
    DateTime CreatedAt,
    int? CreatedByUserId,
    DateTime? UpdatedAt,
    int? UpdatedByUserId
)
{
    public TodoItemResponse(int id, string title, string? description, bool isCompleted,
        DateTime createdAt, DateTime? updatedAt)
        : this(id, title, description, isCompleted, createdAt, null, updatedAt, null) { }
}

public record PaginatedResponse<T>(
    IEnumerable<T> Items,
    int Total,
    int Page,
    int PageSize,
    int TotalPages
);

public record ImportResult(
    int Imported,
    int Failed,
    IReadOnlyList<ImportRowError> Errors
);

public record ImportRowError(
    int Row,
    string Error
);
