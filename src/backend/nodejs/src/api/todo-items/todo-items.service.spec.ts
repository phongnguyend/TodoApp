import { NotFoundException } from '@nestjs/common';
import { Test, TestingModule } from '@nestjs/testing';
import { TodoItem } from '@prisma/client';

import { TodoItemRepository } from './todo-items.repository';
import { TodoItemsService } from './todo-items.service';

const makeTodoItem = (overrides: Partial<TodoItem> = {}): TodoItem => ({
  id: 1,
  title: 'Test Todo',
  description: null,
  isCompleted: false,
  createdAt: new Date('2024-01-01T00:00:00Z'),
  updatedAt: null,
  ...overrides,
});

describe('TodoItemsService', () => {
  let service: TodoItemsService;
  let repository: jest.Mocked<TodoItemRepository>;

  beforeEach(async () => {
    const mockRepository = {
      findAll: jest.fn(),
      findIncomplete: jest.fn(),
      findById: jest.fn(),
      findAllOrdered: jest.fn(),
      create: jest.fn(),
      update: jest.fn(),
      delete: jest.fn(),
    } as unknown as jest.Mocked<TodoItemRepository>;

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        TodoItemsService,
        { provide: TodoItemRepository, useValue: mockRepository },
      ],
    }).compile();

    service = module.get<TodoItemsService>(TodoItemsService);
    repository = module.get(TodoItemRepository);
  });

  // ── getAll ────────────────────────────────────────────────────────────────────

  describe('getAll', () => {
    it('should return paginated todo items', async () => {
      const items = [makeTodoItem(), makeTodoItem({ id: 2, title: 'Second' })];
      repository.findAll.mockResolvedValue({ items, total: 2 });

      const result = await service.getAll(1, 20);

      expect(repository.findAll).toHaveBeenCalledWith(0, 20);
      expect(result.items).toHaveLength(2);
      expect(result.total).toBe(2);
      expect(result.page).toBe(1);
      expect(result.pageSize).toBe(20);
      expect(result.totalPages).toBe(1);
    });

    it('should calculate skip correctly for page 2', async () => {
      repository.findAll.mockResolvedValue({ items: [], total: 25 });

      await service.getAll(2, 10);

      expect(repository.findAll).toHaveBeenCalledWith(10, 10);
    });

    it('should calculate totalPages correctly', async () => {
      repository.findAll.mockResolvedValue({ items: [], total: 25 });

      const result = await service.getAll(1, 10);

      expect(result.totalPages).toBe(3);
    });

    it('should map items to DTOs', async () => {
      const item = makeTodoItem({ id: 7, title: 'Mapped', description: 'desc', isCompleted: true });
      repository.findAll.mockResolvedValue({ items: [item], total: 1 });

      const result = await service.getAll(1, 20);

      expect(result.items[0]).toEqual({
        id: 7,
        title: 'Mapped',
        description: 'desc',
        isCompleted: true,
        createdAt: item.createdAt,
        updatedAt: item.updatedAt,
      });
    });
  });

  // ── getIncomplete ─────────────────────────────────────────────────────────────

  describe('getIncomplete', () => {
    it('should return paginated incomplete items', async () => {
      const items = [makeTodoItem({ isCompleted: false })];
      repository.findIncomplete.mockResolvedValue({ items, total: 1 });

      const result = await service.getIncomplete(1, 20);

      expect(repository.findIncomplete).toHaveBeenCalledWith(0, 20);
      expect(result.total).toBe(1);
      expect(result.items[0].isCompleted).toBe(false);
    });

    it('should calculate skip for page 3 with pageSize 5', async () => {
      repository.findIncomplete.mockResolvedValue({ items: [], total: 0 });

      await service.getIncomplete(3, 5);

      expect(repository.findIncomplete).toHaveBeenCalledWith(10, 5);
    });
  });

  // ── getById ───────────────────────────────────────────────────────────────────

  describe('getById', () => {
    it('should return the todo item when found', async () => {
      const item = makeTodoItem({ id: 5, title: 'Find me' });
      repository.findById.mockResolvedValue(item);

      const result = await service.getById(5);

      expect(repository.findById).toHaveBeenCalledWith(5);
      expect(result.id).toBe(5);
      expect(result.title).toBe('Find me');
    });

    it('should throw NotFoundException when item does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.getById(99)).rejects.toThrow(NotFoundException);
      await expect(service.getById(99)).rejects.toThrow('Todo item 99 not found.');
    });
  });

  // ── create ────────────────────────────────────────────────────────────────────

  describe('create', () => {
    it('should create and return the new todo item', async () => {
      const item = makeTodoItem({ title: 'New Todo', description: 'desc' });
      repository.create.mockResolvedValue(item);

      const result = await service.create({ title: 'New Todo', description: 'desc' });

      expect(repository.create).toHaveBeenCalledWith({ title: 'New Todo', description: 'desc' });
      expect(result.title).toBe('New Todo');
      expect(result.description).toBe('desc');
    });

    it('should create without description when not provided', async () => {
      const item = makeTodoItem({ title: 'No Desc', description: null });
      repository.create.mockResolvedValue(item);

      const result = await service.create({ title: 'No Desc' });

      expect(repository.create).toHaveBeenCalledWith({ title: 'No Desc', description: undefined });
      expect(result.description).toBeNull();
    });
  });

  // ── update ────────────────────────────────────────────────────────────────────

  describe('update', () => {
    it('should update and return the todo item', async () => {
      const existing = makeTodoItem();
      const updated = makeTodoItem({ title: 'Updated' });
      repository.findById.mockResolvedValue(existing);
      repository.update.mockResolvedValue(updated);

      const result = await service.update(1, { title: 'Updated' });

      expect(repository.update).toHaveBeenCalledWith(1, {
        title: 'Updated',
        description: undefined,
        isCompleted: undefined,
      });
      expect(result.title).toBe('Updated');
    });

    it('should throw NotFoundException when item does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.update(99, { title: 'X' })).rejects.toThrow(NotFoundException);
    });
  });

  // ── markComplete ──────────────────────────────────────────────────────────────

  describe('markComplete', () => {
    it('should mark the item as complete', async () => {
      const existing = makeTodoItem({ isCompleted: false });
      const completed = makeTodoItem({ isCompleted: true });
      repository.findById.mockResolvedValue(existing);
      repository.update.mockResolvedValue(completed);

      const result = await service.markComplete(1);

      expect(repository.update).toHaveBeenCalledWith(1, { isCompleted: true });
      expect(result.isCompleted).toBe(true);
    });

    it('should throw NotFoundException when item does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.markComplete(99)).rejects.toThrow(NotFoundException);
    });
  });

  // ── delete ────────────────────────────────────────────────────────────────────

  describe('delete', () => {
    it('should delete the todo item', async () => {
      repository.findById.mockResolvedValue(makeTodoItem());
      repository.delete.mockResolvedValue(undefined);

      await service.delete(1);

      expect(repository.delete).toHaveBeenCalledWith(1);
    });

    it('should throw NotFoundException when item does not exist', async () => {
      repository.findById.mockResolvedValue(null);

      await expect(service.delete(99)).rejects.toThrow(NotFoundException);
    });
  });

  // ── importCsv ─────────────────────────────────────────────────────────────────

  describe('importCsv', () => {
    it('should import valid rows and skip rows with missing titles', async () => {
      repository.create.mockResolvedValue(makeTodoItem());
      const csv =
        'title,description,is_completed\n' +
        'Buy milk,Whole milk,false\n' +
        ',No title,true\n' +
        'Walk dog,,true\n';

      const result = await service.importCsv(Buffer.from(csv, 'utf-8'));

      expect(repository.create).toHaveBeenCalledTimes(2);
      expect(repository.create).toHaveBeenNthCalledWith(1, {
        title: 'Buy milk',
        description: 'Whole milk',
        isCompleted: false,
      });
      expect(repository.create).toHaveBeenNthCalledWith(2, {
        title: 'Walk dog',
        description: null,
        isCompleted: true,
      });
      expect(result).toEqual({
        imported: 2,
        failed: 1,
        errors: [{ row: 3, error: 'Title is required.' }],
      });
    });

    it('should strip a UTF-8 BOM before parsing', async () => {
      repository.create.mockResolvedValue(makeTodoItem());
      const csv = '\uFEFFtitle,description,is_completed\nBuy milk,,false\n';

      const result = await service.importCsv(Buffer.from(csv, 'utf-8'));

      expect(repository.create).toHaveBeenCalledWith({
        title: 'Buy milk',
        description: null,
        isCompleted: false,
      });
      expect(result.imported).toBe(1);
    });

    it('should return an empty result for an empty file', async () => {
      const result = await service.importCsv(Buffer.from('', 'utf-8'));

      expect(repository.create).not.toHaveBeenCalled();
      expect(result).toEqual({ imported: 0, failed: 0, errors: [] });
    });
  });

  // ── exportCsv ─────────────────────────────────────────────────────────────────

  describe('exportCsv', () => {
    it('should render all items as a CSV document', async () => {
      repository.findAllOrdered.mockResolvedValue([
        makeTodoItem({
          id: 1,
          title: 'Buy milk',
          description: 'Whole milk',
          createdAt: new Date('2024-01-01T00:00:00.000Z'),
        }),
      ]);

      const csv = await service.exportCsv();

      expect(csv).toBe(
        'id,title,description,is_completed,created_at,updated_at\r\n' +
          '1,Buy milk,Whole milk,false,2024-01-01T00:00:00.000Z,\r\n',
      );
    });

    it('should quote fields containing commas', async () => {
      repository.findAllOrdered.mockResolvedValue([
        makeTodoItem({ id: 1, title: 'Buy milk, eggs' }),
      ]);

      const csv = await service.exportCsv();

      expect(csv).toContain('"Buy milk, eggs"');
    });
  });
});
