package com.example.todo.security;

import com.example.todo.exception.InvalidPasswordResetTokenException;
import com.fasterxml.jackson.databind.JsonNode;
import com.fasterxml.jackson.databind.ObjectMapper;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.stereotype.Component;
import org.springframework.security.oauth2.jose.jws.MacAlgorithm;
import org.springframework.security.oauth2.jwt.JwtClaimsSet;
import org.springframework.security.oauth2.jwt.JwtEncoder;
import org.springframework.security.oauth2.jwt.JwtEncoderParameters;
import org.springframework.security.oauth2.jwt.JwsHeader;
import org.springframework.security.oauth2.jwt.NimbusJwtEncoder;
import com.nimbusds.jose.jwk.source.ImmutableSecret;

import javax.crypto.Mac;
import javax.crypto.spec.SecretKeySpec;
import javax.crypto.SecretKey;
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
    private final JwtEncoder jwtEncoder;
    private final String resetSecret;

    public UserTokenCodec(
            @Value("${app.user.jwt-secret:change-me-use-at-least-32-bytes-long}") String jwtSecret,
            @Value("${app.user.password-reset.secret:${app.user.jwt-secret:change-me-use-at-least-32-bytes-long}}") String resetSecret) {
        SecretKey accessTokenKey = new SecretKeySpec(jwtSecret.getBytes(StandardCharsets.UTF_8), "HmacSHA256");
        this.jwtEncoder = new NimbusJwtEncoder(new ImmutableSecret<>(accessTokenKey));
        this.resetSecret = resetSecret;
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
        JwsHeader header = JwsHeader.with(MacAlgorithm.HS256).type("JWT").build();
        JwtClaimsSet claims = JwtClaimsSet.builder()
                .subject(Long.toString(userId))
                .issuedAt(issuedAt)
                .expiresAt(expiresAt)
                .build();
        return jwtEncoder.encode(JwtEncoderParameters.from(header, claims)).getTokenValue();
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
