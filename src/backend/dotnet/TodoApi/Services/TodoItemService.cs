using System.Text;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoShared.Models;

namespace TodoApi.Services;

public class TodoItemService(ITodoItemRepository repository, AppDbContext db) : ITodoItemService
{
    private static readonly string[] CsvFieldNames = ["id", "title", "description", "is_completed", "created_at", "updated_at"];
    private static readonly HashSet<string> TrueValues = new(StringComparer.OrdinalIgnoreCase) { "1", "true", "yes", "y" };

    // ── Mapping ───────────────────────────────────────────────────────────────

    private static TodoItemResponse ToResponse(TodoItem item) =>
        new(item.Id, item.Title, item.Description, item.IsCompleted, item.CreatedAt, item.UpdatedAt);

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
}
