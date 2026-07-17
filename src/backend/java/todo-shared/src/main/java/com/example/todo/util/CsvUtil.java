package com.example.todo.util;

import java.util.ArrayList;
import java.util.List;

/**
 * Minimal RFC 4180 CSV parser/writer used for todo item import/export.
 * Handles quoted fields with embedded commas, quotes, and newlines so behavior
 * matches the standard CSV libraries used by the other language
 * implementations.
 */
public final class CsvUtil {

    private CsvUtil() {
    }

    public static List<List<String>> parse(String text) {
        List<List<String>> rows = new ArrayList<>();
        List<String> row = new ArrayList<>();
        StringBuilder field = new StringBuilder();
        boolean inQuotes = false;
        int i = 0;

        while (i < text.length()) {
            char c = text.charAt(i);

            if (inQuotes) {
                if (c == '"') {
                    if (i + 1 < text.length() && text.charAt(i + 1) == '"') {
                        field.append('"');
                        i += 2;
                        continue;
                    }
                    inQuotes = false;
                    i++;
                    continue;
                }
                field.append(c);
                i++;
                continue;
            }

            if (c == '"') {
                inQuotes = true;
                i++;
                continue;
            }
            if (c == ',') {
                row.add(field.toString());
                field.setLength(0);
                i++;
                continue;
            }
            if (c == '\r') {
                i++;
                continue;
            }
            if (c == '\n') {
                row.add(field.toString());
                rows.add(row);
                row = new ArrayList<>();
                field.setLength(0);
                i++;
                continue;
            }
            field.append(c);
            i++;
        }

        // Flush the trailing field/row when the input doesn't end with a newline.
        if (!field.isEmpty() || !row.isEmpty()) {
            row.add(field.toString());
            rows.add(row);
        }

        // Drop blank lines (they parse as a single empty-string field).
        rows.removeIf(r -> r.size() == 1 && r.get(0).isEmpty());
        return rows;
    }

    private static String escapeField(String value) {
        if (value.indexOf('"') >= 0 || value.indexOf(',') >= 0 || value.indexOf('\r') >= 0
                || value.indexOf('\n') >= 0) {
            return "\"" + value.replace("\"", "\"\"") + "\"";
        }
        return value;
    }

    public static String toCsvRow(List<?> values) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < values.size(); i++) {
            if (i > 0) {
                sb.append(',');
            }
            Object value = values.get(i);
            sb.append(escapeField(value == null ? "" : String.valueOf(value)));
        }
        sb.append("\r\n");
        return sb.toString();
    }
}
