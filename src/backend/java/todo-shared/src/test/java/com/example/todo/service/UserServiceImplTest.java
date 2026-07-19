package com.example.todo.service;

import com.example.todo.dto.*;
import com.example.todo.entity.EmailLog;
import com.example.todo.entity.User;
import com.example.todo.exception.InvalidPasswordException;
import com.example.todo.exception.UserConflictException;
import com.example.todo.repository.EmailLogRepository;
import com.example.todo.repository.UserRepository;
import com.example.todo.security.PasswordHasher;
import com.example.todo.security.UserTokenCodec;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.springframework.data.domain.PageImpl;
import org.springframework.data.domain.Pageable;
import org.springframework.security.oauth2.jose.jws.MacAlgorithm;
import org.springframework.security.oauth2.jwt.NimbusJwtDecoder;

import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.util.List;
import java.util.Optional;

import static org.assertj.core.api.Assertions.assertThat;
import static org.assertj.core.api.Assertions.assertThatThrownBy;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

@ExtendWith(MockitoExtension.class)
class UserServiceImplTest {
    @Mock UserRepository repository;
    @Mock EmailLogRepository emailLogRepository;
    private PasswordHasher passwordHasher;
    private UserTokenCodec tokenCodec;
    private UserServiceImpl service;
    private User user;

    @BeforeEach
    void setUp() {
        passwordHasher = new PasswordHasher(1);
        tokenCodec = new UserTokenCodec("test-jwt-secret-at-least-32-bytes-long", "reset-secret");
        service = new UserServiceImpl(repository, emailLogRepository, passwordHasher, tokenCodec,
                60, "https://example.test/reset");
        user = new User("alice", "alice@example.com", passwordHasher.hash("password123"), true);
        user.setId(1L);
        lenient().when(repository.save(any(User.class))).thenAnswer(invocation -> invocation.getArgument(0));
    }

    @Test
    void getAll_clampsPaginationAndMapsUsers() {
        when(repository.findAll(any(Pageable.class))).thenReturn(new PageImpl<>(List.of(user)));

        PaginatedResponse<UserResponse> result = service.getAll(0, 500);

        assertThat(result.page()).isEqualTo(1);
        assertThat(result.pageSize()).isEqualTo(100);
        assertThat(result.items()).extracting(UserResponse::username).containsExactly("alice");
    }

    @Test
    void create_normalizesEmailAndHashesPassword() {
        UserResponse result = service.create(new CreateUserRequest(" alice ", " Alice@Example.com ",
                "password123", null));

        assertThat(result.username()).isEqualTo("alice");
        assertThat(result.email()).isEqualTo("alice@example.com");
        assertThat(userPasswordFromSave()).isNotEqualTo("password123");
        assertThat(passwordHasher.matches("password123", userPasswordFromSave())).isTrue();
    }

    @Test
    void create_duplicateUsernameThrowsConflict() {
        when(repository.existsByUsernameIgnoreCase("alice")).thenReturn(true);

        assertThatThrownBy(() -> service.create(new CreateUserRequest("alice", "a@example.com",
                "password123", true))).isInstanceOf(UserConflictException.class);
    }

    @Test
    void signup_alwaysCreatesActiveUser() {
        UserResponse result = service.signup(new SignUpRequest("bob", "bob@example.com", "password123"));
        assertThat(result.isActive()).isTrue();
    }

    @Test
    void changePassword_rejectsWrongCurrentPassword() {
        when(repository.findById(1L)).thenReturn(Optional.of(user));

        assertThatThrownBy(() -> service.changePassword(1L,
                new ChangePasswordRequest("wrong", "new-password123")))
                .isInstanceOf(InvalidPasswordException.class)
                .hasMessageContaining("incorrect");
    }

    @Test
    void changePassword_replacesHash() {
        when(repository.findById(1L)).thenReturn(Optional.of(user));

        service.changePassword(1L, new ChangePasswordRequest("password123", "new-password123"));

        assertThat(passwordHasher.matches("new-password123", user.getPasswordHash())).isTrue();
    }

    @Test
    void requestPasswordReset_existingActiveUserQueuesEmailWithoutPlainPassword() {
        when(repository.findByEmailIgnoreCase("alice@example.com")).thenReturn(Optional.of(user));
        when(emailLogRepository.save(any())).thenAnswer(invocation -> invocation.getArgument(0));

        service.requestPasswordReset(new ResetPasswordRequest("Alice@Example.com"));

        var captor = org.mockito.ArgumentCaptor.forClass(EmailLog.class);
        verify(emailLogRepository).save(captor.capture());
        assertThat(captor.getValue().getRecipient()).isEqualTo("alice@example.com");
        assertThat(captor.getValue().getBody()).contains("token=").doesNotContain("password123");
    }

    @Test
    void requestPasswordReset_unknownUserDoesNothing() {
        when(repository.findByEmailIgnoreCase("missing@example.com")).thenReturn(Optional.empty());
        service.requestPasswordReset(new ResetPasswordRequest("missing@example.com"));
        verifyNoInteractions(emailLogRepository);
    }

    @Test
    void confirmPasswordReset_acceptsValidTokenAndInvalidatesItAfterUse() {
        when(repository.findById(1L)).thenReturn(Optional.of(user));
        String token = tokenCodec.createPasswordResetToken(1L,
                java.time.Instant.now().plusSeconds(60), user.getPasswordHash());

        service.confirmPasswordReset(new ConfirmPasswordResetRequest(token, "new-password123"));

        assertThat(passwordHasher.matches("new-password123", user.getPasswordHash())).isTrue();
        assertThatThrownBy(() -> service.confirmPasswordReset(
                new ConfirmPasswordResetRequest(token, "another-password")))
                .hasMessageContaining("invalid or expired");
    }

    @Test
    void createToken_issuesJwtValidatedByNimbus() {
        when(repository.findByEmailIgnoreCase("alice@example.com")).thenReturn(Optional.of(user));

        TokenResponse response = service.createToken(new TokenRequest("Alice@Example.com", "password123"));

        var key = new SecretKeySpec("test-jwt-secret-at-least-32-bytes-long"
                .getBytes(StandardCharsets.UTF_8), "HmacSHA256");
        var decoder = NimbusJwtDecoder.withSecretKey(key).macAlgorithm(MacAlgorithm.HS256).build();
        assertThat(decoder.decode(response.accessToken()).getSubject()).isEqualTo("1");
        assertThat(response.tokenType()).isEqualTo("Bearer");
        assertThat(response.expiresIn()).isEqualTo(3600);
    }

    private String userPasswordFromSave() {
        var captor = org.mockito.ArgumentCaptor.forClass(User.class);
        verify(repository, atLeastOnce()).save(captor.capture());
        return captor.getValue().getPasswordHash();
    }
}
