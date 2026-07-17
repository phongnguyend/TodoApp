using System.Text;
using Microsoft.AspNetCore.Http;
using Microsoft.EntityFrameworkCore;
using DocumentFormat.OpenXml.Packaging;
using DocumentFormat.OpenXml.Spreadsheet;
using Moq;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoApi.Services;
using TodoShared.Models;

namespace TodoApi.Tests.Services;

public class TodoItemServiceTests
{
    private readonly Mock<ITodoItemRepository> _repoMock;
    private readonly AppDbContext _db;
    private readonly TodoItemService _sut;

    public TodoItemServiceTests()
    {
        _repoMock = new Mock<ITodoItemRepository>();

        var options = new DbContextOptionsBuilder<AppDbContext>()
            .UseInMemoryDatabase(databaseName: Guid.NewGuid().ToString())
            .Options;
        _db = new AppDbContext(options);

        _sut = new TodoItemService(_repoMock.Object, _db);
    }

    // ── GetAllAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task GetAllAsync_ReturnsPaginatedResponse()
    {
        var items = new List<TodoItem>
        {
            new() { Id = 1, Title = "Task 1", IsCompleted = false, CreatedAt = DateTime.UtcNow },
            new() { Id = 2, Title = "Task 2", IsCompleted = true,  CreatedAt = DateTime.UtcNow },
        };
        _repoMock.Setup(r => r.GetAllAsync(0, 20, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((items, 2));

        var result = await _sut.GetAllAsync(1, 20);

        Assert.Equal(2, result.Total);
        Assert.Equal(1, result.Page);
        Assert.Equal(20, result.PageSize);
        Assert.Equal(1, result.TotalPages);
        Assert.Equal(2, result.Items.Count());
    }

    [Fact]
    public async Task GetAllAsync_CalculatesSkipFromPage()
    {
        _repoMock.Setup(r => r.GetAllAsync(20, 10, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((Enumerable.Empty<TodoItem>(), 0));

        await _sut.GetAllAsync(page: 3, pageSize: 10);

        _repoMock.Verify(r => r.GetAllAsync(20, 10, It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task GetAllAsync_CalculatesTotalPages()
    {
        var items = Enumerable.Range(1, 5).Select(i => new TodoItem { Id = i, Title = $"Task {i}", CreatedAt = DateTime.UtcNow }).ToList();
        _repoMock.Setup(r => r.GetAllAsync(0, 3, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((items.Take(3), 5));

        var result = await _sut.GetAllAsync(1, 3);

        Assert.Equal(2, result.TotalPages);
    }

    // ── GetIncompleteAsync ────────────────────────────────────────────────────

    [Fact]
    public async Task GetIncompleteAsync_ReturnsPaginatedResponse()
    {
        var items = new List<TodoItem>
        {
            new() { Id = 1, Title = "Task 1", IsCompleted = false, CreatedAt = DateTime.UtcNow },
        };
        _repoMock.Setup(r => r.GetIncompleteAsync(0, 20, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((items, 1));

        var result = await _sut.GetIncompleteAsync(1, 20);

        Assert.Equal(1, result.Total);
        Assert.Single(result.Items);
    }

    // ── GetByIdAsync ──────────────────────────────────────────────────────────

    [Fact]
    public async Task GetByIdAsync_ReturnsItem_WhenFound()
    {
        var item = new TodoItem { Id = 1, Title = "Task 1", CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(item);

        var result = await _sut.GetByIdAsync(1);

        Assert.Equal(1, result.Id);
        Assert.Equal("Task 1", result.Title);
    }

    [Fact]
    public async Task GetByIdAsync_ThrowsKeyNotFoundException_WhenNotFound()
    {
        _repoMock.Setup(r => r.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.GetByIdAsync(99));
    }

    // ── CreateAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task CreateAsync_AddsItemAndSavesChanges()
    {
        var request = new CreateTodoItemRequest("New Task", "Some description");
        var savedItem = new TodoItem { Id = 1, Title = "New Task", Description = "Some description", CreatedAt = DateTime.UtcNow };

        _repoMock.Setup(r => r.AddAsync(It.IsAny<TodoItem>(), It.IsAny<CancellationToken>()))
                 .ReturnsAsync(savedItem);

        var result = await _sut.CreateAsync(request);

        Assert.Equal(1, result.Id);
        Assert.Equal("New Task", result.Title);
        Assert.Equal("Some description", result.Description);
        _repoMock.Verify(r => r.AddAsync(It.IsAny<TodoItem>(), It.IsAny<CancellationToken>()), Times.Once);
    }

    [Fact]
    public async Task CreateAsync_SetsCreatedAtToUtcNow()
    {
        var request = new CreateTodoItemRequest("Task", null);
        var before = DateTime.UtcNow;

        _repoMock.Setup(r => r.AddAsync(It.IsAny<TodoItem>(), It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem item, CancellationToken _) => item);

        await _sut.CreateAsync(request);

        _repoMock.Verify(r => r.AddAsync(
            It.Is<TodoItem>(t => t.CreatedAt >= before && t.CreatedAt <= DateTime.UtcNow),
            It.IsAny<CancellationToken>()), Times.Once);
    }

    // ── UpdateAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task UpdateAsync_UpdatesProvidedFields()
    {
        var item = new TodoItem { Id = 1, Title = "Old Title", IsCompleted = false, CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(item);
        _repoMock.Setup(r => r.Update(It.IsAny<TodoItem>())).Returns(item);

        var request = new UpdateTodoItemRequest("New Title", "Desc", true);
        var result = await _sut.UpdateAsync(1, request);

        Assert.Equal("New Title", result.Title);
        Assert.Equal("Desc", result.Description);
        Assert.True(result.IsCompleted);
        _repoMock.Verify(r => r.Update(It.IsAny<TodoItem>()), Times.Once);
    }

    [Fact]
    public async Task UpdateAsync_DoesNotOverwriteNullFields()
    {
        var item = new TodoItem { Id = 1, Title = "Original", Description = "Orig desc", IsCompleted = false, CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(item);
        _repoMock.Setup(r => r.Update(It.IsAny<TodoItem>())).Returns(item);

        var request = new UpdateTodoItemRequest(null, null, null);
        var result = await _sut.UpdateAsync(1, request);

        Assert.Equal("Original", result.Title);
        Assert.Equal("Orig desc", result.Description);
        Assert.False(result.IsCompleted);
    }

    [Fact]
    public async Task UpdateAsync_ThrowsKeyNotFoundException_WhenNotFound()
    {
        _repoMock.Setup(r => r.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(
            () => _sut.UpdateAsync(99, new UpdateTodoItemRequest("X", null, null)));
    }

    // ── DeleteAsync ───────────────────────────────────────────────────────────

    [Fact]
    public async Task DeleteAsync_DeletesItem()
    {
        var item = new TodoItem { Id = 1, Title = "Task", CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(item);

        await _sut.DeleteAsync(1);

        _repoMock.Verify(r => r.Delete(item), Times.Once);
    }

    [Fact]
    public async Task DeleteAsync_ThrowsKeyNotFoundException_WhenNotFound()
    {
        _repoMock.Setup(r => r.GetByIdAsync(5, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.DeleteAsync(5));
    }

    // ── MarkCompleteAsync ─────────────────────────────────────────────────────

    [Fact]
    public async Task MarkCompleteAsync_SetsIsCompletedTrue()
    {
        var item = new TodoItem { Id = 1, Title = "Task", IsCompleted = false, CreatedAt = DateTime.UtcNow };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(item);
        _repoMock.Setup(r => r.Update(It.IsAny<TodoItem>())).Returns(item);

        var result = await _sut.MarkCompleteAsync(1);

        Assert.True(result.IsCompleted);
        _repoMock.Verify(r => r.Update(It.IsAny<TodoItem>()), Times.Once);
    }

    [Fact]
    public async Task MarkCompleteAsync_ThrowsKeyNotFoundException_WhenNotFound()
    {
        _repoMock.Setup(r => r.GetByIdAsync(99, It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem?)null);

        await Assert.ThrowsAsync<KeyNotFoundException>(() => _sut.MarkCompleteAsync(99));
    }

    [Fact]
    public async Task MarkCompleteAsync_SetsUpdatedAt()
    {
        var item = new TodoItem { Id = 1, Title = "Task", IsCompleted = false, CreatedAt = DateTime.UtcNow, UpdatedAt = null };
        _repoMock.Setup(r => r.GetByIdAsync(1, It.IsAny<CancellationToken>()))
                 .ReturnsAsync(item);
        _repoMock.Setup(r => r.Update(It.IsAny<TodoItem>())).Returns(item);

        var before = DateTime.UtcNow;
        await _sut.MarkCompleteAsync(1);

        Assert.NotNull(item.UpdatedAt);
        Assert.True(item.UpdatedAt >= before);
    }

    // ── ImportCsvAsync ──────────────────────────────────────────────────────

    [Fact]
    public async Task ImportCsvAsync_CreatesValidRowsAndCollectsErrors()
    {
        var csv = "title,description,is_completed\nBuy milk,Whole milk,true\n,Missing title,false\n";
        using var stream = new MemoryStream(Encoding.UTF8.GetBytes(csv));
        var file = new FormFile(stream, 0, stream.Length, "file", "todo-items.csv")
        {
            Headers = new HeaderDictionary(),
            ContentType = "text/csv"
        };

        _repoMock.Setup(r => r.AddAsync(It.IsAny<TodoItem>(), It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem item, CancellationToken _) => item);

        var result = await _sut.ImportCsvAsync(file);

        Assert.Equal(1, result.Imported);
        Assert.Equal(1, result.Failed);
        Assert.Single(result.Errors);
        Assert.Equal(3, result.Errors[0].Row);
    }

    // ── ExportCsvAsync ─────────────────────────────────────────────────────

    [Fact]
    public async Task ExportCsvAsync_ReturnsCsvContent()
    {
        var items = new List<TodoItem>
        {
            new() { Id = 1, Title = "Buy milk", Description = "Whole milk", IsCompleted = false, CreatedAt = DateTime.UtcNow },
        };
        _repoMock.Setup(r => r.GetAllItemsAsync(It.IsAny<CancellationToken>()))
                 .ReturnsAsync(items);

        var csv = await _sut.ExportCsvAsync();

        Assert.StartsWith("id,title,description,is_completed,created_at,updated_at", csv);
        Assert.Contains("Buy milk", csv);
    }

    // ── ImportExcelAsync ────────────────────────────────────────────────────

    [Fact]
    public async Task ImportExcelAsync_CreatesValidRowsAndCollectsErrors()
    {
        using var stream = new MemoryStream();
        using (var package = SpreadsheetDocument.Create(stream, DocumentFormat.OpenXml.SpreadsheetDocumentType.Workbook, true))
        {
            var workbookPart = package.AddWorkbookPart();
            workbookPart.Workbook = new Workbook(new Sheets());

            var worksheetPart = workbookPart.AddNewPart<WorksheetPart>();
            worksheetPart.Worksheet = new Worksheet(new SheetData());
            var sheetData = worksheetPart.Worksheet.GetFirstChild<SheetData>()!;

            AppendRow(sheetData, new[] { "title", "description", "is_completed" });
            AppendRow(sheetData, new[] { "Buy milk", "Whole milk", "true" });
            AppendRow(sheetData, new[] { string.Empty, "Missing title", "false" });

            var sheets = workbookPart.Workbook.GetFirstChild<Sheets>()!;
            sheets.Append(new Sheet { Id = workbookPart.GetIdOfPart(worksheetPart), SheetId = 1, Name = "Todo Items" });
            workbookPart.Workbook.Save();
        }

        stream.Position = 0;

        var file = new FormFile(stream, 0, stream.Length, "file", "todo-items.xlsx")
        {
            Headers = new HeaderDictionary(),
            ContentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
        };

        _repoMock.Setup(r => r.AddAsync(It.IsAny<TodoItem>(), It.IsAny<CancellationToken>()))
                 .ReturnsAsync((TodoItem item, CancellationToken _) => item);

        var result = await _sut.ImportExcelAsync(file);

        Assert.Equal(1, result.Imported);
        Assert.Equal(1, result.Failed);
        Assert.Single(result.Errors);
        Assert.Equal(3, result.Errors[0].Row);
    }

    // ── ExportExcelAsync ────────────────────────────────────────────────────

    [Fact]
    public async Task ExportExcelAsync_ReturnsWorkbookBytes()
    {
        var items = new List<TodoItem>
        {
            new() { Id = 1, Title = "Buy milk", Description = "Whole milk", IsCompleted = false, CreatedAt = DateTime.UtcNow },
        };
        _repoMock.Setup(r => r.GetAllItemsAsync(It.IsAny<CancellationToken>()))
                 .ReturnsAsync(items);

        var bytes = await _sut.ExportExcelAsync();

        using var stream = new MemoryStream(bytes);
        using var package = SpreadsheetDocument.Open(stream, false);
        var workbookPart = package.WorkbookPart!;
        var sheet = workbookPart.Workbook.Descendants<Sheet>().First();
        var worksheetPart = (WorksheetPart)workbookPart.GetPartById(sheet.Id!);
        var rows = worksheetPart.Worksheet.Descendants<Row>().ToList();

        Assert.Equal("Todo Items", sheet.Name);
        Assert.Equal("id", GetCellValue(rows[0].Elements<Cell>().ElementAt(0), workbookPart.SharedStringTablePart));
        Assert.Equal("Buy milk", GetCellValue(rows[1].Elements<Cell>().ElementAt(1), workbookPart.SharedStringTablePart));
    }

    private static void AppendRow(SheetData sheetData, IEnumerable<string> values)
    {
        var row = new Row();
        foreach (var value in values)
        {
            row.AppendChild(new Cell(new InlineString(new Text(value))) { DataType = CellValues.InlineString });
        }

        sheetData.AppendChild(row);
    }

    private static string GetCellValue(Cell cell, SharedStringTablePart? sharedStringTablePart)
    {
        if (cell.CellValue is null)
        {
            return string.Empty;
        }

        var text = cell.CellValue.Text;
        if (cell.DataType?.Value == CellValues.SharedString && sharedStringTablePart is not null)
        {
            if (int.TryParse(text, out var index) && index >= 0 && index < sharedStringTablePart.SharedStringTable.Count())
            {
                return sharedStringTablePart.SharedStringTable.ElementAt(index).InnerText;
            }
        }

        return text;
    }
}
