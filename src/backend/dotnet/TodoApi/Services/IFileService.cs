using Microsoft.AspNetCore.Http;
using TodoApi.DTOs;

namespace TodoApi.Services;

public interface IFileService
{
    Task<PaginatedResponse<FileResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default);
    Task<FileResponse> GetByIdAsync(int id, CancellationToken ct = default);
    Task<FileDownloadTarget> GetDownloadTargetAsync(int id, CancellationToken ct = default);
    Task<FileResponse> UploadAsync(IFormFile file, CancellationToken ct = default);
    Task DeleteAsync(int id, CancellationToken ct = default);
}
