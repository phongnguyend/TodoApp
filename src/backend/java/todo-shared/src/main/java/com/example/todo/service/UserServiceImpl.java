package com.example.todo.service;

import com.example.todo.dto.*;
import com.example.todo.entity.EmailLog;
import com.example.todo.entity.User;
import com.example.todo.exception.InvalidPasswordException;
import com.example.todo.exception.InvalidPasswordResetTokenException;
import com.example.todo.exception.UserConflictException;
import com.example.todo.exception.UnauthorizedException;
import com.example.todo.repository.EmailLogRepository;
import com.example.todo.repository.UserRepository;
import com.example.todo.security.PasswordHasher;
import com.example.todo.security.UserTokenCodec;
import jakarta.persistence.EntityNotFoundException;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.data.domain.Page;
import org.springframework.data.domain.PageRequest;
import org.springframework.data.domain.Sort;
import org.springframework.stereotype.Service;
import org.springframework.transaction.annotation.Transactional;

import java.net.URLEncoder;
import java.nio.charset.StandardCharsets;
import java.time.Instant;
import java.time.temporal.ChronoUnit;
import java.util.Locale;

@Service
@Transactional(readOnly = true)
public class UserServiceImpl implements UserService {
    private final UserRepository repository;
    private final EmailLogRepository emailLogRepository;
    private final PasswordHasher passwordHasher;
    private final UserTokenCodec tokenCodec;
    private final long resetLifetimeMinutes;
    private final String resetConfirmationUrl;
    private final long accessTokenLifetimeMinutes;
    private final String dummyPasswordHash;

    @Autowired
    public UserServiceImpl(UserRepository repository, EmailLogRepository emailLogRepository,
            PasswordHasher passwordHasher, UserTokenCodec tokenCodec,
            @Value("${app.user.password-reset.token-lifetime-minutes:60}") long resetLifetimeMinutes,
            @Value("${app.user.password-reset.confirmation-url:/reset-password}") String resetConfirmationUrl,
            @Value("${app.user.jwt-token-lifetime-minutes:60}") long accessTokenLifetimeMinutes) {
        this.repository = repository;
        this.emailLogRepository = emailLogRepository;
        this.passwordHasher = passwordHasher;
        this.tokenCodec = tokenCodec;
        this.resetLifetimeMinutes = Math.max(1, resetLifetimeMinutes);
        this.resetConfirmationUrl = resetConfirmationUrl;
        this.accessTokenLifetimeMinutes = Math.max(1, accessTokenLifetimeMinutes);
        this.dummyPasswordHash = passwordHasher.hash("not-a-real-password");
    }

    public UserServiceImpl(UserRepository repository, EmailLogRepository emailLogRepository,
            PasswordHasher passwordHasher, UserTokenCodec tokenCodec,
            long resetLifetimeMinutes, String resetConfirmationUrl) {
        this(repository, emailLogRepository, passwordHasher, tokenCodec,
                resetLifetimeMinutes, resetConfirmationUrl, 60);
    }

    @Override
    public PaginatedResponse<UserResponse> getAll(int page, int pageSize) {
        page = Math.max(1, page);
        pageSize = Math.min(100, Math.max(1, pageSize));
        Page<User> result = repository.findAll(PageRequest.of(page - 1, pageSize, Sort.by("createdAt").descending()));
        return new PaginatedResponse<>(result.getContent().stream().map(UserResponse::from).toList(),
                result.getTotalElements(), page, pageSize, result.getTotalPages());
    }

    @Override
    public UserResponse getById(Long id) { return UserResponse.from(getOrThrow(id)); }

    @Override
    @Transactional
    public UserResponse create(CreateUserRequest request) {
        String username = request.username().trim();
        String email = normalizeEmail(request.email());
        ensureUnique(username, email, null);
        User user = new User(username, email, passwordHasher.hash(request.password()),
                request.isActive() == null || request.isActive());
        return UserResponse.from(repository.save(user));
    }

    @Override
    @Transactional
    public UserResponse update(Long id, UpdateUserRequest request) {
        User user = getOrThrow(id);
        String username = request.username() == null ? user.getUsername() : request.username().trim();
        String email = request.email() == null ? user.getEmail() : normalizeEmail(request.email());
        ensureUnique(username, email, id);
        user.setUsername(username);
        user.setEmail(email);
        if (request.password() != null) user.setPasswordHash(passwordHasher.hash(request.password()));
        return UserResponse.from(repository.save(user));
    }

