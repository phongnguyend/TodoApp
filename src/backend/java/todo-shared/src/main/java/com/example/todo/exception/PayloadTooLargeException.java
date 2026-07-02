package com.example.todo.exception;

/**
 * Thrown when an uploaded file exceeds the configured maximum size.
 * Mapped to HTTP 413 (Payload Too Large) by the API's global exception handler.
 */
public class PayloadTooLargeException extends RuntimeException {
    public PayloadTooLargeException(String message) {
        super(message);
    }
}
