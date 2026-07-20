import { Test, TestingModule } from '@nestjs/testing';
import { TodoItemAttachmentResponseDto } from './dto/todo-item-attachment-response.dto';
import { TodoItemAttachmentsController } from './todo-item-attachments.controller';
import { TodoItemAttachmentsService } from './todo-item-attachments.service';

const makeAttachment = (
  overrides: Partial<TodoItemAttachmentResponseDto> = {},
): TodoItemAttachmentResponseDto => ({
  id: 3,
  todoItemId: 10,
  fileId: 5,
  createdAt: new Date('2026-07-18T00:00:00Z'),
  createdByUserId: null,
  updatedAt: null,
  updatedByUserId: null,
  ...overrides,
});

describe('TodoItemAttachmentsController', () => {
  let controller: TodoItemAttachmentsController;
  let service: jest.Mocked<TodoItemAttachmentsService>;

  beforeEach(async () => {
    const mockService = {
      getAll: jest.fn(),
      getById: jest.fn(),
      create: jest.fn(),
      update: jest.fn(),
      delete: jest.fn(),
    } as unknown as jest.Mocked<TodoItemAttachmentsService>;

    const module: TestingModule = await Test.createTestingModule({
      controllers: [TodoItemAttachmentsController],
      providers: [{ provide: TodoItemAttachmentsService, useValue: mockService }],
    }).compile();

    controller = module.get(TodoItemAttachmentsController);
    service = module.get(TodoItemAttachmentsService);
  });

  describe('getAll', () => {
    it('returns all attachment references for the todo item', async () => {
      const attachments = [makeAttachment(), makeAttachment({ id: 4, fileId: 6 })];
      service.getAll.mockResolvedValue(attachments);

      await expect(controller.getAll(10)).resolves.toBe(attachments);
      expect(service.getAll).toHaveBeenCalledWith(10);
    });
  });

  describe('create', () => {
    it('creates an attachment reference using the route id and request DTO', async () => {
      const dto = { fileId: 5 };
      const created = makeAttachment();
      service.create.mockResolvedValue(created);

      await expect(controller.create(10, dto)).resolves.toBe(created);
      expect(service.create).toHaveBeenCalledWith(10, dto);
    });
  });

  describe('getById', () => {
    it('forwards both the todo item and attachment ids', async () => {
      const attachment = makeAttachment();
      service.getById.mockResolvedValue(attachment);

      await expect(controller.getById(10, 3)).resolves.toBe(attachment);
      expect(service.getById).toHaveBeenCalledWith(10, 3);
    });
  });

  describe('update', () => {
    it('updates an attachment reference using both route ids and the DTO', async () => {
      const dto = { fileId: 6 };
      const updated = makeAttachment({ fileId: 6, updatedAt: new Date('2026-07-18T01:00:00Z') });
      service.update.mockResolvedValue(updated);

      await expect(controller.update(10, 3, dto)).resolves.toBe(updated);
      expect(service.update).toHaveBeenCalledWith(10, 3, dto);
    });
  });

  describe('delete', () => {
    it('removes the attachment reference scoped to the todo item', async () => {
      service.delete.mockResolvedValue(undefined);

      await expect(controller.delete(10, 3)).resolves.toBeUndefined();
      expect(service.delete).toHaveBeenCalledWith(10, 3);
    });
  });
});
