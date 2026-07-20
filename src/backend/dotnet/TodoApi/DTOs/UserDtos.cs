using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace TodoApi.DTOs;

public record CreateUserRequest(
    [Required, StringLength(50, MinimumLength = 1)] string Username,
    [Required, EmailAddress, StringLength(255)] string Email,
    [Required, StringLength(128, MinimumLength = 8)] string Password,
    bool IsActive = true);

public record UpdateUserRequest(
    [StringLength(50, MinimumLength = 1)] string? Username,
    [EmailAddress, StringLength(255)] string? Email,
    [StringLength(128, MinimumLength = 8)] string? Password);

public record SignUpRequest(
    [Required, StringLength(50, MinimumLength = 1)] string Username,
    [Required, EmailAddress, StringLength(255)] string Email,
    [Required, StringLength(128, MinimumLength = 8)] string Password);

public record ChangePasswordRequest(
    [Required] string CurrentPassword,
    [Required, StringLength(128, MinimumLength = 8)] string NewPassword);

public record ResetPasswordRequest([Required, EmailAddress] string Email);

public record ConfirmPasswordResetRequest(
    [Required] string Token,
    [Required, StringLength(128, MinimumLength = 8)] string NewPassword);

public record UpdateProfileRequest(
    [StringLength(50, MinimumLength = 1)] string? Username,
    [EmailAddress, StringLength(255)] string? Email);

public record UserResponse(
    int Id,
    string Username,
    string Email,
    bool IsActive,
    DateTime CreatedAt,
    int? CreatedByUserId,
    DateTime? UpdatedAt,
    int? UpdatedByUserId)
{
    public UserResponse(int id, string username, string email, bool isActive,
        DateTime createdAt, DateTime? updatedAt)
        : this(id, username, email, isActive, createdAt, null, updatedAt, null) { }
}

public record TokenRequest(
    [Required, EmailAddress, StringLength(255)] string Email,
    [Required, StringLength(128, MinimumLength = 1)] string Password);

public record TokenResponse(
    [property: JsonPropertyName("access_token")] string AccessToken,
    [property: JsonPropertyName("token_type")] string TokenType,
    [property: JsonPropertyName("expires_in")] int ExpiresIn);
