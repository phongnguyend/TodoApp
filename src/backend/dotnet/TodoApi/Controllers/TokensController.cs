using Microsoft.AspNetCore.Mvc;
using Microsoft.AspNetCore.Authorization;
using TodoApi.DTOs;
using TodoApi.Services;

namespace TodoApi.Controllers;

[ApiController]
[Route("api/tokens")]
[Produces("application/json")]
[AllowAnonymous]
public class TokensController(IUserService service) : ControllerBase
{
    [HttpPost]
    [ProducesResponseType(typeof(TokenResponse), StatusCodes.Status200OK)]
    [ProducesResponseType(StatusCodes.Status400BadRequest)]
    [ProducesResponseType(StatusCodes.Status401Unauthorized)]
    public async Task<IActionResult> Create([FromBody] TokenRequest request, CancellationToken ct = default)
    {
        try
        {
            var token = await service.CreateTokenAsync(request, ct);
            Response.Headers.CacheControl = "no-store";
            Response.Headers.Pragma = "no-cache";
            return Ok(token);
        }
        catch (InvalidCredentialsException)
        {
            Response.Headers.WWWAuthenticate = "Bearer";
            return Unauthorized(new { error = "Invalid email or password." });
        }
    }
}
