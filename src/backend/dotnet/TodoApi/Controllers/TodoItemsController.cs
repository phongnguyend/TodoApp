using System.Text;
using Microsoft.AspNetCore.Mvc;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Controllers;

[ApiController]
[Route("api/todo-items")]
[Produces("application/json")]
public class TodoItemsController(ITodoItemService service, ITodoItemAttachmentService attachmentService) : ControllerBase
{
    // GET api/todo-items?page=1&pageSize=20
    [HttpGet]
    [ProducesResponseType(typeof(PaginatedResponse<TodoItemResponse>), StatusCodes.Status200OK)]
    public async Task<IActionResult> GetAll(
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20,
        CancellationToken ct = default)
    {
        var result = await service.GetAllAsync(page, pageSize, ct);
        return Ok(result);
    }

    // GET api/todo-items/incomplete
    [HttpGet("incomplete")]
    [ProducesResponseType(typeof(PaginatedResponse<TodoItemResponse>), StatusCodes.Status200OK)]
    public async Task<IActionResult> GetIncomplete(
        [FromQuery] int page = 1,
        [FromQuery] int pageSize = 20,
        CancellationToken ct = default)
    {
        var result = await service.GetIncompleteAsync(page, pageSize, ct);
        return Ok(result);
    }

    // GET api/todo-items/{id}
    [HttpGet("{id:int}")]
    [ProducesResponseType(typeof(TodoItemResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> GetById(int id, CancellationToken ct = default)
    {
        try
        {
            var result = await service.GetByIdAsync(id, ct);
            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // POST api/todo-items
    [HttpPost]
    [ProducesResponseType(typeof(TodoItemResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> Create([FromBody] CreateTodoItemRequest request, CancellationToken ct = default)
    {
        var result = await service.CreateAsync(request, ct);
        return CreatedAtAction(nameof(GetById), new { id = result.Id }, result);
    }

    // PUT api/todo-items/{id}
    [HttpPut("{id:int}")]
    [ProducesResponseType(typeof(TodoItemResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> Update(int id, [FromBody] UpdateTodoItemRequest request, CancellationToken ct = default)
    {
        try
        {
            var result = await service.UpdateAsync(id, request, ct);
            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // PATCH api/todo-items/{id}/complete
    [HttpPatch("{id:int}/complete")]
    [ProducesResponseType(typeof(TodoItemResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> MarkComplete(int id, CancellationToken ct = default)
    {
        try
        {
            var result = await service.MarkCompleteAsync(id, ct);
            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // GET api/todo-items/{id}/attachments
    [HttpGet("{id:int}/attachments")]
    [ProducesResponseType(typeof(IReadOnlyList<TodoItemAttachmentResponse>), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> GetAttachments(int id, CancellationToken ct = default)
    {
        try
        {
            var result = await attachmentService.GetAllAsync(id, ct);
            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // POST api/todo-items/{id}/attachments
    [HttpPost("{id:int}/attachments")]
    [ProducesResponseType(typeof(TodoItemAttachmentResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> CreateAttachment(int id, [FromBody] CreateTodoItemAttachmentRequest request, CancellationToken ct = default)
    {
        try
        {
            var result = await attachmentService.CreateAsync(id, request, ct);
            return CreatedAtAction(nameof(GetAttachmentById), new { id, attachmentId = result.Id }, result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // GET api/todo-items/{id}/attachments/{attachmentId}
    [HttpGet("{id:int}/attachments/{attachmentId:int}")]
    [ProducesResponseType(typeof(TodoItemAttachmentResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> GetAttachmentById(int id, int attachmentId, CancellationToken ct = default)
    {
        try
        {
            var result = await attachmentService.GetByIdAsync(id, attachmentId, ct);
            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // PUT api/todo-items/{id}/attachments/{attachmentId}
    [HttpPut("{id:int}/attachments/{attachmentId:int}")]
    [ProducesResponseType(typeof(TodoItemAttachmentResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> UpdateAttachment(int id, int attachmentId, [FromBody] CreateTodoItemAttachmentRequest request, CancellationToken ct = default)
    {
        try
        {
            var result = await attachmentService.UpdateAsync(id, attachmentId, request, ct);
            return Ok(result);
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // DELETE api/todo-items/{id}/attachments/{attachmentId}
    [HttpDelete("{id:int}/attachments/{attachmentId:int}")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> DeleteAttachment(int id, int attachmentId, CancellationToken ct = default)
    {
        try
        {
            await attachmentService.DeleteAsync(id, attachmentId, ct);
            return NoContent();
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }

    // POST api/todo-items/import/csv
    [HttpPost("import/csv")]
    [Consumes("multipart/form-data")]
    [ProducesResponseType(typeof(ImportResult), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> ImportCsv(IFormFile? file, CancellationToken ct = default)
    {
        if (file is null || file.Length == 0)
        {
            return BadRequest(new { error = "file is required" });
        }

        var result = await service.ImportCsvAsync(file, ct);
        return Ok(result);
    }

    // GET api/todo-items/export/csv
    [HttpGet("export/csv")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public async Task<IActionResult> ExportCsv(CancellationToken ct = default)
    {
        var csv = await service.ExportCsvAsync(ct);
        var bytes = Encoding.UTF8.GetBytes(csv);
        return File(bytes, "text/csv; charset=utf-8", "todo-items.csv");
    }

    // POST api/todo-items/import/excel
    [HttpPost("import/excel")]
    [Consumes("multipart/form-data")]
    [ProducesResponseType(typeof(ImportResult), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    public async Task<IActionResult> ImportExcel(IFormFile? file, CancellationToken ct = default)
    {
        if (file is null || file.Length == 0)
        {
            return BadRequest(new { error = "file is required" });
        }

        var result = await service.ImportExcelAsync(file, ct);
        return Ok(result);
    }

    // GET api/todo-items/export/excel
    [HttpGet("export/excel")]
    [ProducesResponseType(StatusCodes.Status200OK)]
    public async Task<IActionResult> ExportExcel(CancellationToken ct = default)
    {
        var bytes = await service.ExportExcelAsync(ct);
        return File(bytes, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", "todo-items.xlsx");
    }

    // DELETE api/todo-items/{id}
    [HttpDelete("{id:int}")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> Delete(int id, CancellationToken ct = default)
    {
        try
        {
            await service.DeleteAsync(id, ct);
            return NoContent();
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
    }
}
