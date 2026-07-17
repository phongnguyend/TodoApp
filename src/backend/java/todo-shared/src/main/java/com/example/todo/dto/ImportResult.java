package com.example.todo.dto;

import java.util.List;

/**
 * Summary returned after importing todo items from a CSV/Excel file.
 */
public record ImportResult(
        int imported,
        int failed,
        List<ImportRowError> errors) {
}
