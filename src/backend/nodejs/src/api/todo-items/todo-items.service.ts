import { Injectable, NotFoundException } from '@nestjs/common';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { CreateTodoItemDto } from './dto/create-todo-item.dto';
import { TodoItemResponseDto } from './dto/todo-item-response.dto';
import { UpdateTodoItemDto } from './dto/update-todo-item.dto';
import { TodoItemRepository } from './todo-items.repository';
import { TodoItem } from '@prisma/client';

/**
 * TodoItemsService contains all business logic.
 * Mirrors a service class registered via DI in ASP.NET Core.
 */
@Injectable()
export class TodoItemsService {
  constructor(private readonly repository: TodoItemRepository) {}

  // ── Helpers ──────────────────────────────────────────────────────────────────

  private static toDto(item: TodoItem): TodoItemResponseDto {
    return {
      id: item.id,
      title: item.title,
      description: item.description,
      isCompleted: item.isCompleted,
      createdAt: item.createdAt,
      updatedAt: item.updatedAt,
    };
  }

  private static toPaginated(
    items: TodoItem[],
    total: number,
    page: number,
    pageSize: number,
  ): PaginatedResponseDto<TodoItemResponseDto> {
    return {
      items: items.map(TodoItemsService.toDto),
      total,
      page,
      pageSize,
      totalPages: Math.ceil(total / pageSize),
    };
  }

  private async getOrThrow(id: number): Promise<TodoItem> {
    const item = await this.repository.findById(id);
    if (!item) {
      throw new NotFoundException(`Todo item ${id} not found.`);
    }
    return item;
  }

  // ── Queries ───────────────────────────────────────────────────────────────────

  async getAll(page: number, pageSize: number): Promise<PaginatedResponseDto<TodoItemResponseDto>> {
    const skip = (page - 1) * pageSize;
    const { items, total } = await this.repository.findAll(skip, pageSize);
    return TodoItemsService.toPaginated(items, total, page, pageSize);
  }

  async getIncomplete(page: number, pageSize: number): Promise<PaginatedResponseDto<TodoItemResponseDto>> {
    const skip = (page - 1) * pageSize;
    const { items, total } = await this.repository.findIncomplete(skip, pageSize);
    return TodoItemsService.toPaginated(items, total, page, pageSize);
  }

  async getById(id: number): Promise<TodoItemResponseDto> {
    const item = await this.getOrThrow(id);
    return TodoItemsService.toDto(item);
  }

  // ── Commands ──────────────────────────────────────────────────────────────────

  async create(dto: CreateTodoItemDto): Promise<TodoItemResponseDto> {
    const item = await this.repository.create({
      title: dto.title,
      description: dto.description,
    });
    return TodoItemsService.toDto(item);
  }

  async update(id: number, dto: UpdateTodoItemDto): Promise<TodoItemResponseDto> {
    await this.getOrThrow(id);
    const item = await this.repository.update(id, {
      title: dto.title,
      description: dto.description,
      isCompleted: dto.isCompleted,
    });
    return TodoItemsService.toDto(item);
  }

  async markComplete(id: number): Promise<TodoItemResponseDto> {
    await this.getOrThrow(id);
    const item = await this.repository.update(id, { isCompleted: true });
    return TodoItemsService.toDto(item);
  }

  async delete(id: number): Promise<void> {
    await this.getOrThrow(id);
    await this.repository.delete(id);
  }
}
