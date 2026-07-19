using Microsoft.AspNetCore.DataProtection;
using Microsoft.AspNetCore.Identity;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Configuration;
using Microsoft.IdentityModel.Tokens;
using System.IdentityModel.Tokens.Jwt;
using System.Text;
using TodoApi.Data;
using TodoApi.DTOs;
using TodoApi.Repositories;
using TodoApi.Services;
using TodoShared.Models;

namespace TodoApi.Tests.Services;

public class UserServiceTests
{
    private static (UserService Service, AppDbContext Db, PasswordHasher<User> Hasher) CreateSut()
    {
        var options = new DbContextOptionsBuilder<AppDbContext>()
            .UseInMemoryDatabase(Guid.NewGuid().ToString())
            .Options;
        var db = new AppDbContext(options);
        var repository = new UserRepository(db);
        var hasher = new PasswordHasher<User>();
        var keyDirectory = new DirectoryInfo(Path.Combine(Path.GetTempPath(), "todo-api-tests", Guid.NewGuid().ToString()));
        var protectionProvider = DataProtectionProvider.Create(keyDirectory);
        var configuration = new ConfigurationBuilder().AddInMemoryCollection(new Dictionary<string, string?>
        {
            ["PasswordReset:TokenLifetimeMinutes"] = "60",
            ["PasswordReset:ConfirmationUrl"] = "https://example.test/reset",
            ["Authentication:Secret"] = "test-secret-at-least-32-bytes-long",
            ["Authentication:TokenLifetimeMinutes"] = "60"
        }).Build();
        return (new UserService(repository, db, hasher, protectionProvider, configuration), db, hasher);
    }

    [Fact]
    public async Task CreateAsync_HashesPasswordAndDoesNotExposeIt()
    {
        var (service, db, hasher) = CreateSut();

        var response = await service.CreateAsync(new CreateUserRequest("alice", "Alice@Example.com", "password123"));
        var stored = await db.Users.SingleAsync();

        Assert.Equal("alice@example.com", response.Email);
        Assert.NotEqual("password123", stored.PasswordHash);
        Assert.Equal(PasswordVerificationResult.Success,
            hasher.VerifyHashedPassword(stored, stored.PasswordHash, "password123"));
    }

    [Fact]
    public async Task CreateAsync_RejectsDuplicateUsernameCaseInsensitively()
    {
        var (service, _, _) = CreateSut();
        await service.CreateAsync(new CreateUserRequest("Alice", "alice@example.com", "password123"));

        await Assert.ThrowsAsync<UserConflictException>(() =>
            service.CreateAsync(new CreateUserRequest("ALICE", "other@example.com", "password123")));
    }

    [Fact]
    public async Task PasswordReset_QueuesEmailAndTokenCanOnlyBeUsedOnce()
    {
        var (service, db, hasher) = CreateSut();
        await service.CreateAsync(new CreateUserRequest("alice", "alice@example.com", "password123"));

        await service.RequestPasswordResetAsync(new ResetPasswordRequest("alice@example.com"));
        var email = await db.EmailLogs.SingleAsync();
        var tokenText = email.Body.Split("token=", 2)[1].Split('\n', 2)[0];
        var token = Uri.UnescapeDataString(tokenText);

        await service.ConfirmPasswordResetAsync(new ConfirmPasswordResetRequest(token, "new-password123"));
        var user = await db.Users.SingleAsync();
        Assert.Equal(PasswordVerificationResult.Success,
            hasher.VerifyHashedPassword(user, user.PasswordHash, "new-password123"));
        await Assert.ThrowsAsync<InvalidPasswordResetTokenException>(() =>
            service.ConfirmPasswordResetAsync(new ConfirmPasswordResetRequest(token, "another-password123")));
    }

    [Fact]
    public async Task RequestPasswordReset_DoesNotRevealOrQueueUnknownEmail()
    {
        var (service, db, _) = CreateSut();

        await service.RequestPasswordResetAsync(new ResetPasswordRequest("missing@example.com"));

        Assert.Empty(db.EmailLogs);
    }

    [Fact]
    public async Task CreateTokenAsync_IssuesJwtValidatedByIdentityModel()
    {
        var (service, _, _) = CreateSut();
        await service.CreateAsync(new CreateUserRequest("alice", "alice@example.com", "password123"));

        var response = await service.CreateTokenAsync(new TokenRequest("Alice@Example.com", "password123"));

        var handler = new JwtSecurityTokenHandler { MapInboundClaims = false };
        var principal = handler.ValidateToken(response.AccessToken, new TokenValidationParameters
        {
            ValidateIssuer = false,
            ValidateAudience = false,
            ValidateIssuerSigningKey = true,
            IssuerSigningKey = new SymmetricSecurityKey(
                Encoding.UTF8.GetBytes("test-secret-at-least-32-bytes-long")),
            ValidateLifetime = true,
            ClockSkew = TimeSpan.Zero
        }, out _);
        Assert.Equal("1", principal.FindFirst(JwtRegisteredClaimNames.Sub)?.Value);
        Assert.Equal("Bearer", response.TokenType);
        Assert.Equal(3600, response.ExpiresIn);
    }

    [Fact]
    public async Task CreateTokenAsync_RejectsInvalidCredentialsAndInactiveUsersIdentically()
    {
        var (service, _, _) = CreateSut();
        var user = await service.CreateAsync(new CreateUserRequest("alice", "alice@example.com", "password123"));

        await Assert.ThrowsAsync<InvalidCredentialsException>(() =>
            service.CreateTokenAsync(new TokenRequest("alice@example.com", "wrong")));
        await service.SetActiveAsync(user.Id, false);
        await Assert.ThrowsAsync<InvalidCredentialsException>(() =>
            service.CreateTokenAsync(new TokenRequest("alice@example.com", "password123")));
    }
}
