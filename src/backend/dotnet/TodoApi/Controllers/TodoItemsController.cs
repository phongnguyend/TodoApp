using Microsoft.AspNetCore.Mvc;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Controllers;

[ApiController]
[Route("api/todo-items")]
[Produces("application/json")]
public class TodoItemsController(ITodoItemService service) : ControllerBase
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