    @Override
    @Transactional
    public UserResponse setActive(Long id, boolean active) {
        User user = getOrThrow(id);
        user.setActive(active);
        return UserResponse.from(repository.save(user));
    }

    @Override
    @Transactional
    public UserResponse signup(SignUpRequest request) {
        return create(new CreateUserRequest(request.username(), request.email(), request.password(), true));
    }

    @Override
    public UserResponse getProfile(Long userId) { return getById(userId); }

    @Override
    @Transactional
    public UserResponse updateProfile(Long userId, UpdateProfileRequest request) {
        return update(userId, new UpdateUserRequest(request.username(), request.email(), null));
    }

    @Override
    @Transactional
    public void changePassword(Long userId, ChangePasswordRequest request) {
        User user = getOrThrow(userId);
        if (!user.isActive()) throw new InvalidPasswordException("The user account is inactive.");
        if (!passwordHasher.matches(request.currentPassword(), user.getPasswordHash()))
            throw new InvalidPasswordException("The current password is incorrect.");
        user.setPasswordHash(passwordHasher.hash(request.newPassword()));
        repository.save(user);
    }

    @Override
    @Transactional
    public void requestPasswordReset(ResetPasswordRequest request) {
        repository.findByEmailIgnoreCase(normalizeEmail(request.email()))
                .filter(User::isActive)
                .ifPresent(user -> {
                    String token = tokenCodec.createPasswordResetToken(user.getId(),
                            Instant.now().plus(resetLifetimeMinutes, ChronoUnit.MINUTES), user.getPasswordHash());
                    String separator = resetConfirmationUrl.contains("?") ? "&" : "?";
                    String url = resetConfirmationUrl + separator + "token="
                            + URLEncoder.encode(token, StandardCharsets.UTF_8);
                    emailLogRepository.save(new EmailLog(user.getEmail(), "Reset your Todo API password",
                            "Use this link to reset your password: " + url + "\n\nThis link expires in "
                                    + resetLifetimeMinutes + " minutes."));
                });
    }

    @Override
    @Transactional
    public void confirmPasswordReset(ConfirmPasswordResetRequest request) {
        UserTokenCodec.ResetTokenPayload payload = tokenCodec.decodePasswordResetToken(request.token());
        User user = repository.findById(payload.userId())
                .orElseThrow(() -> invalidResetToken());
        if (!user.isActive() || !tokenCodec.passwordMatchesToken(user.getPasswordHash(), payload))
            throw invalidResetToken();
        user.setPasswordHash(passwordHasher.hash(request.newPassword()));
        repository.save(user);
    }

    @Override
    public TokenResponse createToken(TokenRequest request) {
        User user = repository.findByEmailIgnoreCase(normalizeEmail(request.email())).orElse(null);
        boolean passwordValid = passwordHasher.matches(request.password(),
                user == null ? dummyPasswordHash : user.getPasswordHash());
        if (user == null || !passwordValid || !user.isActive())
            throw new UnauthorizedException("Invalid email or password.");
        Instant issuedAt = Instant.now();
        long expiresIn = accessTokenLifetimeMinutes * 60;
        String token = tokenCodec.createAccessToken(user.getId(), issuedAt, issuedAt.plusSeconds(expiresIn));
        return new TokenResponse(token, "Bearer", expiresIn);
    }

    private User getOrThrow(Long id) {
        return repository.findById(id)
                .orElseThrow(() -> new EntityNotFoundException("User " + id + " not found."));
    }

    private void ensureUnique(String username, String email, Long excludingId) {
        boolean usernameExists = excludingId == null ? repository.existsByUsernameIgnoreCase(username)
                : repository.existsByUsernameIgnoreCaseAndIdNot(username, excludingId);
        if (usernameExists) throw new UserConflictException("Username is already in use.");
        boolean emailExists = excludingId == null ? repository.existsByEmailIgnoreCase(email)
                : repository.existsByEmailIgnoreCaseAndIdNot(email, excludingId);
        if (emailExists) throw new UserConflictException("Email is already in use.");
    }

    private static String normalizeEmail(String email) { return email.trim().toLowerCase(Locale.ROOT); }

    private static InvalidPasswordResetTokenException invalidResetToken() {
        return new InvalidPasswordResetTokenException("The password reset token is invalid or expired.");
    }
}
