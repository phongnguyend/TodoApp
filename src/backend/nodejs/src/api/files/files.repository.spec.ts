import { Test, TestingModule } from '@nestjs/testing';
import { File } from '@prisma/client';

import { PrismaService } from '../../shared/prisma/prisma.service';
import { FileRepository } from './files.repository';

const makeFile = (overrides: Partial<File> = {}): File => ({
  id: 1,
  name: 'report.pdf',
  extension: 'pdf',
  size: 1024,
  contentType: 'application/pdf',
  location: '/uploads/abc_report.pdf',
  createdAt: new Date('2024-01-01T00:00:00Z'),
  createdByUserId: null,
  updatedAt: null,
  updatedByUserId: null,
  ...overrides,
});

describe('FileRepository', () => {
  let repository: FileRepository;
  let prismaMock: {
    file: {
      findMany: jest.Mock;
      count: jest.Mock;
      findUnique: jest.Mock;
      create: jest.Mock;
      delete: jest.Mock;
    };
    $transaction: jest.Mock;
  };

  beforeEach(async () => {
    prismaMock = {
      file: {
        findMany: jest.fn(),
        count: jest.fn(),
        findUnique: jest.fn(),
        create: jest.fn(),
        delete: jest.fn(),
      },
      $transaction: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        FileRepository,
        { provide: PrismaService, useValue: prismaMock },
      ],
    }).compile();

    repository = module.get<FileRepository>(FileRepository);
  });

  // ── findAll ───────────────────────────────────────────────────────────────────

  describe('findAll', () => {
    it('should return items and total count', async () => {
      const items = [makeFile(), makeFile({ id: 2 })];
      prismaMock.file.findMany.mockReturnValue('findManyQuery');
      prismaMock.file.count.mockReturnValue('countQuery');
      prismaMock.$transaction.mockResolvedValue([items, 2]);

      const result = await repository.findAll(0, 20);

      expect(prismaMock.$transaction).toHaveBeenCalledWith([
        expect.anything(),
        expect.anything(),
      ]);
      expect(result.items).toEqual(items);
      expect(result.total).toBe(2);
    });

    it('should pass skip and take to findMany', async () => {
      prismaMock.$transaction.mockResolvedValue([[], 0]);
      prismaMock.file.findMany.mockReturnValue('findManyQuery');
      prismaMock.file.count.mockReturnValue('countQuery');

      await repository.findAll(10, 5);

      expect(prismaMock.file.findMany).toHaveBeenCalledWith({
        skip: 10,
        take: 5,
        orderBy: { createdAt: 'desc' },
      });
    });
  });

  // ── findById ──────────────────────────────────────────────────────────────────

  describe('findById', () => {
    it('should return the file when found', async () => {
      const file = makeFile({ id: 5 });
      prismaMock.file.findUnique.mockResolvedValue(file);

      const result = await repository.findById(5);

      expect(prismaMock.file.findUnique).toHaveBeenCalledWith({ where: { id: 5 } });
      expect(result).toEqual(file);
    });

    it('should return null when file is not found', async () => {
      prismaMock.file.findUnique.mockResolvedValue(null);

      const result = await repository.findById(99);

      expect(result).toBeNull();
    });
  });

  // ── create ────────────────────────────────────────────────────────────────────

  describe('create', () => {
    it('should create and return a new file', async () => {
      const file = makeFile({ name: 'new.txt', extension: 'txt' });
      prismaMock.file.create.mockResolvedValue(file);

      const data = {
        name: 'new.txt',
        extension: 'txt',
        size: 10,
        contentType: 'text/plain',
        location: '/uploads/xyz_new.txt',
      };
      const result = await repository.create(data);

      expect(prismaMock.file.create).toHaveBeenCalledWith({ data });
      expect(result).toEqual(file);
    });
  });

  // ── delete ────────────────────────────────────────────────────────────────────

  describe('delete', () => {
    it('should delete the file', async () => {
      prismaMock.file.delete.mockResolvedValue(undefined);

      await repository.delete(1);

      expect(prismaMock.file.delete).toHaveBeenCalledWith({ where: { id: 1 } });
    });
  });
});
