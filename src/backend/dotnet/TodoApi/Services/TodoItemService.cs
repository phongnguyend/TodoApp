using System.Text;
using DocumentFormat.OpenXml;
using DocumentFormat.OpenXml.Packaging;
using DocumentFormat.OpenXml.Spreadsheet;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoShared.Models;

namespace TodoApi.Services;

public class TodoItemService(
    ITodoItemRepository repository,
    AppDbContext db,
    IHttpContextAccessor? httpContextAccessor = null) : ITodoItemService
{
    private int? ActorUserId => AuditActor.GetUserId(httpContextAccessor);
    private static readonly string[] CsvFieldNames = ["id", "title", "description", "is_completed", "created_at", "updated_at"];
    private static readonly string[] ExcelFieldNames = ["id", "title", "description", "is_completed", "created_at", "updated_at"];
    private static readonly HashSet<string> TrueValues = new(StringComparer.OrdinalIgnoreCase) { "1", "true", "yes", "y" };

    // ── Mapping ───────────────────────────────────────────────────────────────

    private static TodoItemResponse ToResponse(TodoItem item) =>
        new(item.Id, item.Title, item.Description, item.IsCompleted, item.CreatedAt,
            item.CreatedByUserId, item.UpdatedAt, item.UpdatedByUserId);

    private static PaginatedResponse<TodoItemResponse> ToPaginated(
        IEnumerable<TodoItem> items, int total, int page, int pageSize) =>
        new(
            items.Select(ToResponse),
            total,
            page,
            pageSize,
            (int)Math.Ceiling(total / (double)pageSize)
        );

    private async Task<TodoItem> GetOrThrowAsync(int id, CancellationToken ct)
    {
        var item = await repository.GetByIdAsync(id, ct)
            ?? throw new KeyNotFoundException($"Todo item {id} not found.");
        return item;
    }

    // ── Queries ───────────────────────────────────────────────────────────────

    public async Task<PaginatedResponse<TodoItemResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var (items, total) = await repository.GetAllAsync((page - 1) * pageSize, pageSize, ct);
        return ToPaginated(items, total, page, pageSize);
    }

    public async Task<PaginatedResponse<TodoItemResponse>> GetIncompleteAsync(int page, int pageSize, CancellationToken ct = default)
    {
        var (items, total) = await repository.GetIncompleteAsync((page - 1) * pageSize, pageSize, ct);
        return ToPaginated(items, total, page, pageSize);
    }

    public async Task<TodoItemResponse> GetByIdAsync(int id, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        return ToResponse(item);
    }

    // ── Commands ──────────────────────────────────────────────────────────────

    public async Task<TodoItemResponse> CreateAsync(CreateTodoItemRequest request, CancellationToken ct = default)
    {
        var item = new TodoItem
        {
            Title = request.Title,
            Description = request.Description,
            CreatedAt = DateTime.UtcNow,
            CreatedByUserId = ActorUserId,
        };
        var created = await repository.AddAsync(item, ct);
        await db.SaveChangesAsync(ct);
        return ToResponse(created);
    }

    public async Task<TodoItemResponse> UpdateAsync(int id, UpdateTodoItemRequest request, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        if (request.Title is not null) item.Title = request.Title;
        if (request.Description is not null) item.Description = request.Description;
        if (request.IsCompleted is not null) item.IsCompleted = request.IsCompleted.Value;
        item.UpdatedAt = DateTime.UtcNow;
        item.UpdatedByUserId = ActorUserId;
        repository.Update(item);
        await db.SaveChangesAsync(ct);
        return ToResponse(item);
    }

    public async Task DeleteAsync(int id, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        repository.Delete(item);
        await db.SaveChangesAsync(ct);
    }

    public async Task<TodoItemResponse> MarkCompleteAsync(int id, CancellationToken ct = default)
    {
        var item = await GetOrThrowAsync(id, ct);
        item.IsCompleted = true;
        item.UpdatedAt = DateTime.UtcNow;
        item.UpdatedByUserId = ActorUserId;
        repository.Update(item);
        await db.SaveChangesAsync(ct);
        return ToResponse(item);
    }

    // ── CSV import/export ─────────────────────────────────────────────────────

    public async Task<ImportResult> ImportCsvAsync(IFormFile file, CancellationToken ct = default)
    {
        if (file.Length == 0)
        {
            return new ImportResult(0, 0, []);
        }

        using var stream = file.OpenReadStream();
        using var reader = new StreamReader(stream, Encoding.UTF8);

        var headerLine = await reader.ReadLineAsync(ct);
        if (string.IsNullOrWhiteSpace(headerLine))
        {
            return new ImportResult(0, 0, []);
        }

        var headers = ParseCsvLine(headerLine).Select(value => value.Trim().ToLowerInvariant()).ToList();
        var indexByName = headers.Select((name, index) => new { name, index })
            .Where(entry => !string.IsNullOrWhiteSpace(entry.name))
            .ToDictionary(entry => entry.name, entry => entry.index, StringComparer.OrdinalIgnoreCase);

        var imported = 0;
        var errors = new List<ImportRowError>();

        string? line;
        var rowNumber = 2;
        while ((line = await reader.ReadLineAsync(ct)) is not null)
        {
            if (string.IsNullOrWhiteSpace(line))
            {
                continue;
            }

            var values = ParseCsvLine(line);
            var title = GetColumnValue(values, indexByName, "title")?.Trim();
            if (string.IsNullOrWhiteSpace(title))
            {
                errors.Add(new ImportRowError(rowNumber, "Title is required."));
                rowNumber++;
                continue;
            }

            var description = GetColumnValue(values, indexByName, "description")?.Trim();
            var isCompleted = ParseBool(GetColumnValue(values, indexByName, "is_completed"));

            var item = new TodoItem
            {
                Title = title,
                Description = string.IsNullOrWhiteSpace(description) ? null : description,
                IsCompleted = isCompleted,
                CreatedAt = DateTime.UtcNow,
                CreatedByUserId = ActorUserId,
            };

            await repository.AddAsync(item, ct);
            imported++;
            rowNumber++;
        }

        await db.SaveChangesAsync(ct);
        return new ImportResult(imported, errors.Count, errors);
    }

    public async Task<string> ExportCsvAsync(CancellationToken ct = default)
    {
        var items = await repository.GetAllItemsAsync(ct);
        var builder = new StringBuilder();

        builder.AppendLine(string.Join(",", CsvFieldNames.Select(EscapeCsvValue)));
        foreach (var item in items)
        {
            builder.AppendLine(string.Join(",", new[]
            {
                item.Id.ToString(),
                item.Title ?? string.Empty,
                item.Description ?? string.Empty,
                item.IsCompleted.ToString().ToLowerInvariant(),
                item.CreatedAt.ToString("O"),
                item.UpdatedAt?.ToString("O") ?? string.Empty,
            }.Select(EscapeCsvValue)));
        }

        return builder.ToString();
    }

    private static string? GetColumnValue(IReadOnlyList<string> values, IReadOnlyDictionary<string, int> indexByName, string columnName)
    {
        if (!indexByName.TryGetValue(columnName, out var index) || index >= values.Count)
        {
            return null;
        }

        return values[index];
    }

    private static bool ParseBool(string? raw)
    {
        if (string.IsNullOrWhiteSpace(raw))
        {
            return false;
        }

        return TrueValues.Contains(raw.Trim());
    }

    private static IReadOnlyList<string> ParseCsvLine(string line)
    {
        var values = new List<string>();
        var current = new StringBuilder();
        var inQuotes = false;

        for (var i = 0; i < line.Length; i++)
        {
            var character = line[i];
            if (character == '"')
            {
                if (inQuotes && i + 1 < line.Length && line[i + 1] == '"')
                {
                    current.Append('"');
                    i++;
                }
                else
                {
                    inQuotes = !inQuotes;
                }
            }
            else if (character == ',' && !inQuotes)
            {
                values.Add(current.ToString());
                current.Clear();
            }
            else
            {
                current.Append(character);
            }
        }

        values.Add(current.ToString());
        return values;
    }

    private static string EscapeCsvValue(string value)
    {
        if (string.IsNullOrEmpty(value))
        {
            return string.Empty;
        }

        if (value.Contains(',') || value.Contains('"') || value.Contains('\n') || value.Contains('\r'))
        {
            return $"\"{value.Replace("\"", "\"\"")}\"";
        }

        return value;
    }

    public async Task<ImportResult> ImportExcelAsync(IFormFile file, CancellationToken ct = default)
    {
        if (file.Length == 0)
        {
            return new ImportResult(0, 0, []);
        }

        using var stream = file.OpenReadStream();
        using var package = SpreadsheetDocument.Open(stream, false);
        var workbookPart = package.WorkbookPart ?? throw new InvalidOperationException("The uploaded file is not a valid workbook.");
        var sheet = workbookPart.Workbook.Descendants<Sheet>().FirstOrDefault();
        if (sheet is null)
        {
            return new ImportResult(0, 0, []);
        }

        var worksheetPart = (WorksheetPart)workbookPart.GetPartById(sheet.Id!);
        var sharedStringPart = workbookPart.GetPartsOfType<SharedStringTablePart>().FirstOrDefault();
        var rows = worksheetPart.Worksheet.Descendants<Row>().ToList();
        if (rows.Count == 0)
        {
            return new ImportResult(0, 0, []);
        }

        var headerCells = ReadRowValues(rows[0], sharedStringPart).Select(value => value.Trim().ToLowerInvariant()).ToList();
        var indexByName = headerCells.Select((name, index) => new { name, index })
            .Where(entry => !string.IsNullOrWhiteSpace(entry.name))
            .ToDictionary(entry => entry.name, entry => entry.index, StringComparer.OrdinalIgnoreCase);

        var imported = 0;
        var errors = new List<ImportRowError>();

        for (var rowIndex = 1; rowIndex < rows.Count; rowIndex++)
        {
            var rowNumber = rowIndex + 1;
            var values = ReadRowValues(rows[rowIndex], sharedStringPart);
            var title = GetColumnValue(values, indexByName, "title")?.Trim();
            if (string.IsNullOrWhiteSpace(title))
            {
                errors.Add(new ImportRowError(rowNumber, "Title is required."));
                continue;
            }

            var description = GetColumnValue(values, indexByName, "description")?.Trim();
            var isCompleted = ParseBool(GetColumnValue(values, indexByName, "is_completed"));

            var item = new TodoItem
            {
                Title = title,
                Description = string.IsNullOrWhiteSpace(description) ? null : description,
                IsCompleted = isCompleted,
                CreatedAt = DateTime.UtcNow,
                CreatedByUserId = ActorUserId,
            };

            await repository.AddAsync(item, ct);
            imported++;
        }

        await db.SaveChangesAsync(ct);
        return new ImportResult(imported, errors.Count, errors);
    }

    public async Task<byte[]> ExportExcelAsync(CancellationToken ct = default)
    {
        var items = await repository.GetAllItemsAsync(ct);
        using var stream = new MemoryStream();
        using var package = SpreadsheetDocument.Create(stream, SpreadsheetDocumentType.Workbook, true);
        var workbookPart = package.AddWorkbookPart();
        workbookPart.Workbook = new Workbook(new Sheets());

        var worksheetPart = workbookPart.AddNewPart<WorksheetPart>();
        worksheetPart.Worksheet = new Worksheet(new SheetData());
        var sheets = workbookPart.Workbook.GetFirstChild<Sheets>()!;
        sheets.Append(new Sheet { Id = workbookPart.GetIdOfPart(worksheetPart), SheetId = 1, Name = "Todo Items" });

        var sheetData = worksheetPart.Worksheet.GetFirstChild<SheetData>()!;
        AppendRow(sheetData, ExcelFieldNames);
        foreach (var item in items)
        {
            AppendRow(sheetData, new[]
            {
                item.Id.ToString(),
                item.Title ?? string.Empty,
                item.Description ?? string.Empty,
                item.IsCompleted.ToString().ToLowerInvariant(),
                item.CreatedAt.ToString("O"),
                item.UpdatedAt?.ToString("O") ?? string.Empty,
            });
        }

        worksheetPart.Worksheet.Save();
        workbookPart.Workbook.Save();
        package.Save();
        stream.Position = 0;
        return stream.ToArray();
    }

    private static IReadOnlyList<string> ReadRowValues(Row row, SharedStringTablePart? sharedStringTablePart)
    {
        var values = new List<string>();
        var cells = row.Elements<Cell>().ToList();
        foreach (var cell in cells)
        {
            values.Add(GetCellValue(cell, sharedStringTablePart));
        }

        return values;
    }

    private static string GetCellValue(Cell cell, SharedStringTablePart? sharedStringTablePart)
    {
        if (cell.DataType?.Value == CellValues.SharedString && sharedStringTablePart is not null)
        {
            if (cell.CellValue is not null && int.TryParse(cell.CellValue.Text, out var index) && index >= 0 && index < sharedStringTablePart.SharedStringTable.Count())
            {
                return sharedStringTablePart.SharedStringTable.ElementAt(index).InnerText;
            }
        }

        if (cell.DataType?.Value == CellValues.InlineString)
        {
            return cell.InlineString?.InnerText ?? string.Empty;
        }

        if (cell.CellValue is not null)
        {
            return cell.CellValue.Text;
        }

        return string.Empty;
    }

    private static void AppendRow(SheetData sheetData, IEnumerable<string> values)
    {
        var row = new Row();
        foreach (var value in values)
        {
            row.AppendChild(new Cell(new CellValue(value)) { DataType = CellValues.String });
        }

        sheetData.AppendChild(row);
    }
}
