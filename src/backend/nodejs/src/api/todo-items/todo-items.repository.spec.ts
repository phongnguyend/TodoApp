import { Test, TestingModule } from '@nestjs/testing';
import { TodoItem } from '@prisma/client';

import { PrismaService } from '../../shared/prisma/prisma.service';
import { TodoItemRepository } from './todo-items.repository';

const makeTodoItem = (overrides: Partial<TodoItem> = {}): TodoItem => ({
  id: 1,
  title: 'Test Todo',
  description: null,
  isCompleted: false,
  createdAt: new Date('2024-01-01T00:00:00Z'),
  updatedAt: null,
  ...overrides,
});

describe('TodoItemRepository', () => {
  let repository: TodoItemRepository;
  let prismaMock: {
    todoItem: {
      findMany: jest.Mock;
      count: jest.Mock;
      findUnique: jest.Mock;
      create: jest.Mock;
      update: jest.Mock;
      delete: jest.Mock;
    };
    $transaction: jest.Mock;
  };

  beforeEach(async () => {
    prismaMock = {
      todoItem: {
        findMany: jest.fn(),
        count: jest.fn(),
        findUnique: jest.fn(),
        create: jest.fn(),
        update: jest.fn(),
        delete: jest.fn(),
      },
      $transaction: jest.fn(),
    };

    const module: TestingModule = await Test.createTestingModule({
      providers: [
        TodoItemRepository,
        { provide: PrismaService, useValue: prismaMock },
      ],
    }).compile();

    repository = module.get<TodoItemRepository>(TodoItemRepository);
  });

  // ── findAll ───────────────────────────────────────────────────────────────────

  describe('findAll', () => {
    it('should return items and total count', async () => {
      const items = [makeTodoItem(), makeTodoItem({ id: 2 })];
      prismaMock.todoItem.findMany.mockReturnValue('findManyQuery');
      prismaMock.todoItem.count.mockReturnValue('countQuery');
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
      prismaMock.todoItem.findMany.mockReturnValue('findManyQuery');
      prismaMock.todoItem.count.mockReturnValue('countQuery');

      await repository.findAll(10, 5);

      expect(prismaMock.todoItem.findMany).toHaveBeenCalledWith({
        skip: 10,
        take: 5,
        orderBy: { createdAt: 'desc' },
      });
    });
  });

  // ── findIncomplete ────────────────────────────────────────────────────────────

  describe('findIncomplete', () => {
    it('should return incomplete items and total', async () => {
      const items = [makeTodoItem({ isCompleted: false })];
      prismaMock.$transaction.mockResolvedValue([items, 1]);

      const result = await repository.findIncomplete(0, 20);

      expect(result.items).toEqual(items);
      expect(result.total).toBe(1);
    });

    it('should filter by isCompleted: false', async () => {
      prismaMock.$transaction.mockResolvedValue([[], 0]);
      prismaMock.todoItem.findMany.mockReturnValue('findManyQuery');
      prismaMock.todoItem.count.mockReturnValue('countQuery');

      await repository.findIncomplete(0, 10);

      expect(prismaMock.todoItem.findMany).toHaveBeenCalledWith(
        expect.objectContaining({ where: { isCompleted: false } }),
      );
      expect(prismaMock.todoItem.count).toHaveBeenCalledWith(
        expect.objectContaining({ where: { isCompleted: false } }),
      );
    });
  });

  // ── findAllOrdered ────────────────────────────────────────────────────────────

  describe('findAllOrdered', () => {
    it('should return all items ordered by createdAt descending', async () => {
      const items = [makeTodoItem(), makeTodoItem({ id: 2 })];
      prismaMock.todoItem.findMany.mockResolvedValue(items);

      const result = await repository.findAllOrdered();

      expect(prismaMock.todoItem.findMany).toHaveBeenCalledWith({
        orderBy: { createdAt: 'desc' },
      });
      expect(result).toEqual(items);
    });
  });

  // ── findById ──────────────────────────────────────────────────────────────────

  describe('findById', () => {
    it('should return the item when found', async () => {
      const item = makeTodoItem({ id: 5 });
      prismaMock.todoItem.findUnique.mockResolvedValue(item);

      const result = await repository.findById(5);

      expect(prismaMock.todoItem.findUnique).toHaveBeenCalledWith({ where: { id: 5 } });
      expect(result).toEqual(item);
    });

    it('should return null when item is not found', async () => {
      prismaMock.todoItem.findUnique.mockResolvedValue(null);

      const result = await repository.findById(99);

      expect(result).toBeNull();
    });
  });

  // ── create ────────────────────────────────────────────────────────────────────

  describe('create', () => {
    it('should create and return a new todo item', async () => {
      const item = makeTodoItem({ title: 'New Todo', description: 'details' });
      prismaMock.todoItem.create.mockResolvedValue(item);

      const result = await repository.create({ title: 'New Todo', description: 'details' });

      expect(prismaMock.todoItem.create).toHaveBeenCalledWith({
        data: { title: 'New Todo', description: 'details' },
      });
      expect(result).toEqual(item);
    });
  });

  // ── update ────────────────────────────────────────────────────────────────────

  describe('update', () => {
    it('should update and return the todo item', async () => {
      const updated = makeTodoItem({ title: 'Updated Title' });
      prismaMock.todoItem.update.mockResolvedValue(updated);

      const result = await repository.update(1, { title: 'Updated Title' });

      expect(prismaMock.todoItem.update).toHaveBeenCalledWith({
        where: { id: 1 },
        data: { title: 'Updated Title' },
      });
      expect(result).toEqual(updated);
    });

    it('should update isCompleted flag', async () => {
      const completed = makeTodoItem({ isCompleted: true });
      prismaMock.todoItem.update.mockResolvedValue(completed);

      const result = await repository.update(1, { isCompleted: true });

      expect(prismaMock.todoItem.update).toHaveBeenCalledWith({
        where: { id: 1 },
        data: { isCompleted: true },
      });
      expect(result.isCompleted).toBe(true);
    });
  });

  // ── delete ────────────────────────────────────────────────────────────────────

  describe('delete', () => {
    it('should delete the todo item', async () => {
      prismaMock.todoItem.delete.mockResolvedValue(undefined);

      await repository.delete(1);

      expect(prismaMock.todoItem.delete).toHaveBeenCalledWith({ where: { id: 1 } });
    });
  });
});
