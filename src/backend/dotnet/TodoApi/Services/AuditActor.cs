using System.Security.Claims;

namespace TodoApi.Services;

internal static class AuditActor
{
    public static int? GetUserId(IHttpContextAccessor? accessor)
    {
        var principal = accessor?.HttpContext?.User;
        var value = principal?.FindFirstValue("sub")
            ?? principal?.FindFirstValue(ClaimTypes.NameIdentifier);
        return int.TryParse(value, out var userId) && userId > 0 ? userId : null;
    }
}
