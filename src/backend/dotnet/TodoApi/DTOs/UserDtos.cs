using System.ComponentModel.DataAnnotations;

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
    DateTime? UpdatedAt);
