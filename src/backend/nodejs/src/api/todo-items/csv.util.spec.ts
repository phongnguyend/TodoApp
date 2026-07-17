import { parseCsv, toCsvRow } from './csv.util';

describe('csv.util', () => {
  describe('parseCsv', () => {
    it('should parse simple comma-separated rows', () => {
      const rows = parseCsv('title,description\nBuy milk,Whole milk\n');

      expect(rows).toEqual([
        ['title', 'description'],
        ['Buy milk', 'Whole milk'],
      ]);
    });

    it('should handle quoted fields with embedded commas', () => {
      const rows = parseCsv('title,description\n"Buy milk, eggs",Whole milk\n');

      expect(rows).toEqual([
        ['title', 'description'],
        ['Buy milk, eggs', 'Whole milk'],
      ]);
    });

    it('should unescape doubled quotes inside quoted fields', () => {
      const rows = parseCsv('title\n"Say ""hi"""\n');

      expect(rows).toEqual([['title'], ['Say "hi"']]);
    });

    it('should handle embedded newlines inside quoted fields', () => {
      const rows = parseCsv('title\n"Line1\nLine2"\n');

      expect(rows).toEqual([['title'], ['Line1\nLine2']]);
    });

    it('should skip blank lines', () => {
      const rows = parseCsv('title\nBuy milk\n\nWalk dog\n');

      expect(rows).toEqual([['title'], ['Buy milk'], ['Walk dog']]);
    });

    it('should parse a trailing row without a final newline', () => {
      const rows = parseCsv('title\nBuy milk');

      expect(rows).toEqual([['title'], ['Buy milk']]);
    });

    it('should return an empty array for an empty string', () => {
      expect(parseCsv('')).toEqual([]);
    });
  });

  describe('toCsvRow', () => {
    it('should join values with commas and a trailing CRLF', () => {
      expect(toCsvRow(['a', 'b', 1, true])).toBe('a,b,1,true\r\n');
    });

    it('should quote values containing commas, quotes, or newlines', () => {
      expect(toCsvRow(['a,b'])).toBe('"a,b"\r\n');
      expect(toCsvRow(['a"b'])).toBe('"a""b"\r\n');
      expect(toCsvRow(['a\nb'])).toBe('"a\nb"\r\n');
    });
  });
});
