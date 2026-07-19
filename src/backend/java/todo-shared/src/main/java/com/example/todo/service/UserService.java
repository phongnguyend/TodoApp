package com.example.todo.service;

import com.example.todo.dto.*;

public interface UserService {
    PaginatedResponse<UserResponse> getAll(int page, int pageSize);
    UserResponse getById(Long id);
    UserResponse create(CreateUserRequest request);
    UserResponse update(Long id, UpdateUserRequest request);
    UserResponse setActive(Long id, boolean active);
    UserResponse signup(SignUpRequest request);
    UserResponse getProfile(Long userId);
    UserResponse updateProfile(Long userId, UpdateProfileRequest request);
    void changePassword(Long userId, ChangePasswordRequest request);
    void requestPasswordReset(ResetPasswordRequest request);
    void confirmPasswordReset(ConfirmPasswordResetRequest request);
    TokenResponse createToken(TokenRequest request);
}
