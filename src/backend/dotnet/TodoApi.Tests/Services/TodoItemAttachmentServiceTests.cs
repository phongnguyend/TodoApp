using Microsoft.EntityFrameworkCore;
using Moq;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoApi.Services;
using TodoShared.Models;

namespace TodoApi.Tests.Services;

public class TodoItemAttachmentServiceTests
{
    private readonly Mock<ITodoItemAttachmentRepository> _attachmentRepoMock = new();
    private readonly Mock<ITodoItemRepository> _todoItemRepoMock = new();
    private readonly Mock<IFileRepository> _fileRepoMock = new();
    private readonly AppDbContext _db;
    private readonly TodoItemAttachmentService _sut;

    public TodoItemAttachmentServiceTests()
    {
        var options = new DbContextOptionsBuilder<AppDbContext>()
            .UseInMemoryDatabase(databaseName: Guid.NewGuid().ToString())
            .Options;
        _db = new AppDbContext(options);

        _sut = new TodoItemAttachmentService(
            _attachmentRepoMock.Object,
            _todoItemRepoMock.Object,
            _fileRepoMock.Object,
            _db);
    }

    [Fact]
    public async Task GetAllAsync_ReturnsAttachments_WhenTodoItemExists()
    {
        _todoItemRepoMock.Setup(r => r.GetByIdAsync(10, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItem { Id = 10, Title = "Task", CreatedAt = DateTime.UtcNow });
        _attachmentRepoMock.Setup(r => r.GetByTodoItemIdAsync(10, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new List<TodoItemAttachment>
            {
                new() { Id = 1, TodoItemId = 10, FileId = 5, CreatedAt = DateTime.UtcNow }
            });

        var result = await _sut.GetAllAsync(10);

        Assert.Single(result);
        Assert.Equal(10, result[0].TodoItemId);
        Assert.Equal(5, result[0].FileId);
    }

    [Fact]
    public async Task GetByIdAsync_Throws_WhenAttachmentNotFound()
    {
        _todoItemRepoMock.Setup(r => r.GetByIdAsync(10, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItem { Id = 10, Title = "Task", CreatedAt = DateTime.UtcNow });
        _attachmentRepoMock.Setup(r => r.GetByIdForTodoItemAsync(10, 99, It.IsAny<CancellationToken>()))
            .ReturnsAsync((TodoItemAttachment?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.GetByIdAsync(10, 99));
    }

    [Fact]
    public async Task CreateAsync_CreatesAttachment_WhenTodoItemAndFileExist()
    {
        _todoItemRepoMock.Setup(r => r.GetByIdAsync(10, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItem { Id = 10, Title = "Task", CreatedAt = DateTime.UtcNow });
        _fileRepoMock.Setup(r => r.GetByIdAsync(5, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new FileEntity { Id = 5, Name = "a.txt", Extension = "txt", Size = 1, Location = "/tmp/a.txt", CreatedAt = DateTime.UtcNow });
        _attachmentRepoMock.Setup(r => r.GetByTodoItemAndFileAsync(10, 5, It.IsAny<CancellationToken>()))
            .ReturnsAsync((TodoItemAttachment?)null);
        _attachmentRepoMock.Setup(r => r.AddAsync(It.IsAny<TodoItemAttachment>(), It.IsAny<CancellationToken>()))
            .ReturnsAsync((TodoItemAttachment attachment, CancellationToken _) => attachment);

        var result = await _sut.CreateAsync(10, new CreateTodoItemAttachmentRequest(5));

        Assert.Equal(10, result.TodoItemId);
        Assert.Equal(5, result.FileId);
        _attachmentRepoMock.Verify(r => r.AddAsync(It.IsAny<TodoItemAttachment>(), It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task UpdateAsync_UpdatesAttachmentFile_WhenValid()
    {
        _todoItemRepoMock.Setup(r => r.GetByIdAsync(10, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItem { Id = 10, Title = "Task", CreatedAt = DateTime.UtcNow });
        _fileRepoMock.Setup(r => r.GetByIdAsync(6, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new FileEntity { Id = 6, Name = "b.txt", Extension = "txt", Size = 2, Location = "/tmp/b.txt", CreatedAt = DateTime.UtcNow });
        _attachmentRepoMock.Setup(r => r.GetByIdForTodoItemAsync(10, 3, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItemAttachment { Id = 3, TodoItemId = 10, FileId = 5, CreatedAt = DateTime.UtcNow });
        _attachmentRepoMock.Setup(r => r.GetByTodoItemAndFileAsync(10, 6, It.IsAny<CancellationToken>()))
            .ReturnsAsync((TodoItemAttachment?)null);
        _attachmentRepoMock.Setup(r => r.Update(It.IsAny<TodoItemAttachment>()))
            .Returns((TodoItemAttachment attachment) => attachment);

        var result = await _sut.UpdateAsync(10, 3, new CreateTodoItemAttachmentRequest(6));

        Assert.Equal(6, result.FileId);
        _attachmentRepoMock.Verify(r => r.Update(It.IsAny<TodoItemAttachment>()), Times.Once);
    }

    [Fact]
    public async Task DeleteAsync_RemovesAttachment_WhenFound()
    {
        _todoItemRepoMock.Setup(r => r.GetByIdAsync(10, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItem { Id = 10, Title = "Task", CreatedAt = DateTime.UtcNow });
        _attachmentRepoMock.Setup(r => r.GetByIdForTodoItemAsync(10, 3, It.IsAny<CancellationToken>()))
            .ReturnsAsync(new TodoItemAttachment { Id = 3, TodoItemId = 10, FileId = 5, CreatedAt = DateTime.UtcNow });

        await _sut.DeleteAsync(10, 3);

        _attachmentRepoMock.Verify(r => r.Delete(It.IsAny<TodoItemAttachment>()), Times.Once);
    }
}
