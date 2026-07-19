package com.example.todo.api.controller;

import com.example.todo.dto.TokenRequest;
import com.example.todo.dto.TokenResponse;
import com.example.todo.exception.UnauthorizedException;
import com.example.todo.service.UserService;
import io.swagger.v3.oas.annotations.tags.Tag;
import jakarta.validation.Valid;
import lombok.RequiredArgsConstructor;
import org.springframework.http.CacheControl;
import org.springframework.http.HttpHeaders;
import org.springframework.http.HttpStatus;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.PostMapping;
import org.springframework.web.bind.annotation.RequestBody;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.util.Map;

@RestController
@RequestMapping("/api/tokens")
@RequiredArgsConstructor
@Tag(name = "Tokens")
public class TokenController {
    private final UserService service;

    @PostMapping
    public ResponseEntity<?> create(@Valid @RequestBody TokenRequest request) {
        try {
            TokenResponse token = service.createToken(request);
            return ResponseEntity.ok()
                    .cacheControl(CacheControl.noStore())
                    .header("Pragma", "no-cache")
                    .body(token);
        } catch (UnauthorizedException ex) {
            return ResponseEntity.status(HttpStatus.UNAUTHORIZED)
                    .header(HttpHeaders.WWW_AUTHENTICATE, "Bearer")
                    .body(Map.of("error", "Invalid email or password."));
        }
    }
}
