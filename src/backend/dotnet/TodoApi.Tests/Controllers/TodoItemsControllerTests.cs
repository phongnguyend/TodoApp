using Microsoft.AspNetCore.Http;
using Microsoft.AspNetCore.Mvc;
using Moq;
using TodoApi.Controllers;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Tests.Controllers;

public class TodoItemsControllerTests
{
    private readonly Mock<ITodoItemService> _serviceMock;
    private readonly Mock<ITodoItemAttachmentService> _attachmentServiceMock;
    private readonly TodoItemsController _sut;

    public TodoItemsControllerTests()
    {
        _serviceMock = new Mock<ITodoItemService>();
        _attachmentServiceMock = new Mock<ITodoItemAttachmentService>();
        _sut = new TodoItemsController(_serviceMock.Object, _attachmentServiceMock.Object);
    }

    // ── GetAll ────────────────────────────────────────────────────────────────

    [Fact]
    public async Task GetAll_Returns200WithPaginatedResponse()
    {
        var paginated = new PaginatedResponse<TodoItemResponse>([], 0, 1, 20, 0);
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
        var paginated = new PaginatedResponse<TodoItemResponse>([], 0, 2, 5, 0);
        _serviceMock.Setup(s => s.GetAllAsync(2, 5, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(paginated);

        await _sut.GetAll(page: 2, pageSize: 5);

        _serviceMock.Verify(s => s.GetAllAsync(2, 5, It.IsAny<CancellationToken>()), Times.Once);
    }

    // ── GetIncomplete ─────────────────────────────────────────────────────────

    [Fact]
    public async Task GetIncomplete_Returns200WithPaginatedResponse()
    {
        var paginated = new PaginatedResponse<TodoItemResponse>([], 0, 1, 20, 0);
        _serviceMock.Setup(s => s.GetIncompleteAsync(1, 20, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(paginated);

        var result = await _sut.GetIncomplete();

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(paginated, ok.Value);
    }

    // ── GetById ───────────────────────────────────────────────────────────────

    [Fact]
    public async Task GetById_Returns200_WhenFound()
    {
        var item = new TodoItemResponse(1, "Task", null, false, DateTime.UtcNow, null);
        _serviceMock.Setup(s => s.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(item);

        var result = await _sut.GetById(1);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(item, ok.Value);
    }

    [Fact]
    public async Task GetById_Returns404_WhenNotFound()
    {
        _serviceMock.Setup(s => s.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new KeyNotFoundException("Todo item 99 not found."));

        var result = await _sut.GetById(99);

        var notFound = Assert.IsType<NotFoundObjectResult>(result);
        Assert.Equal(404, notFound.StatusCode);
    }

    // ── Attachments ────────────────────────────────────────────────────────

    [Fact]
    public async Task GetAttachments_Returns200WithAttachments()
    {
        var attachments = new List<TodoItemAttachmentResponse>
        {
            new(1, 10, 1, DateTime.UtcNow, null)
        };
        _attachmentServiceMock.Setup(s => s.GetAllAsync(10, It.IsAny<CancellationToken>()))
                               .ReturnsAsync(attachments);

        var result = await _sut.GetAttachments(10);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Equal(200, ok.StatusCode);
        Assert.Same(attachments, ok.Value);
    }

    [Fact]
    public async Task CreateAttachment_Returns201WithCreatedAttachment()
    {
        var request = new CreateTodoItemAttachmentRequest(5);
        var created = new TodoItemAttachmentResponse(1, 10, 5, DateTime.UtcNow, null);
        _attachmentServiceMock.Setup(s => s.CreateAsync(10, request, It.IsAny<CancellationToken>()))
                               .ReturnsAsync(created);

        var result = await _sut.CreateAttachment(10, request);

        var createdResult = Assert.IsType<CreatedAtActionResult>(result);
        Assert.Equal(201, createdResult.StatusCode);
        Assert.Same(created, createdResult.Value);
        Assert.Equal(nameof(TodoItemsController.GetAttachmentById), createdResult.ActionName);
        Assert.Equal(10, createdResult.RouteValues!["id"]);
        Assert.Equal(1, createdResult.RouteValues!["attachmentId"]);
    }

    // ── Create ────────────────────────────────────────────────────────────────

    [Fact]
    public async Task Create_Returns201WithCreatedItem()
    {
        var request = new CreateTodoItemRequest("New Task", null);
        var created = new TodoItemResponse(1, "New Task", null, false, DateTime.UtcNow, null);
        _serviceMock.Setup(s => s.CreateAsync(request, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(created);

        var result = await _sut.Create(request);

        var createdResult = Assert.IsType<CreatedAtActionResult>(result);
        Assert.Equal(201, createdResult.StatusCode);
        Assert.Same(created, createdResult.Value);
        Assert.Equal(nameof(TodoItemsController.GetById), createdResult.ActionName);
        Assert.Equal(1, createdResult.RouteValues!["id"]);
    }

    // ── Update ────────────────────────────────────────────────────────────────

    [Fact]
    public async Task Update_Returns200WithUpdatedItem()
    {
        var request = new UpdateTodoItemRequest("Updated", null, null);
        var updated = new TodoItemResponse(1, "Updated", null, false, DateTime.UtcNow, DateTime.UtcNow);
        _serviceMock.Setup(s => s.UpdateAsync(1, request, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(updated);

        var result = await _sut.Update(1, request);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Same(updated, ok.Value);
    }

    [Fact]
    public async Task Update_Returns404_WhenNotFound()
    {
        _serviceMock.Setup(s => s.UpdateAsync(99, It.IsAny<UpdateTodoItemRequest>(), It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new KeyNotFoundException("Todo item 99 not found."));

        var result = await _sut.Update(99, new UpdateTodoItemRequest(null, null, null));

        Assert.IsType<NotFoundObjectResult>(result);
    }

    // ── MarkComplete ──────────────────────────────────────────────────────────

    [Fact]
    public async Task MarkComplete_Returns200WithCompletedItem()
    {
        var completed = new TodoItemResponse(1, "Task", null, true, DateTime.UtcNow, DateTime.UtcNow);
        _serviceMock.Setup(s => s.MarkCompleteAsync(1, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(completed);

        var result = await _sut.MarkComplete(1);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.True(((TodoItemResponse)ok.Value!).IsCompleted);
    }

    [Fact]
    public async Task MarkComplete_Returns404_WhenNotFound()
    {
        _serviceMock.Setup(s => s.MarkCompleteAsync(99, It.IsAny<CancellationToken>()))
                    .ThrowsAsync(new KeyNotFoundException("Todo item 99 not found."));

        var result = await _sut.MarkComplete(99);

        Assert.IsType<NotFoundObjectResult>(result);
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
                    .ThrowsAsync(new KeyNotFoundException("Todo item 99 not found."));

        var result = await _sut.Delete(99);

        Assert.IsType<NotFoundObjectResult>(result);
    }

    // ── CSV import/export ───────────────────────────────────────────────────

    [Fact]
    public async Task ImportCsv_Returns400_WhenFileMissing()
    {
        var result = await _sut.ImportCsv(null);

        var badRequest = Assert.IsType<BadRequestObjectResult>(result);
        Assert.Equal(400, badRequest.StatusCode);
        _serviceMock.Verify(s => s.ImportCsvAsync(It.IsAny<IFormFile>(), It.IsAny<CancellationToken>()), Times.Never);
    }

    [Fact]
    public async Task ImportCsv_Returns200_WithImportResult()
    {
        var formFile = new Mock<IFormFile>();
        formFile.Setup(f => f.Length).Returns(10);
        var importResult = new ImportResult(2, 1, [new ImportRowError(3, "Title is required.")]);
        _serviceMock.Setup(s => s.ImportCsvAsync(formFile.Object, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(importResult);

        var result = await _sut.ImportCsv(formFile.Object);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Equal(200, ok.StatusCode);
        Assert.Same(importResult, ok.Value);
    }

    [Fact]
    public async Task ExportCsv_ReturnsFileResult()
    {
        _serviceMock.Setup(s => s.ExportCsvAsync(It.IsAny<CancellationToken>()))
                    .ReturnsAsync("id,title\n1,Test");

        var result = await _sut.ExportCsv();

        var file = Assert.IsType<FileContentResult>(result);
        Assert.Equal("text/csv; charset=utf-8", file.ContentType);
        Assert.Equal("todo-items.csv", file.FileDownloadName);
        Assert.Equal("id,title\n1,Test", System.Text.Encoding.UTF8.GetString(file.FileContents));
    }

    [Fact]
    public async Task ImportExcel_Returns400_WhenFileMissing()
    {
        var result = await _sut.ImportExcel(null);

        var badRequest = Assert.IsType<BadRequestObjectResult>(result);
        Assert.Equal(400, badRequest.StatusCode);
        _serviceMock.Verify(s => s.ImportExcelAsync(It.IsAny<IFormFile>(), It.IsAny<CancellationToken>()), Times.Never);
    }

    [Fact]
    public async Task ImportExcel_Returns200_WithImportResult()
    {
        var formFile = new Mock<IFormFile>();
        formFile.Setup(f => f.Length).Returns(10);
        var importResult = new ImportResult(1, 0, []);
        _serviceMock.Setup(s => s.ImportExcelAsync(formFile.Object, It.IsAny<CancellationToken>()))
                    .ReturnsAsync(importResult);

        var result = await _sut.ImportExcel(formFile.Object);

        var ok = Assert.IsType<OkObjectResult>(result);
        Assert.Equal(200, ok.StatusCode);
        Assert.Same(importResult, ok.Value);
    }

    [Fact]
    public async Task ExportExcel_ReturnsFileResult()
    {
        _serviceMock.Setup(s => s.ExportExcelAsync(It.IsAny<CancellationToken>()))
                    .ReturnsAsync(new byte[] { 1, 2, 3 });

        var result = await _sut.ExportExcel();

        var file = Assert.IsType<FileContentResult>(result);
        Assert.Equal("application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file.ContentType);
        Assert.Equal("todo-items.xlsx", file.FileDownloadName);
        Assert.Equal(new byte[] { 1, 2, 3 }, file.FileContents);
    }
}
