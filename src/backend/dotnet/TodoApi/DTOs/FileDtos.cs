namespace TodoApi.DTOs;

public record FileResponse(
    int Id,
    string Name,
    string Extension,
    long Size,
    string? ContentType,
    DateTime CreatedAt,
    int? CreatedByUserId,
    DateTime? UpdatedAt,
    int? UpdatedByUserId
)
{
    public FileResponse(int id, string name, string extension, long size, string? contentType,
        DateTime createdAt, DateTime? updatedAt)
        : this(id, name, extension, size, contentType, createdAt, null, updatedAt, null) { }
}

public record FileDownloadTarget(
    string Path,
    string Name,
    string ContentType
);
