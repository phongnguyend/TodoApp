package com.example.todo.api.config;

import org.springframework.security.core.Authentication;
import org.springframework.security.core.context.SecurityContextHolder;

public final class AuditActor {
    private AuditActor() {}

    public static Long currentUserId() {
        Authentication authentication = SecurityContextHolder.getContext().getAuthentication();
        if (authentication == null || !authentication.isAuthenticated()) return null;
        try {
            long userId = Long.parseLong(authentication.getName());
            return userId > 0 ? userId : null;
        } catch (NumberFormatException ignored) {
            return null;
        }
    }
}
