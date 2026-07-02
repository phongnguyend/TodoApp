using Microsoft.AspNetCore.Http;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoShared.Models;

namespace TodoApi.Services;

public class FileService(IFileRepository repository, AppDbContext db, IConfiguration configuration) : IFileService
{
    private readonly string _storageDir = configuration["FileStorage:Path"] ?? "uploads";
    private readonly long _maxUploadSizeBytes = configuration.GetValue<long?>("FileStorage:MaxUploadSizeBytes") ?? 10 * 1024 * 1024;

    // ── Mapping ───────────────────────────────────────────────────────────────

    private static FileResponse ToResponse(FileEntity file) =>
        new(file.Id, file.Name, file.Extension, file.Size, file.ContentType, file.CreatedAt, file.UpdatedAt);

    private static PaginatedResponse<FileResponse> ToPaginated(
        IEnumerable<FileEntity> items, int total, int page, int pageSize) =>
        new(
            items.Select(ToResponse),
            total,
            page,
            pageSize,
            (int)Math.Ceiling(total / (double)pageSize)
        );

    private async Task<FileEntity> GetOrThrowAsync(int id, CancellationToken ct)
    {
        var file = await repository.GetByIdAsync(id, ct)
            ?? throw new KeyNotFoundException($"File {id} not found.");
        return file;
    }

    // ── Queries ───────────────────────────────────────────────────────────────

    public async Task<PaginatedResponse<FileResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var (items, total) = await repository.GetAllAsync((page - 1) * pageSize, pageSize, ct);
        return ToPaginated(items, total, page, pageSize);
    }

    public async Task<FileResponse> GetByIdAsync(int id, CancellationToken ct = default)
    {
        var file = await GetOrThrowAsync(id, ct);
        return ToResponse(file);
    }

    public async Task<FileDownloadTarget> GetDownloadTargetAsync(int id, CancellationToken ct = default)
    {
        var file = await GetOrThrowAsync(id, ct);
        if (!System.IO.File.Exists(file.Location))
        {
            throw new KeyNotFoundException($"File {id} content not found on disk.");
        }
        return new FileDownloadTarget(file.Location, file.Name, file.ContentType ?? "application/octet-stream");
    }

    // ── Commands ──────────────────────────────────────────────────────────────

    public async Task<FileResponse> UploadAsync(IFormFile file, CancellationToken ct = default)
    {
        if (file.Length > _maxUploadSizeBytes)
        {
            throw new FileTooLargeException($"File exceeds the maximum allowed size of {_maxUploadSizeBytes} bytes.");
        }

        // Strip any directory components from the client-supplied name to prevent path traversal.
        var originalName = Path.GetFileName(file.FileName);
        if (string.IsNullOrEmpty(originalName))
        {
            originalName = "unnamed";
        }
        var extension = Path.GetExtension(originalName).TrimStart('.').ToLowerInvariant();

        Directory.CreateDirectory(_storageDir);

        // A random prefix avoids collisions/overwrites between uploads that share a name.
        var storedName = $"{Guid.NewGuid():N}_{originalName}";
        var location = Path.GetFullPath(Path.Combine(_storageDir, storedName));

        await using (var stream = System.IO.File.Create(location))
        {
            await file.CopyToAsync(stream, ct);
        }

        var entity = new FileEntity
        {
            Name = originalName,
            Extension = extension,
            Size = file.Length,
            ContentType = file.ContentType,
            Location = location,
            CreatedAt = DateTime.UtcNow,
        };
        var created = await repository.AddAsync(entity, ct);
        await db.SaveChangesAsync(ct);
        return ToResponse(created);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var file = await GetOrThrowAsync(id, ct);
        repository.Delete(file);
        await db.SaveChangesAsync(ct);
        try
        {
            if (System.IO.File.Exists(file.Location))
            {
                System.IO.File.Delete(file.Location);
            }
        }
        catch (IOException)
        {
            // Content already missing or inaccessible on disk - nothing left to clean up.
        }
    }
}
