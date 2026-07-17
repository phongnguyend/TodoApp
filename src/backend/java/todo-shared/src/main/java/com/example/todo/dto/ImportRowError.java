package com.example.todo.dto;

/**
 * Describes a single row that failed validation during a CSV/Excel import.
 */
public record ImportRowError(
        int row,
        String error) {
}
