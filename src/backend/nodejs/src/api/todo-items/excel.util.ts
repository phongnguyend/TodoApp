import * as ExcelJS from 'exceljs';

/**
 * Excel (.xlsx) reader/writer helpers used for todo item import/export.
 * Mirrors the shape of csv.util.ts but backed by the `exceljs` library.
 */

export const EXCEL_SHEET_NAME = 'Todo Items';

/**
 * Reads the first worksheet of an .xlsx workbook and returns its rows as
 * arrays of cell values (row 0 is the header row), matching the shape
 * produced by `parseCsv`.
 */
export async function parseExcel(buffer: Buffer): Promise<unknown[][]> {
  const workbook = new ExcelJS.Workbook();
  await workbook.xlsx.load(buffer as unknown as ArrayBuffer);
  const sheet = workbook.worksheets[0];
  if (!sheet) {
    return [];
  }

  const rows: unknown[][] = [];
  sheet.eachRow((row) => {
    // `row.values` is 1-indexed with index 0 unused; slice it off.
    const values = (row.values as unknown[]).slice(1);
    rows.push(values);
  });
  return rows;
}

/**
 * Builds an .xlsx workbook containing a single sheet with the given header
 * and data rows, returning the raw file bytes.
 */
export async function buildExcelWorkbook(
  header: string[],
  rows: Array<Array<string | number | boolean | null>>,
): Promise<Buffer> {
  const workbook = new ExcelJS.Workbook();
  const sheet = workbook.addWorksheet(EXCEL_SHEET_NAME);
  sheet.addRow(header);
  for (const row of rows) {
    sheet.addRow(row);
  }
  const arrayBuffer = await workbook.xlsx.writeBuffer();
  return Buffer.from(arrayBuffer);
}
