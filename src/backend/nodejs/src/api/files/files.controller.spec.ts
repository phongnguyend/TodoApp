import { BadRequestException } from '@nestjs/common';
import { Test, TestingModule } from '@nestjs/testing';
import { Response } from 'express';

import { FileResponseDto } from './dto/file-response.dto';
import { FilesController } from './files.controller';
import { FilesService } from './files.service';

const makeResponseDto = (overrides: Partial<FileResponseDto> = {}): FileResponseDto => ({
  id: 1,
  name: 'report.pdf',
  extension: 'pdf',
  size: 1024,
  contentType: 'application/pdf',
  createdAt: new Date('2024-01-01T00:00:00Z'),
  createdByUserId: null,
  updatedAt: null,
  updatedByUserId: null,
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

describe('FilesController', () => {
  let controller: FilesController;
  let service: jest.Mocked<FilesService>;

  beforeEach(async () => {
    const mockService = {
      getAll: jest.fn(),
      getById: jest.fn(),
      getDownloadTarget: jest.fn(),
      upload: jest.fn(),
      delete: jest.fn(),
    } as unknown as jest.Mocked<FilesService>;

    const module: TestingModule = await Test.createTestingModule({
      controllers: [FilesController],
      providers: [{ provide: FilesService, useValue: mockService }],
    }).compile();

    controller = module.get<FilesController>(FilesController);
    service = module.get(FilesService);
  });

  // ── getAll ────────────────────────────────────────────────────────────────────

  describe('getAll', () => {
    it('should call service.getAll and return the result', async () => {
      const response = {
        items: [makeResponseDto(), makeResponseDto({ id: 2 })],
        total: 2,
        page: 1,
        pageSize: 20,
        totalPages: 1,
      };
      service.getAll.mockResolvedValue(response);

      const result = await controller.getAll(1, 20);

      expect(service.getAll).toHaveBeenCalledWith(1, 20);
      expect(result).toBe(response);
    });
  });

  // ── getById ───────────────────────────────────────────────────────────────────

  describe('getById', () => {
    it('should return a single file', async () => {
      const dto = makeResponseDto({ id: 3, name: 'specific.txt' });
      service.getById.mockResolvedValue(dto);

      const result = await controller.getById(3);

      expect(service.getById).toHaveBeenCalledWith(3);
      expect(result).toBe(dto);
    });
  });

  // ── download ──────────────────────────────────────────────────────────────────

  describe('download', () => {
    it('should stream the file content via res.download', async () => {
      service.getDownloadTarget.mockResolvedValue({
        path: '/uploads/abc_report.pdf',
        name: 'report.pdf',
        contentType: 'application/pdf',
      });
      const res = { download: jest.fn() } as unknown as Response;

      await controller.download(1, res);

      expect(service.getDownloadTarget).toHaveBeenCalledWith(1);
      expect(res.download).toHaveBeenCalledWith('/uploads/abc_report.pdf', 'report.pdf', {
        headers: { 'Content-Type': 'application/pdf' },
      });
    });
  });

  // ── create ────────────────────────────────────────────────────────────────────

  describe('create', () => {
    it('should upload the file and return it', async () => {
      const dto = makeResponseDto();
      service.upload.mockResolvedValue(dto);
      const file = makeMulterFile();

      const result = await controller.create(file);

      expect(service.upload).toHaveBeenCalledWith(file);
      expect(result).toBe(dto);
    });

    it('should throw BadRequestException when no file is provided', async () => {
      expect(() => controller.create(undefined)).toThrow(BadRequestException);
      expect(service.upload).not.toHaveBeenCalled();
    });
  });

  // ── delete ────────────────────────────────────────────────────────────────────

  describe('delete', () => {
    it('should delete a file', async () => {
      service.delete.mockResolvedValue(undefined);

      await controller.delete(1);

      expect(service.delete).toHaveBeenCalledWith(1);
    });
  });
});
