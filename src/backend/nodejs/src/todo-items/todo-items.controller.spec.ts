import { Test, TestingModule } from '@nestjs/testing';

import { PaginatedResponseDto } from '../common/dto/paginated-response.dto';
import { TodoItemResponseDto } from './dto/todo-item-response.dto';
import { TodoItemsController } from './todo-items.controller';
import { TodoItemsService } from './todo-items.service';

const makeResponseDto = (overrides: Partial<TodoItemResponseDto> = {}): TodoItemResponseDto => ({
  id: 1,
  title: 'Test Todo',
  description: null,
  isCompleted: false,
  createdAt: new Date('2024-01-01T00:00:00Z'),
  updatedAt: null,
  ...overrides,
});

const makePaginated = (items: TodoItemResponseDto[]): PaginatedResponseDto<TodoItemResponseDto> => ({
  items,
  total: items.length,
  page: 1,
  pageSize: 20,
  totalPages: 1,
});

describe('TodoItemsController', () => {
  let controller: TodoItemsController;
  let service: jest.Mocked<TodoItemsService>;

  beforeEach(async () => {
    const mockService = {
      getAll: jest.fn(),
      getIncomplete: jest.fn(),
      getById: jest.fn(),
      create: jest.fn(),
      update: jest.fn(),
      markComplete: jest.fn(),
      delete: jest.fn(),
    } as unknown as jest.Mocked<TodoItemsService>;

    const module: TestingModule = await Test.createTestingModule({
      controllers: [TodoItemsController],
      providers: [{ provide: TodoItemsService, useValue: mockService }],
    }).compile();

    controller = module.get<TodoItemsController>(TodoItemsController);
    service = module.get(TodoItemsService);
  });

  // ── getAll ────────────────────────────────────────────────────────────────────

  describe('getAll', () => {
    it('should call service.getAll and return the result', async () => {
      const response = makePaginated([makeResponseDto(), makeResponseDto({ id: 2 })]);
      service.getAll.mockResolvedValue(response);

      const result = await controller.getAll(1, 20);

      expect(service.getAll).toHaveBeenCalledWith(1, 20);
      expect(result).toBe(response);
    });

    it('should forward custom pagination parameters', async () => {
      service.getAll.mockResolvedValue(makePaginated([]));

      await controller.getAll(3, 5);

      expect(service.getAll).toHaveBeenCalledWith(3, 5);
    });
  });

  // ── getIncomplete ─────────────────────────────────────────────────────────────

  describe('getIncomplete', () => {
    it('should call service.getIncomplete and return the result', async () => {
      const response = makePaginated([makeResponseDto({ isCompleted: false })]);
      service.getIncomplete.mockResolvedValue(response);

      const result = await controller.getIncomplete(1, 20);

      expect(service.getIncomplete).toHaveBeenCalledWith(1, 20);
      expect(result).toBe(response);
    });
  });

  // ── getById ───────────────────────────────────────────────────────────────────

  describe('getById', () => {
    it('should return a single todo item', async () => {
      const dto = makeResponseDto({ id: 3, title: 'Specific item' });
      service.getById.mockResolvedValue(dto);

      const result = await controller.getById(3);

      expect(service.getById).toHaveBeenCalledWith(3);
      expect(result).toBe(dto);
    });
  });

  // ── create ────────────────────────────────────────────────────────────────────

  describe('create', () => {
    it('should create a todo item and return it', async () => {
      const dto = makeResponseDto({ title: 'New Todo', description: 'details' });
      service.create.mockResolvedValue(dto);

      const result = await controller.create({ title: 'New Todo', description: 'details' });

      expect(service.create).toHaveBeenCalledWith({ title: 'New Todo', description: 'details' });
      expect(result).toBe(dto);
    });
  });

  // ── update ────────────────────────────────────────────────────────────────────

  describe('update', () => {
    it('should update a todo item and return it', async () => {
      const dto = makeResponseDto({ title: 'Updated Title' });
      service.update.mockResolvedValue(dto);

      const result = await controller.update(1, { title: 'Updated Title' });

      expect(service.update).toHaveBeenCalledWith(1, { title: 'Updated Title' });
      expect(result).toBe(dto);
    });
  });

  // ── markComplete ──────────────────────────────────────────────────────────────

  describe('markComplete', () => {
    it('should mark a todo item as complete', async () => {
      const dto = makeResponseDto({ isCompleted: true });
      service.markComplete.mockResolvedValue(dto);

      const result = await controller.markComplete(1);

      expect(service.markComplete).toHaveBeenCalledWith(1);
      expect(result.isCompleted).toBe(true);
    });
  });

  // ── delete ────────────────────────────────────────────────────────────────────

  describe('delete', () => {
    it('should delete a todo item', async () => {
      service.delete.mockResolvedValue(undefined);

      await controller.delete(1);

      expect(service.delete).toHaveBeenCalledWith(1);
    });
  });
});
