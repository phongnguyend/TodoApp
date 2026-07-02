namespace TodoApi.DTOs;

public record FileResponse(
    int Id,
    string Name,
    string Extension,
    long Size,
    string? ContentType,
    DateTime CreatedAt,
    DateTime? UpdatedAt
);

public record FileDownloadTarget(
    string Path,
    string Name,
    string ContentType
);
