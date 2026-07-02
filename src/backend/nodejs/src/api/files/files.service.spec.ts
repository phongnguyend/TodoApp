import { NotFoundException, PayloadTooLargeException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Test, TestingModule } from '@nestjs/testing';
import { File } from '@prisma/client';
import { promises as fs } from 'fs';

import { FileRepository } from './files.repository';
import { FilesService } from './files.service';

jest.mock('fs', () => ({
  ...jest.requireActual('fs'),
  promises: {
    mkdir: jest.fn(),
    writeFile: jest.fn(),
    unlink: jest.fn(),
    access: jest.fn(),
  },
}));

const fsMock = fs as unknown as {
  mkdir: jest.Mock;
  writeFile: jest.Mock;
  unlink: jest.Mock;
  access: jest.Mock;
};

const makeFile = (overrides: Partial<File> = {}): File => ({
  id: 1,
  name: 'report.pdf',
  extension: 'pdf',
  size: 1024,
  contentType: 'application/pdf',
  location: '/uploads/abc_report.pdf',
  createdAt: new Date('2024-01-01T00:00:00Z'),
  updatedAt: null,
  ...overrides,
});

const makeMulterFile = (overrides: Partial<Express.Multer.File> = {}): Express.Multer.File =>
  ({
    fieldname: 'file',
    originalname: 'report.pdf',
    encoding: '7bit',
    mimetype: 'application/pdf',
    size: 1024,
    buffer: Buffer.from('file-content'),
    destination: '',
    filename: '',
    path: '',
    stream: undefined as never,
    ...overrides,
  }) as Express.Multer.File;

