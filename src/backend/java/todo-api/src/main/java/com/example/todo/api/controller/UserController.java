package com.example.todo.api.controller;

import com.example.todo.api.config.AuditActor;
import com.example.todo.dto.*;
import com.example.todo.service.UserService;
import io.swagger.v3.oas.annotations.Operation;
import io.swagger.v3.oas.annotations.security.SecurityRequirement;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.HttpStatus;
import org.springframework.web.bind.annotation.*;
import org.springframework.security.core.annotation.AuthenticationPrincipal;
import org.springframework.security.oauth2.jwt.Jwt;

@RestController
@RequestMapping("/api/users")
@RequiredArgsConstructor
@Tag(name = "Users")
public class UserController {
    private final UserService service;

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
    public UserResponse create(@Valid @RequestBody CreateUserRequest request) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.create(request) : service.create(request, actor);
    }

    @PutMapping("/{id}")
    @Operation(summary = "Update a user")
    public UserResponse update(@PathVariable Long id, @Valid @RequestBody UpdateUserRequest request) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.update(id, request) : service.update(id, request, actor);
    }

    @PatchMapping("/{id}/activate")
    @Operation(summary = "Activate a user")
    public UserResponse activate(@PathVariable Long id) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.setActive(id, true) : service.setActive(id, true, actor);
    }

    @PatchMapping("/{id}/deactivate")
    @Operation(summary = "Deactivate a user")
    public UserResponse deactivate(@PathVariable Long id) {
        Long actor = AuditActor.currentUserId();
        return actor == null ? service.setActive(id, false) : service.setActive(id, false, actor);
    }

    @PostMapping("/signup")
    @ResponseStatus(HttpStatus.CREATED)
    @Operation(summary = "Register a new account")
    public UserResponse signup(@Valid @RequestBody SignUpRequest request) { return service.signup(request); }

    @GetMapping("/profile")
    @SecurityRequirement(name = "bearerAuth")
    @Operation(summary = "Read the authenticated user's profile")
    public UserResponse profile(@AuthenticationPrincipal Jwt jwt) {
        return service.getProfile(Long.parseLong(jwt.getSubject()));
    }

    @PutMapping("/profile")
    @SecurityRequirement(name = "bearerAuth")
    @Operation(summary = "Update the authenticated user's profile")
    public UserResponse updateProfile(
            @AuthenticationPrincipal Jwt jwt,
            @Valid @RequestBody UpdateProfileRequest request) {
        return service.updateProfile(Long.parseLong(jwt.getSubject()), request);
    }

    @PostMapping("/password/change")
    @ResponseStatus(HttpStatus.NO_CONTENT)
    @SecurityRequirement(name = "bearerAuth")
    @Operation(summary = "Change the authenticated user's password")
    public void changePassword(
            @AuthenticationPrincipal Jwt jwt,
            @Valid @RequestBody ChangePasswordRequest request) {
        service.changePassword(Long.parseLong(jwt.getSubject()), request);
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
