using Microsoft.AspNetCore.Mvc;
using Moq;
using TodoApi.Controllers;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Tests.Controllers;

public class TodoItemsControllerTests
{
    private readonly Mock<ITodoItemService> _serviceMock;
    private readonly TodoItemsController _sut;

    public TodoItemsControllerTests()
    {
        _serviceMock = new Mock<ITodoItemService>();
        _sut = new TodoItemsController(_serviceMock.Object);
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
}
