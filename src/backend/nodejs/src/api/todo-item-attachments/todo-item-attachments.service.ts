import { Injectable, NotFoundException } from '@nestjs/common';
import { TodoItemAttachment } from '@prisma/client';
import { FileRepository } from '../files/files.repository';
import { TodoItemRepository } from '../todo-items/todo-items.repository';
import { SaveTodoItemAttachmentDto } from './dto/save-todo-item-attachment.dto';
import { TodoItemAttachmentResponseDto } from './dto/todo-item-attachment-response.dto';
import { TodoItemAttachmentRepository } from './todo-item-attachments.repository';

@Injectable()
export class TodoItemAttachmentsService {
  constructor(
    private readonly repository: TodoItemAttachmentRepository,
    private readonly todoItems: TodoItemRepository,
    private readonly files: FileRepository,
  ) {}

  private static toDto(item: TodoItemAttachment): TodoItemAttachmentResponseDto {
    return { ...item };
  }

  private async requireTodo(todoItemId: number): Promise<void> {
    if (!(await this.todoItems.findById(todoItemId))) {
      throw new NotFoundException(`Todo item ${todoItemId} not found.`);
    }
  }

  private async requireFile(fileId: number): Promise<void> {
    if (!(await this.files.findById(fileId))) {
      throw new NotFoundException(`File ${fileId} not found.`);
    }
  }

  private async requireAttachment(todoItemId: number, attachmentId: number): Promise<TodoItemAttachment> {
    const attachment = await this.repository.findByIdForTodoItem(todoItemId, attachmentId);
    if (!attachment) {
      throw new NotFoundException(
        `Attachment ${attachmentId} not found for todo item ${todoItemId}.`,
      );
    }
    return attachment;
  }

  async getAll(todoItemId: number): Promise<TodoItemAttachmentResponseDto[]> {
    await this.requireTodo(todoItemId);
    return (await this.repository.findByTodoItemId(todoItemId)).map(TodoItemAttachmentsService.toDto);
  }

  async getById(todoItemId: number, attachmentId: number): Promise<TodoItemAttachmentResponseDto> {
    await this.requireTodo(todoItemId);
    return TodoItemAttachmentsService.toDto(await this.requireAttachment(todoItemId, attachmentId));
  }

  async create(todoItemId: number, dto: SaveTodoItemAttachmentDto): Promise<TodoItemAttachmentResponseDto> {
    await this.requireTodo(todoItemId);
    await this.requireFile(dto.fileId);
    const existing = await this.repository.findByTodoItemAndFile(todoItemId, dto.fileId);
    const attachment = existing ?? await this.repository.create({ todoItemId, fileId: dto.fileId });
    return TodoItemAttachmentsService.toDto(attachment);
  }

  async update(todoItemId: number, attachmentId: number, dto: SaveTodoItemAttachmentDto): Promise<TodoItemAttachmentResponseDto> {
    await this.requireTodo(todoItemId);
    await this.requireFile(dto.fileId);
    const current = await this.requireAttachment(todoItemId, attachmentId);
    const existing = await this.repository.findByTodoItemAndFile(todoItemId, dto.fileId);
    if (existing && existing.id !== current.id) {
      return TodoItemAttachmentsService.toDto(existing);
    }
    return TodoItemAttachmentsService.toDto(await this.repository.update(current.id, dto.fileId));
  }

  async delete(todoItemId: number, attachmentId: number): Promise<void> {
    await this.requireTodo(todoItemId);
    const attachment = await this.requireAttachment(todoItemId, attachmentId);
    await this.repository.delete(attachment.id);
  }
}
