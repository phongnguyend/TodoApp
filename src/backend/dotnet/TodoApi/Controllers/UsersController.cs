using System.Security.Claims;
using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authorization;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Controllers;

[ApiController]
[Route("api/users")]
[Produces("application/json")]
public class UsersController(IUserService service) : ControllerBase
{
    [HttpGet]
    [ProducesResponseType(typeof(PaginatedResponse<UserResponse>), StatusCodes.Status200OK)]
    public async Task<IActionResult> GetAll([FromQuery] int page = 1, [FromQuery] int pageSize = 20, CancellationToken ct = default) =>
        Ok(await service.GetAllAsync(page, pageSize, ct));

    [HttpGet("{id:int}")]
    [ProducesResponseType(typeof(UserResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status404NotFound)]
    public async Task<IActionResult> GetById(int id, CancellationToken ct = default) =>
        await HandleNotFoundAsync(() => service.GetByIdAsync(id, ct));

    [HttpPost]
    [ProducesResponseType(typeof(UserResponse), StatusCodes.Status201Created)]
    [ProducesResponseType(StatusCodes.Status409Conflict)]
    public async Task<IActionResult> Create([FromBody] CreateUserRequest request, CancellationToken ct = default)
    {
        try
        {
            var user = await service.CreateAsync(request, ct);
            return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
        }
        catch (UserConflictException ex)
        {
            return Conflict(new { error = ex.Message });
        }
    }

    [HttpPut("{id:int}")]
    public Task<IActionResult> Update(int id, [FromBody] UpdateUserRequest request, CancellationToken ct = default) =>
        HandleUserWriteAsync(() => service.UpdateAsync(id, request, ct));

    [HttpPatch("{id:int}/activate")]
    public Task<IActionResult> Activate(int id, CancellationToken ct = default) =>
        HandleUserWriteAsync(() => service.SetActiveAsync(id, true, ct));

    [HttpPatch("{id:int}/deactivate")]
    public Task<IActionResult> Deactivate(int id, CancellationToken ct = default) =>
        HandleUserWriteAsync(() => service.SetActiveAsync(id, false, ct));

    [HttpPost("signup")]
    [ProducesResponseType(typeof(UserResponse), StatusCodes.Status201Created)]
    public async Task<IActionResult> SignUp([FromBody] SignUpRequest request, CancellationToken ct = default)
    {
        try
        {
            var user = await service.SignUpAsync(request, ct);
            return CreatedAtAction(nameof(GetById), new { id = user.Id }, user);
        }
        catch (UserConflictException ex)
        {
            return Conflict(new { error = ex.Message });
        }
    }

    [HttpGet("profile")]
    [Authorize]
    public async Task<IActionResult> GetProfile(CancellationToken ct = default)
    {
        if (!TryGetCurrentUserId(out var userId)) return Unauthorized();
        return await HandleNotFoundAsync(() => service.GetProfileAsync(userId, ct));
    }

    [HttpPut("profile")]
    [Authorize]
    public async Task<IActionResult> UpdateProfile([FromBody] UpdateProfileRequest request, CancellationToken ct = default)
    {
        if (!TryGetCurrentUserId(out var userId)) return Unauthorized();
        return await HandleUserWriteAsync(() => service.UpdateProfileAsync(userId, request, ct));
    }

    [HttpPost("password/change")]
    [Authorize]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    public async Task<IActionResult> ChangePassword([FromBody] ChangePasswordRequest request, CancellationToken ct = default)
    {
        if (!TryGetCurrentUserId(out var userId)) return Unauthorized();
        try
        {
            await service.ChangePasswordAsync(userId, request, ct);
            return NoContent();
        }
        catch (KeyNotFoundException ex)
        {
            return NotFound(new { error = ex.Message });
        }
        catch (InvalidPasswordException ex)
        {
            return BadRequest(new { error = ex.Message });
        }
    }

    [HttpPost("password/reset")]
    [ProducesResponseType(StatusCodes.Status202Accepted)]
    public async Task<IActionResult> RequestPasswordReset([FromBody] ResetPasswordRequest request, CancellationToken ct = default)
    {
        await service.RequestPasswordResetAsync(request, ct);
        return Accepted(new { message = "If the account exists, a password reset email has been queued." });
    }

    [HttpPost("password/confirm")]
    [ProducesResponseType(StatusCodes.Status204NoContent)]
    public async Task<IActionResult> ConfirmPasswordReset([FromBody] ConfirmPasswordResetRequest request, CancellationToken ct = default)
    {
        try
        {
            await service.ConfirmPasswordResetAsync(request, ct);
            return NoContent();
        }
        catch (InvalidPasswordResetTokenException ex)
        {
            return BadRequest(new { error = ex.Message });
        }
    }

    private bool TryGetCurrentUserId(out int userId)
    {
        var principal = HttpContext?.User;
        var value = principal?.FindFirstValue(ClaimTypes.NameIdentifier) ?? principal?.FindFirstValue("sub");
        return int.TryParse(value, out userId) && userId > 0;
    }

    private static async Task<IActionResult> HandleNotFoundAsync(Func<Task<UserResponse>> action)
    {
        try { return new OkObjectResult(await action()); }
        catch (KeyNotFoundException ex) { return new NotFoundObjectResult(new { error = ex.Message }); }
    }

    private static async Task<IActionResult> HandleUserWriteAsync(Func<Task<UserResponse>> action)
    {
        try { return new OkObjectResult(await action()); }
        catch (KeyNotFoundException ex) { return new NotFoundObjectResult(new { error = ex.Message }); }
        catch (UserConflictException ex) { return new ConflictObjectResult(new { error = ex.Message }); }
    }
}
