using System.Security.Claims;
using System.Security.Cryptography;
using System.Text;
using System.Text.Json;

namespace TodoApi.Security;

public sealed class JwtAuthenticationMiddleware(RequestDelegate next, IConfiguration configuration)
{
    public async Task InvokeAsync(HttpContext context)
    {
        var authorization = context.Request.Headers.Authorization.ToString();
        if (authorization.StartsWith("Bearer ", StringComparison.OrdinalIgnoreCase))
        {
            var token = authorization[7..].Trim();
            if (TryValidate(token, out var userId))
            {
                context.User = new ClaimsPrincipal(new ClaimsIdentity([
                    new Claim(ClaimTypes.NameIdentifier, userId),
                    new Claim("sub", userId)
                ], "Bearer"));
            }
        }
        await next(context);
    }

    private bool TryValidate(string token, out string userId)
    {
        userId = string.Empty;
        try
        {
            var parts = token.Split('.');
            if (parts.Length != 3) return false;
            using var header = JsonDocument.Parse(Decode(parts[0]));
            if (header.RootElement.GetProperty("alg").GetString() != "HS256") return false;
            var secret = configuration["JWT_SECRET_KEY"] ?? configuration["Authentication:Secret"] ?? "change-me";
            using var hmac = new HMACSHA256(Encoding.UTF8.GetBytes(secret));
            var expected = hmac.ComputeHash(Encoding.UTF8.GetBytes($"{parts[0]}.{parts[1]}"));
            var supplied = Decode(parts[2]);
            if (!CryptographicOperations.FixedTimeEquals(expected, supplied)) return false;
            using var payload = JsonDocument.Parse(Decode(parts[1]));
            userId = payload.RootElement.GetProperty("sub").GetString() ?? string.Empty;
            var exp = payload.RootElement.GetProperty("exp").GetInt64();
            return int.TryParse(userId, out var id) && id > 0 && exp >= DateTimeOffset.UtcNow.ToUnixTimeSeconds();
        }
        catch (Exception ex) when (ex is FormatException or JsonException or CryptographicException or KeyNotFoundException or InvalidOperationException)
        {
            return false;
        }
    }

    private static byte[] Decode(string value)
    {
        value = value.Replace('-', '+').Replace('_', '/');
        value += new string('=', (4 - value.Length % 4) % 4);
        return Convert.FromBase64String(value);
    }
}
