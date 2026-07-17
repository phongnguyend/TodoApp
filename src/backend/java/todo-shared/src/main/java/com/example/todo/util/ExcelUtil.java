package com.example.todo.util;

import org.apache.poi.ss.usermodel.Cell;
import org.apache.poi.ss.usermodel.DataFormatter;
import org.apache.poi.ss.usermodel.Row;
import org.apache.poi.ss.usermodel.Sheet;
import org.apache.poi.ss.usermodel.Workbook;
import org.apache.poi.ss.usermodel.WorkbookFactory;
import org.apache.poi.xssf.usermodel.XSSFWorkbook;

import java.io.ByteArrayOutputStream;
import java.io.IOException;
import java.io.InputStream;
import java.io.UncheckedIOException;
import java.util.ArrayList;
import java.util.List;

/**
 * Minimal .xlsx reader/writer (backed by Apache POI) used for todo item
 * import/export. Mirrors the shape of {@link CsvUtil} - rows are returned as
 * {@code List<List<String>>} (row 0 is the header row) - so the same
 * row-processing logic can consume either format.
 */
public final class ExcelUtil {

    public static final String SHEET_NAME = "Todo Items";

    private ExcelUtil() {
    }

    public static List<List<String>> parse(InputStream inputStream) {
        List<List<String>> rows = new ArrayList<>();
        DataFormatter formatter = new DataFormatter();

        try (Workbook workbook = WorkbookFactory.create(inputStream)) {
            if (workbook.getNumberOfSheets() == 0) {
                return rows;
            }
            Sheet sheet = workbook.getSheetAt(0);

            for (Row row : sheet) {
                List<String> values = new ArrayList<>();
                for (int c = 0; c < row.getLastCellNum(); c++) {
                    Cell cell = row.getCell(c, Row.MissingCellPolicy.RETURN_BLANK_AS_NULL);
                    values.add(cell != null ? formatter.formatCellValue(cell) : "");
                }
                // Skip blank rows, matching CsvUtil's handling of blank lines.
                if (values.stream().anyMatch(value -> !value.isBlank())) {
                    rows.add(values);
                }
            }
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }

        return rows;
    }

    public static byte[] write(List<String> header, List<List<Object>> rows) {
        try (Workbook workbook = new XSSFWorkbook(); ByteArrayOutputStream out = new ByteArrayOutputStream()) {
            Sheet sheet = workbook.createSheet(SHEET_NAME);

            Row headerRow = sheet.createRow(0);
            for (int c = 0; c < header.size(); c++) {
                headerRow.createCell(c).setCellValue(header.get(c));
            }

            for (int r = 0; r < rows.size(); r++) {
                Row dataRow = sheet.createRow(r + 1);
                List<Object> values = rows.get(r);
                for (int c = 0; c < values.size(); c++) {
                    setCellValue(dataRow.createCell(c), values.get(c));
                }
            }

            workbook.write(out);
            return out.toByteArray();
        } catch (IOException e) {
            throw new UncheckedIOException(e);
        }
    }

    private static void setCellValue(Cell cell, Object value) {
        switch (value) {
            case null -> cell.setBlank();
            case Boolean bool -> cell.setCellValue(bool);
            case Number number -> cell.setCellValue(number.doubleValue());
            default -> cell.setCellValue(String.valueOf(value));
        }
    }
}
