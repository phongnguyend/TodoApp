using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Moq;
using TodoApi.Controllers;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Tests.Controllers;

public class FilesControllerTests
{
    private readonly Mock<IFileService> _serviceMock;
    private readonly FilesController _sut;

    public FilesControllerTests()
    {
        _serviceMock = new Mock<IFileService>();
        _sut = new FilesController(_serviceMock.Object);
    }

    // ── GetAll ────────────────────────────────────────────────────────────────

    [Fact]
    public async Task GetAll_Returns200WithPaginatedResponse()
    {
        var paginated = new PaginatedResponse<FileResponse>([], 0, 1, 20, 0);
        _serviceMock.Setup(s => s.GetAllAsync(1, 20, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(paginated);

        var result = await _sut.GetAll();

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Equal(200, ok.StatusCode);
        Assert.Same(paginated, ok.Value);
    }

    [Fact]
    public async Task GetAll_ForwardsPageAndPageSizeToService()
    {
        var paginated = new PaginatedResponse<FileResponse>([], 0, 2, 5, 0);
        _serviceMock.Setup(s => s.GetAllAsync(2, 5, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(paginated);

        await _sut.GetAll(page: 2, pageSize: 5);

        _serviceMock.Verify(s => s.GetAllAsync(2, 5, It.IsAny<CancellationToken>()), Times.Once);
    }

    // ── GetById ───────────────────────────────────────────────────────────────

    [Fact]
    public async Task GetById_Returns200_WhenFound()
    {
        var file = new FileResponse(1, "a.txt", "txt", 10, "text/plain", DateTime.UtcNow, null);
        _serviceMock.Setup(s => s.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(file);

        var result = await _sut.GetById(1);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(file, ok.Value);
    }

    [Fact]
    public async Task GetById_Returns404_WhenNotFound()
    {
        _serviceMock.Setup(s => s.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new KeyNotFoundException("File 99 not found."));

        var result = await _sut.GetById(99);

        var notFound = Assert.IsType<NotFoundObjectResult>(result);
        Assert.Equal(404, notFound.StatusCode);
    }

    // ── Download ──────────────────────────────────────────────────────────────

    [Fact]
    public async Task Download_Returns404_WhenNotFound()
    {
        _serviceMock.Setup(s => s.GetDownloadTargetAsync(99, It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new KeyNotFoundException("File 99 not found."));

        var result = await _sut.Download(99);

        Assert.IsType<NotFoundObjectResult>(result);
    }

    [Fact]
    public async Task Download_ReturnsPhysicalFile_WhenFound()
    {
        var target = new FileDownloadTarget(Path.Combine(Path.GetTempPath(), "a.txt"), "a.txt", "text/plain");
        _serviceMock.Setup(s => s.GetDownloadTargetAsync(1, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(target);

        var result = await _sut.Download(1);

        var fileResult = Assert.IsType<PhysicalFileResult>(result);
        Assert.Equal(target.Path, fileResult.FileName);
        Assert.Equal("text/plain", fileResult.ContentType);
        Assert.Equal("a.txt", fileResult.FileDownloadName);
    }

    // ── Upload ────────────────────────────────────────────────────────────────

    [Fact]
    public async Task Upload_Returns400_WhenFileMissing()
    {
        var result = await _sut.Upload(null);

        var badRequest = Assert.IsType<BadRequestObjectResult>(result);
        Assert.Equal(400, badRequest.StatusCode);
        _serviceMock.Verify(s => s.UploadAsync(It.IsAny<IFormFile>(), It.IsAny<CancellationToken>()), Times.Never);
    }

    [Fact]
    public async Task Upload_Returns201WithCreatedFile()
    {
        var formFile = new Mock<IFormFile>();
        formFile.Setup(f => f.Length).Returns(10);
        var created = new FileResponse(1, "a.txt", "txt", 10, "text/plain", DateTime.UtcNow, null);
        _serviceMock.Setup(s => s.UploadAsync(formFile.Object, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(created);

        var result = await _sut.Upload(formFile.Object);

        var createdResult = Assert.IsType<CreatedAtActionResult>(result);
        Assert.Equal(201, createdResult.StatusCode);
        Assert.Same(created, createdResult.Value);
        Assert.Equal(nameof(FilesController.GetById), createdResult.ActionName);
        Assert.Equal(1, createdResult.RouteValues!["id"]);
    }

    [Fact]
    public async Task Upload_Returns413_WhenFileTooLarge()
    {
        var formFile = new Mock<IFormFile>();
        formFile.Setup(f => f.Length).Returns(999);
        _serviceMock.Setup(s => s.UploadAsync(formFile.Object, It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new FileTooLargeException("File exceeds the maximum allowed size of 999 bytes."));

        var result = await _sut.Upload(formFile.Object);

        var statusResult = Assert.IsType<ObjectResult>(result);
        Assert.Equal(StatusCodes.Status413PayloadTooLarge, statusResult.StatusCode);
    }

    // ── Delete ────────────────────────────────────────────────────────────────

    [Fact]
    public async Task Delete_Returns204_WhenDeleted()
    {
        _serviceMock.Setup(s => s.DeleteAsync(1, It.IsAny<CancellationToken>()))
                    .Returns(Task.CompletedTask);

        var result = await _sut.Delete(1);

        var noContent = Assert.IsType<NoContentResult>(result);
        Assert.Equal(204, noContent.StatusCode);
    }

    [Fact]
    public async Task Delete_Returns404_WhenNotFound()
    {
        _serviceMock.Setup(s => s.DeleteAsync(99, It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new KeyNotFoundException("File 99 not found."));

        var result = await _sut.Delete(99);

        Assert.IsType<NotFoundObjectResult>(result);
    }
}
