import { buildExcelWorkbook, parseExcel } from './excel.util';

describe('excel.util', () => {
  describe('buildExcelWorkbook + parseExcel', () => {
    it('should round-trip a header and data rows', async () => {
      const buffer = await buildExcelWorkbook(
        ['id', 'title', 'description', 'is_completed', 'created_at', 'updated_at'],
        [[1, 'Buy milk', 'Whole milk', true, '2024-01-01T00:00:00.000Z', '']],
      );

      const rows = await parseExcel(buffer);

      expect(rows).toHaveLength(2);
      expect(rows[0]).toEqual(['id', 'title', 'description', 'is_completed', 'created_at', 'updated_at']);
      expect(rows[1][0]).toBe(1);
      expect(rows[1][1]).toBe('Buy milk');
      expect(rows[1][2]).toBe('Whole milk');
      expect(rows[1][3]).toBe(true);
    });

    it('should return only the header row when no data rows are provided', async () => {
      const buffer = await buildExcelWorkbook(['id', 'title'], []);

      const rows = await parseExcel(buffer);

      expect(rows).toEqual([['id', 'title']]);
    });
  });
});
