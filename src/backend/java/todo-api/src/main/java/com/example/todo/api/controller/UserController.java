package com.example.todo.api.controller;

import com.example.todo.dto.*;
import com.example.todo.security.UserTokenCodec;
import com.example.todo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;

@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
@Tag(name = "Users")
public class UserController {
    private final UserService service;
    private final UserTokenCodec tokenCodec;

    @GetMapping
    @Operation(summary = "List users (paginated)")
    public PaginatedResponse<UserResponse> getAll(
            @RequestParam(defaultValue = "1") int page,
            @RequestParam(defaultValue = "20") int pageSize) {
        return service.getAll(page, pageSize);
    }

    @GetMapping("/{id}")
    @Operation(summary = "Get a user by ID")
    public UserResponse getById(@PathVariable Long id) { return service.getById(id); }

    @PostMapping
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Create a user")
    public UserResponse create(@Valid @RequestBody CreateUserRequest request) { return service.create(request); }

    @PutMapping("/{id}")
    @Operation(summary = "Update a user")
    public UserResponse update(@PathVariable Long id, @Valid @RequestBody UpdateUserRequest request) {
        return service.update(id, request);
    }

    @PatchMapping("/{id}/activate")
    @Operation(summary = "Activate a user")
    public UserResponse activate(@PathVariable Long id) { return service.setActive(id, true); }

    @PatchMapping("/{id}/deactivate")
    @Operation(summary = "Deactivate a user")
    public UserResponse deactivate(@PathVariable Long id) { return service.setActive(id, false); }

    @PostMapping("/signup")
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Register a new account")
    public UserResponse signup(@Valid @RequestBody SignUpRequest request) { return service.signup(request); }

    @GetMapping("/profile")
    @SecurityRequirement(name = "bearerAuth")
    @Operation(summary = "Read the authenticated user's profile")
    public UserResponse profile(@RequestHeader(value = "Authorization", required = false) String authorization) {
        return service.getProfile(tokenCodec.authenticatedUserId(authorization));
    }

    @PutMapping("/profile")
    @SecurityRequirement(name = "bearerAuth")
    @Operation(summary = "Update the authenticated user's profile")
    public UserResponse updateProfile(
            @RequestHeader(value = "Authorization", required = false) String authorization,
            @Valid @RequestBody UpdateProfileRequest request) {
        return service.updateProfile(tokenCodec.authenticatedUserId(authorization), request);
    }

    @PostMapping("/password/change")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @SecurityRequirement(name = "bearerAuth")
    @Operation(summary = "Change the authenticated user's password")
    public void changePassword(
            @RequestHeader(value = "Authorization", required = false) String authorization,
            @Valid @RequestBody ChangePasswordRequest request) {
        service.changePassword(tokenCodec.authenticatedUserId(authorization), request);
    }

    @PostMapping("/password/reset")
    @ResponseStatus(HttpStatus.ACCEPTED)
    @Operation(summary = "Request a password reset email")
    public MessageResponse requestPasswordReset(@Valid @RequestBody ResetPasswordRequest request) {
        service.requestPasswordReset(request);
        return new MessageResponse("If the account exists, a password reset email has been queued.");
    }

    @PostMapping("/password/confirm")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @Operation(summary = "Confirm a password reset")
    public void confirmPasswordReset(@Valid @RequestBody ConfirmPasswordResetRequest request) {
        service.confirmPasswordReset(request);
    }
}
