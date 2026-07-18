import { NotFoundException } from '@nestjs/common';
import { Test } from '@nestjs/testing';
import { TodoItemAttachment } from '@prisma/client';
import { FileRepository } from '../files/files.repository';
import { TodoItemRepository } from '../todo-items/todo-items.repository';
import { TodoItemAttachmentRepository } from './todo-item-attachments.repository';
import { TodoItemAttachmentsService } from './todo-item-attachments.service';

const attachment = (overrides: Partial<TodoItemAttachment> = {}): TodoItemAttachment => ({
  id: 3,
  todoItemId: 10,
  fileId: 5,
  createdAt: new Date('2026-07-18T00:00:00Z'),
  updatedAt: null,
  ...overrides,
});

describe('TodoItemAttachmentsService', () => {
  let service: TodoItemAttachmentsService;
  let repository: jest.Mocked<TodoItemAttachmentRepository>;
  let todos: jest.Mocked<TodoItemRepository>;
  let files: jest.Mocked<FileRepository>;

  beforeEach(async () => {
    repository = {
      findByTodoItemId: jest.fn(), findByIdForTodoItem: jest.fn(),
      findByTodoItemAndFile: jest.fn(), create: jest.fn(), update: jest.fn(), delete: jest.fn(),
    } as unknown as jest.Mocked<TodoItemAttachmentRepository>;
    todos = { findById: jest.fn() } as unknown as jest.Mocked<TodoItemRepository>;
    files = { findById: jest.fn() } as unknown as jest.Mocked<FileRepository>;
    const module = await Test.createTestingModule({
      providers: [
        TodoItemAttachmentsService,
        { provide: TodoItemAttachmentRepository, useValue: repository },
        { provide: TodoItemRepository, useValue: todos },
        { provide: FileRepository, useValue: files },
      ],
    }).compile();
    service = module.get(TodoItemAttachmentsService);
    todos.findById.mockResolvedValue({ id: 10 } as never);
    files.findById.mockResolvedValue({ id: 5 } as never);
  });

  it('lists attachments only after verifying the todo exists', async () => {
    repository.findByTodoItemId.mockResolvedValue([attachment()]);
    await expect(service.getAll(10)).resolves.toEqual([attachment()]);
    expect(repository.findByTodoItemId).toHaveBeenCalledWith(10);
  });

  it('rejects an unknown todo item', async () => {
    todos.findById.mockResolvedValue(null);
    await expect(service.getAll(99)).rejects.toBeInstanceOf(NotFoundException);
    expect(repository.findByTodoItemId).not.toHaveBeenCalled();
  });

  it('creates an attachment reference to an existing file', async () => {
    repository.findByTodoItemAndFile.mockResolvedValue(null);
    repository.create.mockResolvedValue(attachment());
    await expect(service.create(10, { fileId: 5 })).resolves.toEqual(attachment());
    expect(repository.create).toHaveBeenCalledWith({ todoItemId: 10, fileId: 5 });
  });

  it('returns the existing reference when creating a duplicate', async () => {
    repository.findByTodoItemAndFile.mockResolvedValue(attachment());
    await expect(service.create(10, { fileId: 5 })).resolves.toEqual(attachment());
    expect(repository.create).not.toHaveBeenCalled();
  });

  it('rejects an unknown file', async () => {
    files.findById.mockResolvedValue(null);
    await expect(service.create(10, { fileId: 99 })).rejects.toBeInstanceOf(NotFoundException);
    expect(repository.create).not.toHaveBeenCalled();
  });

  it('scopes attachment lookup to its parent todo', async () => {
    repository.findByIdForTodoItem.mockResolvedValue(null);
    await expect(service.getById(10, 3)).rejects.toBeInstanceOf(NotFoundException);
    expect(repository.findByIdForTodoItem).toHaveBeenCalledWith(10, 3);
  });

  it('updates the referenced file', async () => {
    const changed = attachment({ fileId: 6 });
    files.findById.mockResolvedValue({ id: 6 } as never);
    repository.findByIdForTodoItem.mockResolvedValue(attachment());
    repository.findByTodoItemAndFile.mockResolvedValue(null);
    repository.update.mockResolvedValue(changed);
    await expect(service.update(10, 3, { fileId: 6 })).resolves.toEqual(changed);
    expect(repository.update).toHaveBeenCalledWith(3, 6);
  });

  it('deletes only the reference, not the file', async () => {
    repository.findByIdForTodoItem.mockResolvedValue(attachment());
    await service.delete(10, 3);
    expect(repository.delete).toHaveBeenCalledWith(3);
    expect(files.findById).not.toHaveBeenCalled();
  });
});
