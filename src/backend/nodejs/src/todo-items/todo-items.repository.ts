import { Injectable } from '@nestjs/common';
import { Prisma, TodoItem } from '@prisma/client';
import { PrismaService } from '../../prisma/prisma.service';

export interface PaginatedResult<T> {
  items: T[];
  total: number;
}

/**
 * TodoItemRepository encapsulates all Prisma data-access logic.
 * Mirrors a repository implementation backed by EF Core's DbSet<TodoItem>.
 */
@Injectable()
export class TodoItemRepository {
  constructor(private readonly prisma: PrismaService) {}

  async findAll(skip: number, take: number): Promise<PaginatedResult<TodoItem>> {
    const [items, total] = await this.prisma.$transaction([
      this.prisma.todoItem.findMany({ skip, take, orderBy: { createdAt: 'desc' } }),
      this.prisma.todoItem.count(),
    ]);
    return { items, total };
  }

  async findIncomplete(skip: number, take: number): Promise<PaginatedResult<TodoItem>> {
    const where: Prisma.TodoItemWhereInput = { isCompleted: false };
    const [items, total] = await this.prisma.$transaction([
      this.prisma.todoItem.findMany({ where, skip, take, orderBy: { createdAt: 'desc' } }),
      this.prisma.todoItem.count({ where }),
    ]);
    return { items, total };
  }

  async findById(id: number): Promise<TodoItem | null> {
    return this.prisma.todoItem.findUnique({ where: { id } });
  }

  async create(data: Prisma.TodoItemCreateInput): Promise<TodoItem> {
    return this.prisma.todoItem.create({ data });
  }

  async update(id: number, data: Prisma.TodoItemUpdateInput): Promise<TodoItem> {
    return this.prisma.todoItem.update({ where: { id }, data });
  }

  async delete(id: number): Promise<void> {
    await this.prisma.todoItem.delete({ where: { id } });
  }
}
