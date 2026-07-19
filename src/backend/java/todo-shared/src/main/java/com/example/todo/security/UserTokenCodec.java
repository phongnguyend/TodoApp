package com.example.todo.security;

import com.example.todo.exception.InvalidPasswordResetTokenException;
import com.example.todo.exception.UnauthorizedException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import java.nio.charset.StandardCharsets;
import java.security.MessageDigest;
import java.time.Instant;
import java.util.Base64;
import java.util.Map;

@Component
public class UserTokenCodec {
    private static final ObjectMapper JSON = new ObjectMapper();
    private static final Base64.Encoder ENCODER = Base64.getUrlEncoder().withoutPadding();
    private static final Base64.Decoder DECODER = Base64.getUrlDecoder();
    private final String jwtSecret;
    private final String resetSecret;

    public UserTokenCodec(
            @Value("${app.user.jwt-secret:change-me}") String jwtSecret,
            @Value("${app.user.password-reset.secret:${app.user.jwt-secret:change-me}}") String resetSecret) {
        this.jwtSecret = jwtSecret;
        this.resetSecret = resetSecret;
    }

    public long authenticatedUserId(String authorizationHeader) {
        try {
            if (authorizationHeader == null || !authorizationHeader.regionMatches(true, 0, "Bearer ", 0, 7))
                throw new IllegalArgumentException();
            String[] parts = authorizationHeader.substring(7).trim().split("\\.", -1);
            if (parts.length != 3 || !validSignature(parts[0] + "." + parts[1], parts[2], jwtSecret))
                throw new IllegalArgumentException();
            JsonNode payload = JSON.readTree(DECODER.decode(parts[1]));
            long id = Long.parseLong(payload.path("sub").asText());
            if (id <= 0 || (payload.has("exp") && payload.path("exp").asLong() < Instant.now().getEpochSecond()))
                throw new IllegalArgumentException();
            return id;
        } catch (Exception ex) {
            throw new UnauthorizedException("A valid bearer token is required.");
        }
    }

    public String createPasswordResetToken(long userId, Instant expiresAt, String passwordHash) {
        try {
            String payload = ENCODER.encodeToString(JSON.writeValueAsBytes(Map.of(
                    "sub", userId,
                    "exp", expiresAt.getEpochSecond(),
                    "password", fingerprint(passwordHash))));
            return payload + "." + sign(payload, resetSecret);
        } catch (Exception ex) {
            throw new IllegalStateException("Could not create a password reset token.", ex);
        }
    }

    public String createAccessToken(long userId, Instant issuedAt, Instant expiresAt) {
        try {
            String header = ENCODER.encodeToString(JSON.writeValueAsBytes(Map.of("alg", "HS256", "typ", "JWT")));
            String payload = ENCODER.encodeToString(JSON.writeValueAsBytes(Map.of(
                    "sub", Long.toString(userId),
                    "iat", issuedAt.getEpochSecond(),
                    "exp", expiresAt.getEpochSecond())));
            String content = header + "." + payload;
            return content + "." + sign(content, jwtSecret);
        } catch (Exception ex) {
            throw new IllegalStateException("Could not create an access token.", ex);
        }
    }

    public ResetTokenPayload decodePasswordResetToken(String token) {
        try {
            String[] parts = token.split("\\.", -1);
            if (parts.length != 2 || !validSignature(parts[0], parts[1], resetSecret))
                throw new IllegalArgumentException();
            JsonNode payload = JSON.readTree(DECODER.decode(parts[0]));
            long userId = payload.path("sub").asLong();
            long expiresAt = payload.path("exp").asLong();
            String passwordFingerprint = payload.path("password").asText();
            if (userId <= 0 || expiresAt < Instant.now().getEpochSecond() || passwordFingerprint.isBlank())
                throw new IllegalArgumentException();
            return new ResetTokenPayload(userId, passwordFingerprint);
        } catch (Exception ex) {
            throw new InvalidPasswordResetTokenException("The password reset token is invalid or expired.");
        }
    }

    public boolean passwordMatchesToken(String passwordHash, ResetTokenPayload payload) {
        return MessageDigest.isEqual(fingerprint(passwordHash).getBytes(StandardCharsets.UTF_8),
                payload.passwordFingerprint().getBytes(StandardCharsets.UTF_8));
    }

    private static String fingerprint(String value) {
        try {
            return ENCODER.encodeToString(MessageDigest.getInstance("SHA-256")
                    .digest(value.getBytes(StandardCharsets.UTF_8)));
        } catch (Exception ex) {
            throw new IllegalStateException(ex);
        }
    }

    private static boolean validSignature(String content, String supplied, String secret) {
        try {
            return MessageDigest.isEqual(DECODER.decode(sign(content, secret)), DECODER.decode(supplied));
        } catch (RuntimeException ex) {
            return false;
        }
    }

    private static String sign(String content, String secret) {
        try {
            Mac mac = Mac.getInstance("HmacSHA256");
            mac.init(new SecretKeySpec(secret.getBytes(StandardCharsets.UTF_8), "HmacSHA256"));
            return ENCODER.encodeToString(mac.doFinal(content.getBytes(StandardCharsets.UTF_8)));
        } catch (Exception ex) {
            throw new IllegalStateException("Token signing is unavailable.", ex);
        }
    }

    public record ResetTokenPayload(long userId, String passwordFingerprint) {}
}