describe('FilesService', () => {
  let service: FilesService;
  let repository: jest.Mocked<FileRepository>;

  beforeEach(async () => {
    jest.clearAllMocks();

    const mockRepository = {
      findAll: jest.fn(),
      findById: jest.fn(),
      create: jest.fn(),
      delete: jest.fn(),
    } as unknown as jest.Mocked<FileRepository>;

    const mockConfig = {
      get: jest.fn((key: string) => {
        if (key === 'FILE_STORAGE_PATH') return './uploads';
        if (key === 'MAX_UPLOAD_SIZE_BYTES') return String(10 * 1024 * 1024);
        return undefined;
      }),
    } as unknown as ConfigService;

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        FilesService,
        { provide: FileRepository, useValue: mockRepository },
        { provide: ConfigService, useValue: mockConfig },
      ],
    }).compile();

    service = module.get<FilesService>(FilesService);
    repository = module.get(FileRepository);
  });

  // ── getAll ────────────────────────────────────────────────────────────────────

  describe('getAll', () => {
    it('should return paginated files', async () => {
      const items = [makeFile(), makeFile({ id: 2 })];
      repository.findAll.mockResolvedValue({ items, total: 2 });

      const result = await service.getAll(1, 20);

      expect(repository.findAll).toHaveBeenCalledWith(0, 20);
      expect(result.items).toHaveLength(2);
      expect(result.total).toBe(2);
      expect(result.page).toBe(1);
      expect(result.pageSize).toBe(20);
      expect(result.totalPages).toBe(1);
    });

    it('should not expose the on-disk location', async () => {
      repository.findAll.mockResolvedValue({ items: [makeFile()], total: 1 });

      const result = await service.getAll(1, 20);

      expect(result.items[0]).not.toHaveProperty('location');
    });
  });

  // ── getById ───────────────────────────────────────────────────────────────────

  describe('getById', () => {
    it('should return the file when found', async () => {
      const file = makeFile({ id: 5, name: 'find-me.txt' });
      repository.findById.mockResolvedValue(file);

      const result = await service.getById(5);

      expect(result.id).toBe(5);
      expect(result.name).toBe('find-me.txt');
    });

    it('should throw NotFoundException when file does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.getById(99)).rejects.toThrow(NotFoundException);
      await expect(service.getById(99)).rejects.toThrow('File 99 not found.');
    });
  });

  // ── upload ────────────────────────────────────────────────────────────────────

  describe('upload', () => {
    it('should store the file content and create a record', async () => {
      const created = makeFile({ name: 'report.pdf', extension: 'pdf', size: 1024 });
      repository.create.mockResolvedValue(created);
      fsMock.mkdir.mockResolvedValue(undefined);
      fsMock.writeFile.mockResolvedValue(undefined);

      const result = await service.upload(makeMulterFile());

      expect(fsMock.mkdir).toHaveBeenCalledWith('./uploads', { recursive: true });
      expect(fsMock.writeFile).toHaveBeenCalled();
      expect(repository.create).toHaveBeenCalledWith(
        expect.objectContaining({
          name: 'report.pdf',
          extension: 'pdf',
          size: 1024,
          contentType: 'application/pdf',
        }),
      );
      expect(result.name).toBe('report.pdf');
    });

    it('should strip directory components from the original name', async () => {
      repository.create.mockResolvedValue(makeFile());
      fsMock.mkdir.mockResolvedValue(undefined);
      fsMock.writeFile.mockResolvedValue(undefined);

      await service.upload(makeMulterFile({ originalname: '../../etc/passwd' }));

      expect(repository.create).toHaveBeenCalledWith(
        expect.objectContaining({ name: 'passwd', extension: '' }),
      );
    });

    it('should throw PayloadTooLargeException when the file exceeds the max size', async () => {
      await expect(
        service.upload(makeMulterFile({ size: 20 * 1024 * 1024 })),
      ).rejects.toThrow(PayloadTooLargeException);

      expect(repository.create).not.toHaveBeenCalled();
    });
  });

  // ── getDownloadTarget ─────────────────────────────────────────────────────────

  describe('getDownloadTarget', () => {
    it('should return the download target when the file exists on disk', async () => {
      repository.findById.mockResolvedValue(makeFile());
      fsMock.access.mockResolvedValue(undefined);

      const result = await service.getDownloadTarget(1);

      expect(result).toEqual({
        path: '/uploads/abc_report.pdf',
        name: 'report.pdf',
        contentType: 'application/pdf',
      });
    });

    it('should default to application/octet-stream when contentType is missing', async () => {
      repository.findById.mockResolvedValue(makeFile({ contentType: null }));
      fsMock.access.mockResolvedValue(undefined);

      const result = await service.getDownloadTarget(1);

      expect(result.contentType).toBe('application/octet-stream');
    });

    it('should throw NotFoundException when the record does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.getDownloadTarget(99)).rejects.toThrow(NotFoundException);
    });

    it('should throw NotFoundException when the content is missing on disk', async () => {
      repository.findById.mockResolvedValue(makeFile());
      fsMock.access.mockRejectedValue(new Error('ENOENT'));

      await expect(service.getDownloadTarget(1)).rejects.toThrow(
        'File 1 content not found on disk.',
      );
    });
  });

  // ── delete ────────────────────────────────────────────────────────────────────

  describe('delete', () => {
    it('should delete the record and remove the file from disk', async () => {
      repository.findById.mockResolvedValue(makeFile());
      repository.delete.mockResolvedValue(undefined);
      fsMock.unlink.mockResolvedValue(undefined);

      await service.delete(1);

      expect(repository.delete).toHaveBeenCalledWith(1);
      expect(fsMock.unlink).toHaveBeenCalledWith('/uploads/abc_report.pdf');
    });

    it('should not throw when the file is already missing on disk', async () => {
      repository.findById.mockResolvedValue(makeFile());
      repository.delete.mockResolvedValue(undefined);
      fsMock.unlink.mockRejectedValue(new Error('ENOENT'));

      await expect(service.delete(1)).resolves.toBeUndefined();
    });

    it('should throw NotFoundException when file does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.delete(99)).rejects.toThrow(NotFoundException);
      expect(repository.delete).not.toHaveBeenCalled();
    });
  });
});
