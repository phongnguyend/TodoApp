package com.example.todo.dto;

import java.util.List;

/**
 * Generic paginated response - mirrors a PagedResult&lt;T&gt; / PaginatedResponse in C#.
 *
 * @param <T> the item type
 */
public record PaginatedResponse<T>(
        List<T> items,
        long total,
        int page,
        int pageSize,
        int totalPages
) {}
