using TodoApi.DTOs;

namespace TodoApi.Services;

public interface IUserService
{
    Task<PaginatedResponse<UserResponse>> GetAllAsync(int page, int pageSize, CancellationToken ct = default);
    Task<UserResponse> GetByIdAsync(int id, CancellationToken ct = default);
    Task<UserResponse> CreateAsync(CreateUserRequest request, CancellationToken ct = default);
    Task<UserResponse> UpdateAsync(int id, UpdateUserRequest request, CancellationToken ct = default);
    Task<UserResponse> SetActiveAsync(int id, bool isActive, CancellationToken ct = default);
    Task<UserResponse> SignUpAsync(SignUpRequest request, CancellationToken ct = default);
    Task<UserResponse> GetProfileAsync(int userId, CancellationToken ct = default);
    Task<UserResponse> UpdateProfileAsync(int userId, UpdateProfileRequest request, CancellationToken ct = default);
    Task ChangePasswordAsync(int userId, ChangePasswordRequest request, CancellationToken ct = default);
    Task RequestPasswordResetAsync(ResetPasswordRequest request, CancellationToken ct = default);
    Task ConfirmPasswordResetAsync(ConfirmPasswordResetRequest request, CancellationToken ct = default);
    Task<TokenResponse> CreateTokenAsync(TokenRequest request, CancellationToken ct = default);
}

public sealed class UserConflictException(string message) : Exception(message);
public sealed class InvalidPasswordException(string message) : Exception(message);
public sealed class InvalidPasswordResetTokenException(string message) : Exception(message);
public sealed class InvalidCredentialsException(string message) : Exception(message);
