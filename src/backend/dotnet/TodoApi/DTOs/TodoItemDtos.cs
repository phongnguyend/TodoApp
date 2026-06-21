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
    DateTime? UpdatedAt
);

public record PaginatedResponse<T>(
    IEnumerable<T> Items,
    int Total,
    int Page,
    int PageSize,
    int TotalPages
);
