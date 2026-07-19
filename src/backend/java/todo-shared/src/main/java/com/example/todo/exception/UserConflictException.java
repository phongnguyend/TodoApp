package com.example.todo.exception;

public class UserConflictException extends RuntimeException {
    public UserConflictException(String message) { super(message); }
}
