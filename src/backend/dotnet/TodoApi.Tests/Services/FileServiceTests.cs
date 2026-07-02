using Microsoft.AspNetCore.Http;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Moq;
using TodoApi.Data;
using TodoApi.Repositories;
using TodoApi.Services;
using TodoShared.Models;

namespace TodoApi.Tests.Services;

public class FileServiceTests : IDisposable
{
    private readonly Mock<IFileRepository> _repoMock;
    private readonly AppDbContext _db;
    private readonly string _storageDir;
    private readonly FileService _sut;

    public FileServiceTests()
    {
        _repoMock = new Mock<IFileRepository>();

        var options = new DbContextOptionsBuilder<AppDbContext>()
            .UseInMemoryDatabase(databaseName: Guid.NewGuid().ToString())
            .Options;
        _db = new AppDbContext(options);

        _storageDir = Path.Combine(Path.GetTempPath(), "todo-api-tests-" + Guid.NewGuid());

        var configuration = new ConfigurationBuilder()
            .AddInMemoryCollection(new Dictionary<string, string?>
            {
                ["FileStorage:Path"] = _storageDir,
                ["FileStorage:MaxUploadSizeBytes"] = "1024",
            })
            .Build();

        _sut = new FileService(_repoMock.Object, _db, configuration);
    }

    public void Dispose()
    {
        _db.Dispose();
        if (Directory.Exists(_storageDir))
        {
            Directory.Delete(_storageDir, recursive: true);
        }
    }

    private static Mock<IFormFile> CreateFormFile(string fileName, string content, string? contentType = "text/plain")
    {
        var bytes = System.Text.Encoding.UTF8.GetBytes(content);
        var formFile = new Mock<IFormFile>();
        formFile.Setup(f => f.FileName).Returns(fileName);
        formFile.Setup(f => f.Length).Returns(bytes.Length);
        formFile.Setup(f => f.ContentType).Returns(contentType!);
        formFile.Setup(f => f.CopyToAsync(It.IsAny<Stream>(), It.IsAny<CancellationToken>()))
                .Returns((Stream stream, CancellationToken ct) => stream.WriteAsync(bytes, 0, bytes.Length, ct));
        return formFile;
    }

    // ── GetAllAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task GetAllAsync_ReturnsPaginatedResponse()
    {
        var items = new List<FileEntity>
        {
            new() { Id = 1, Name = "a.txt", Extension = "txt", Size = 10, Location = "a", CreatedAt = DateTime.UtcNow },
            new() { Id = 2, Name = "b.txt", Extension = "txt", Size = 20, Location = "b", CreatedAt = DateTime.UtcNow },
        };
        _repoMock.Setup(r => r.GetAllAsync(0, 20, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((items, 2));

        var result = await _sut.GetAllAsync(1, 20);

        Assert.Equal(2, result.Total);
        Assert.Equal(2, result.Items.Count());
    }

    // ── GetByIdAsync ──────────────────────────────────────────────────────────

    [Fact]
    public async Task GetByIdAsync_ReturnsFile_WhenFound()
    {
        var file = new FileEntity { Id = 1, Name = "a.txt", Extension = "txt", Size = 10, Location = "a", CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(file);

        var result = await _sut.GetByIdAsync(1);

        Assert.Equal(1, result.Id);
        Assert.Equal("a.txt", result.Name);
    }

    [Fact]
    public async Task GetByIdAsync_ThrowsKeyNotFoundException_WhenNotFound()
    {
        _repoMock.Setup(r => r.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((FileEntity?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.GetByIdAsync(99));
    }

    // ── UploadAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task UploadAsync_WritesFileToDiskAndSavesMetadata()
    {
        var formFile = CreateFormFile("report.pdf", "hello world");
        _repoMock.Setup(r => r.AddAsync(It.IsAny<FileEntity>(), It.IsAny<CancellationToken>()))
                 .ReturnsAsync((FileEntity f, CancellationToken _) => f);

        var result = await _sut.UploadAsync(formFile.Object);

        Assert.Equal("report.pdf", result.Name);
        Assert.Equal("pdf", result.Extension);
        Assert.Equal("text/plain", result.ContentType);
        _repoMock.Verify(r => r.AddAsync(It.Is<FileEntity>(f => File.Exists(f.Location)), It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task UploadAsync_StripsDirectoryComponentsFromFileName()
    {
        var formFile = CreateFormFile("../../etc/passwd", "malicious");
        _repoMock.Setup(r => r.AddAsync(It.IsAny<FileEntity>(), It.IsAny<CancellationToken>()))
                 .ReturnsAsync((FileEntity f, CancellationToken _) => f);

        var result = await _sut.UploadAsync(formFile.Object);

        Assert.Equal("passwd", result.Name);
        _repoMock.Verify(r => r.AddAsync(
            It.Is<FileEntity>(f => f.Location.StartsWith(_storageDir, StringComparison.Ordinal)),
            It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task UploadAsync_ThrowsFileTooLargeException_WhenExceedsMaxSize()
    {
        var formFile = CreateFormFile("big.bin", new string('x', 2000));

        await Assert.ThrowsAsync<FileTooLargeException>(() => _sut.UploadAsync(formFile.Object));
        _repoMock.Verify(r => r.AddAsync(It.IsAny<FileEntity>(), It.IsAny<CancellationToken>()), Times.Never);
    }

    // ── GetDownloadTargetAsync ────────────────────────────────────────────────

    [Fact]
    public async Task GetDownloadTargetAsync_ThrowsKeyNotFoundException_WhenMetadataMissing()
    {
        _repoMock.Setup(r => r.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((FileEntity?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.GetDownloadTargetAsync(99));
    }

    [Fact]
    public async Task GetDownloadTargetAsync_ThrowsKeyNotFoundException_WhenContentMissingOnDisk()
    {
        var file = new FileEntity { Id = 1, Name = "a.txt", Extension = "txt", Size = 10, Location = Path.Combine(_storageDir, "missing.txt"), CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(file);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.GetDownloadTargetAsync(1));
    }

    [Fact]
    public async Task GetDownloadTargetAsync_ReturnsTarget_WhenFileExists()
    {
        Directory.CreateDirectory(_storageDir);
        var path = Path.Combine(_storageDir, "a.txt");
        await File.WriteAllTextAsync(path, "hello");
        var file = new FileEntity { Id = 1, Name = "a.txt", Extension = "txt", Size = 5, ContentType = "text/plain", Location = path, CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(file);

        var result = await _sut.GetDownloadTargetAsync(1);

        Assert.Equal(path, result.Path);
        Assert.Equal("a.txt", result.Name);
        Assert.Equal("text/plain", result.ContentType);
    }

    [Fact]
    public async Task GetDownloadTargetAsync_DefaultsContentType_WhenNull()
    {
        Directory.CreateDirectory(_storageDir);
        var path = Path.Combine(_storageDir, "a.bin");
        await File.WriteAllTextAsync(path, "hello");
        var file = new FileEntity { Id = 1, Name = "a.bin", Extension = "bin", Size = 5, ContentType = null, Location = path, CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(file);

        var result = await _sut.GetDownloadTargetAsync(1);

        Assert.Equal("application/octet-stream", result.ContentType);
    }

    // ── DeleteAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task DeleteAsync_DeletesEntityAndDiskFile()
    {
        Directory.CreateDirectory(_storageDir);
        var path = Path.Combine(_storageDir, "a.txt");
        await File.WriteAllTextAsync(path, "hello");
        var file = new FileEntity { Id = 1, Name = "a.txt", Extension = "txt", Size = 5, Location = path, CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(file);

        await _sut.DeleteAsync(1);

        _repoMock.Verify(r => r.Delete(file), Times.Once);
        Assert.False(File.Exists(path));
    }

    [Fact]
    public async Task DeleteAsync_ThrowsKeyNotFoundException_WhenNotFound()
    {
        _repoMock.Setup(r => r.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((FileEntity?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.DeleteAsync(99));
    }
}
