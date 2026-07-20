package com.example.todo.service;

import com.example.todo.dto.*;

public interface UserService {
    PaginatedResponse<UserResponse> getAll(int page, int pageSize);
    UserResponse getById(Long id);
    UserResponse create(CreateUserRequest request);
    default UserResponse create(CreateUserRequest request, Long actorUserId) { return create(request); }
    UserResponse update(Long id, UpdateUserRequest request);
    default UserResponse update(Long id, UpdateUserRequest request, Long actorUserId) { return update(id, request); }
    UserResponse setActive(Long id, boolean active);
    default UserResponse setActive(Long id, boolean active, Long actorUserId) { return setActive(id, active); }
    UserResponse signup(SignUpRequest request);
    UserResponse getProfile(Long userId);
    UserResponse updateProfile(Long userId, UpdateProfileRequest request);
    void changePassword(Long userId, ChangePasswordRequest request);
    void requestPasswordReset(ResetPasswordRequest request);
    void confirmPasswordReset(ConfirmPasswordResetRequest request);
    TokenResponse createToken(TokenRequest request);
}
