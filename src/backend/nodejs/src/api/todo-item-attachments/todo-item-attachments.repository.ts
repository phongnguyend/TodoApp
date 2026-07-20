import { Injectable } from '@nestjs/common';
import { Prisma, TodoItemAttachment } from '@prisma/client';
import { PrismaService } from '../../shared/prisma/prisma.service';

@Injectable()
export class TodoItemAttachmentRepository {
  constructor(private readonly prisma: PrismaService) {}

  findByTodoItemId(todoItemId: number): Promise<TodoItemAttachment[]> {
    return this.prisma.todoItemAttachment.findMany({
      where: { todoItemId },
      orderBy: [{ createdAt: 'desc' }, { id: 'desc' }],
    });
  }

  findByIdForTodoItem(todoItemId: number, id: number): Promise<TodoItemAttachment | null> {
    return this.prisma.todoItemAttachment.findFirst({ where: { id, todoItemId } });
  }

  findByTodoItemAndFile(todoItemId: number, fileId: number): Promise<TodoItemAttachment | null> {
    return this.prisma.todoItemAttachment.findUnique({
      where: { todoItemId_fileId: { todoItemId, fileId } },
    });
  }

  create(data: Prisma.TodoItemAttachmentUncheckedCreateInput): Promise<TodoItemAttachment> {
    return this.prisma.todoItemAttachment.create({ data });
  }

  update(id: number, fileId: number, updatedByUserId?: number): Promise<TodoItemAttachment> {
    return this.prisma.todoItemAttachment.update({
      where: { id },
      data: { fileId, ...(updatedByUserId !== undefined ? { updatedByUserId } : {}) },
    });
  }

  async delete(id: number): Promise<void> {
    await this.prisma.todoItemAttachment.delete({ where: { id } });
  }
}
