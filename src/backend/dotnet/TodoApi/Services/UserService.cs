using System.Text.Json;
using Microsoft.AspNetCore.DataProtection;
using Microsoft.AspNetCore.Identity;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoShared.Models;

namespace TodoApi.Services;

public class UserService(
    IUserRepository repository,
    AppDbContext db,
    IPasswordHasher<User> passwordHasher,
    IDataProtectionProvider dataProtectionProvider,
    IConfiguration configuration) : IUserService
{
    private readonly IDataProtector _resetTokenProtector =
        dataProtectionProvider.CreateProtector("TodoApi.UserPasswordReset.v1");

    private static UserResponse ToResponse(User user) =>
        new(user.Id, user.Username, user.Email, user.IsActive, user.CreatedAt, user.UpdatedAt);

    private async Task<User> GetOrThrowAsync(int id, CancellationToken ct) =>
        await repository.GetByIdAsync(id, ct) ?? throw new KeyNotFoundException($"User {id} not found.");

    public async Task<PaginatedResponse<UserResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default)
    {
        page = Math.Max(1, page);
        pageSize = Math.Clamp(pageSize, 1, 100);
        var (users, total) = await repository.GetAllAsync((page - 1) * pageSize, pageSize, ct);
        return new PaginatedResponse<UserResponse>(
            users.Select(ToResponse), total, page, pageSize, (int)Math.Ceiling(total / (double)pageSize));
    }

    public async Task<UserResponse> GetByIdAsync(int id, CancellationToken ct = default) =>
        ToResponse(await GetOrThrowAsync(id, ct));

    public async Task<UserResponse> CreateAsync(CreateUserRequest request, CancellationToken ct = default)
    {
        await EnsureUniqueAsync(request.Username, request.Email, null, ct);
        var user = new User
        {
            Username = request.Username.Trim(),
            Email = NormalizeEmail(request.Email),
            IsActive = request.IsActive,
            CreatedAt = DateTime.UtcNow
        };
        user.PasswordHash = passwordHasher.HashPassword(user, request.Password);
        await repository.AddAsync(user, ct);
        await db.SaveChangesAsync(ct);
        return ToResponse(user);
    }

    public async Task<UserResponse> UpdateAsync(int id, UpdateUserRequest request, CancellationToken ct = default)
    {
        var user = await GetOrThrowAsync(id, ct);
        var username = request.Username?.Trim() ?? user.Username;
        var email = request.Email is null ? user.Email : NormalizeEmail(request.Email);
        await EnsureUniqueAsync(username, email, id, ct);

        user.Username = username;
        user.Email = email;
        if (request.Password is not null)
            user.PasswordHash = passwordHasher.HashPassword(user, request.Password);
        user.UpdatedAt = DateTime.UtcNow;
        repository.Update(user);
        await db.SaveChangesAsync(ct);
        return ToResponse(user);
    }

    public async Task<UserResponse> SetActiveAsync(int id, bool isActive, CancellationToken ct = default)
    {
        var user = await GetOrThrowAsync(id, ct);
        user.IsActive = isActive;
        user.UpdatedAt = DateTime.UtcNow;
        repository.Update(user);
        await db.SaveChangesAsync(ct);
        return ToResponse(user);
    }

    public Task<UserResponse> SignUpAsync(SignUpRequest request, CancellationToken ct = default) =>
        CreateAsync(new CreateUserRequest(request.Username, request.Email, request.Password), ct);

    public Task<UserResponse> GetProfileAsync(int userId, CancellationToken ct = default) => GetByIdAsync(userId, ct);

    public async Task<UserResponse> UpdateProfileAsync(int userId, UpdateProfileRequest request, CancellationToken ct = default)
    {
        var user = await GetOrThrowAsync(userId, ct);
        var username = request.Username?.Trim() ?? user.Username;
        var email = request.Email is null ? user.Email : NormalizeEmail(request.Email);
        await EnsureUniqueAsync(username, email, userId, ct);
        user.Username = username;
        user.Email = email;
        user.UpdatedAt = DateTime.UtcNow;
        repository.Update(user);
        await db.SaveChangesAsync(ct);
        return ToResponse(user);
    }

    public async Task ChangePasswordAsync(int userId, ChangePasswordRequest request, CancellationToken ct = default)
    {
        var user = await GetOrThrowAsync(userId, ct);
        if (!user.IsActive)
            throw new InvalidPasswordException("The user account is inactive.");

        var verification = passwordHasher.VerifyHashedPassword(user, user.PasswordHash, request.CurrentPassword);
        if (verification == PasswordVerificationResult.Failed)
            throw new InvalidPasswordException("The current password is incorrect.");

        user.PasswordHash = passwordHasher.HashPassword(user, request.NewPassword);
        user.UpdatedAt = DateTime.UtcNow;
        repository.Update(user);
        await db.SaveChangesAsync(ct);
    }

    public async Task RequestPasswordResetAsync(ResetPasswordRequest request, CancellationToken ct = default)
    {
        var user = await repository.GetByEmailAsync(NormalizeEmail(request.Email), ct);
        if (user is null || !user.IsActive)
            return; // Do not reveal whether an account exists.

        var lifetimeMinutes = Math.Max(1, configuration.GetValue("PasswordReset:TokenLifetimeMinutes", 60));
        var payload = new PasswordResetTokenPayload(user.Id, DateTime.UtcNow.AddMinutes(lifetimeMinutes), user.PasswordHash);
        var token = _resetTokenProtector.Protect(JsonSerializer.Serialize(payload));
        var baseUrl = configuration["PasswordReset:ConfirmationUrl"] ?? "/reset-password";
        var separator = baseUrl.Contains('?') ? '&' : '?';
        var resetUrl = $"{baseUrl}{separator}token={Uri.EscapeDataString(token)}";

        db.EmailLogs.Add(new EmailLog
        {
            Recipient = user.Email,
            Subject = "Reset your Todo API password",
            Body = $"Use this link to reset your password: {resetUrl}\n\nThis link expires in {lifetimeMinutes} minutes.",
            Status = "pending",
            CreatedAt = DateTime.UtcNow
        });
        await db.SaveChangesAsync(ct);
    }

    public async Task ConfirmPasswordResetAsync(ConfirmPasswordResetRequest request, CancellationToken ct = default)
    {
        PasswordResetTokenPayload payload;
        try
        {
            payload = JsonSerializer.Deserialize<PasswordResetTokenPayload>(_resetTokenProtector.Unprotect(request.Token))
                ?? throw new InvalidPasswordResetTokenException("The password reset token is invalid or expired.");
        }
        catch (InvalidPasswordResetTokenException)
        {
            throw;
        }
        catch (Exception ex) when (ex is System.Security.Cryptography.CryptographicException or JsonException)
        {
            throw new InvalidPasswordResetTokenException("The password reset token is invalid or expired.");
        }

        var user = await repository.GetByIdAsync(payload.UserId, ct);
        if (user is null || !user.IsActive || payload.ExpiresAtUtc < DateTime.UtcNow ||
            !string.Equals(payload.PasswordHash, user.PasswordHash, StringComparison.Ordinal))
            throw new InvalidPasswordResetTokenException("The password reset token is invalid or expired.");

        user.PasswordHash = passwordHasher.HashPassword(user, request.NewPassword);
        user.UpdatedAt = DateTime.UtcNow;
        repository.Update(user);
        await db.SaveChangesAsync(ct);
    }

    private async Task EnsureUniqueAsync(string username, string email, int? excludingId, CancellationToken ct)
    {
        if (await repository.UsernameExistsAsync(username.Trim(), excludingId, ct))
            throw new UserConflictException("Username is already in use.");
        if (await repository.EmailExistsAsync(NormalizeEmail(email), excludingId, ct))
            throw new UserConflictException("Email is already in use.");
    }

    private static string NormalizeEmail(string email) => email.Trim().ToLowerInvariant();

    private sealed record PasswordResetTokenPayload(int UserId, DateTime ExpiresAtUtc, string PasswordHash);
}
