import { Injectable, NotFoundException } from '@nestjs/common';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { parseCsv, toCsvRow } from './csv.util';
import { buildExcelWorkbook, parseExcel } from './excel.util';
import { CreateTodoItemDto } from './dto/create-todo-item.dto';
import { ImportResultDto, ImportRowErrorDto } from './dto/import-result.dto';
import { TodoItemResponseDto } from './dto/todo-item-response.dto';
import { UpdateTodoItemDto } from './dto/update-todo-item.dto';
import { TodoItemRepository } from './todo-items.repository';
import { TodoItem } from '@prisma/client';

const CSV_HEADER = ['id', 'title', 'description', 'is_completed', 'created_at', 'updated_at'];
const TRUE_VALUES = new Set(['1', 'true', 'yes', 'y']);

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
      createdByUserId: item.createdByUserId,
      updatedAt: item.updatedAt,
      updatedByUserId: item.updatedByUserId,
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

  private static parseBool(value: string | undefined): boolean {
    return TRUE_VALUES.has((value ?? '').trim().toLowerCase());
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

  async create(dto: CreateTodoItemDto, actorUserId?: number): Promise<TodoItemResponseDto> {
    const item = await this.repository.create({
      title: dto.title,
      description: dto.description,
      ...(actorUserId !== undefined ? { createdByUserId: actorUserId } : {}),
    });
    return TodoItemsService.toDto(item);
  }

  async update(id: number, dto: UpdateTodoItemDto, actorUserId?: number): Promise<TodoItemResponseDto> {
    await this.getOrThrow(id);
    const item = await this.repository.update(id, {
      title: dto.title,
      description: dto.description,
      isCompleted: dto.isCompleted,
      ...(actorUserId !== undefined ? { updatedByUserId: actorUserId } : {}),
    });
    return TodoItemsService.toDto(item);
  }

  async markComplete(id: number, actorUserId?: number): Promise<TodoItemResponseDto> {
    await this.getOrThrow(id);
    const item = await this.repository.update(id, {
      isCompleted: true,
      ...(actorUserId !== undefined ? { updatedByUserId: actorUserId } : {}),
    });
    return TodoItemsService.toDto(item);
  }

  async delete(id: number): Promise<void> {
    await this.getOrThrow(id);
    await this.repository.delete(id);
  }

  // ── CSV import/export ──────────────────────────────────────────────────────

  async importCsv(buffer: Buffer, actorUserId?: number): Promise<ImportResultDto> {
    const text = buffer.toString('utf-8').replace(/^\uFEFF/, ''); // strip UTF-8 BOM
    const rows = parseCsv(text);

    if (rows.length === 0) {
      return { imported: 0, failed: 0, errors: [] };
    }

    const header = rows[0].map((column) => column.trim().toLowerCase());
    const colIndex = new Map(header.map((name, idx) => [name, idx]));
    const getCell = (row: string[], name: string): string | undefined => {
      const idx = colIndex.get(name);
      return idx !== undefined ? row[idx] : undefined;
    };

    let imported = 0;
    const errors: ImportRowErrorDto[] = [];

    for (let i = 1; i < rows.length; i++) {
      const rowNumber = i + 1; // header occupies row 1
      const row = rows[i];

      const title = (getCell(row, 'title') ?? '').trim();
      if (!title) {
        errors.push({ row: rowNumber, error: 'Title is required.' });
        continue;
      }

      const description = (getCell(row, 'description') ?? '').trim() || null;
      const isCompleted = TodoItemsService.parseBool(getCell(row, 'is_completed'));

      await this.repository.create({ title, description, isCompleted,
        ...(actorUserId !== undefined ? { createdByUserId: actorUserId } : {}) });
      imported++;
    }

    return { imported, failed: errors.length, errors };
  }

  async exportCsv(): Promise<string> {
    const items = await this.repository.findAllOrdered();

    const lines = [toCsvRow(CSV_HEADER)];
    for (const item of items) {
      lines.push(
        toCsvRow([
          item.id,
          item.title,
          item.description ?? '',
          item.isCompleted,
          item.createdAt ? item.createdAt.toISOString() : '',
          item.updatedAt ? item.updatedAt.toISOString() : '',
        ]),
      );
    }
    return lines.join('');
  }

  // ── Excel import/export ─────────────────────────────────────────────────────

  async importExcel(buffer: Buffer, actorUserId?: number): Promise<ImportResultDto> {
    const rows = await parseExcel(buffer);

    if (rows.length === 0) {
      return { imported: 0, failed: 0, errors: [] };
    }

    const header = rows[0].map((column) => String(column ?? '').trim().toLowerCase());
    const colIndex = new Map(header.map((name, idx) => [name, idx]));
    const getCell = (row: unknown[], name: string): unknown => {
      const idx = colIndex.get(name);
      return idx !== undefined ? row[idx] : undefined;
    };

    let imported = 0;
    const errors: ImportRowErrorDto[] = [];

    for (let i = 1; i < rows.length; i++) {
      const rowNumber = i + 1; // header occupies row 1
      const row = rows[i];

      if (row === undefined || row.every((value) => value === null || value === undefined)) {
        continue;
      }

      const title = String(getCell(row, 'title') ?? '').trim();
      if (!title) {
        errors.push({ row: rowNumber, error: 'Title is required.' });
        continue;
      }

      const description = String(getCell(row, 'description') ?? '').trim() || null;
      const isCompleted = TodoItemsService.parseBool(String(getCell(row, 'is_completed') ?? ''));

      await this.repository.create({ title, description, isCompleted,
        ...(actorUserId !== undefined ? { createdByUserId: actorUserId } : {}) });
      imported++;
    }

    return { imported, failed: errors.length, errors };
  }

  async exportExcel(): Promise<Buffer> {
    const items = await this.repository.findAllOrdered();

    const rows = items.map((item) => [
      item.id,
      item.title,
      item.description ?? '',
      item.isCompleted,
      item.createdAt ? item.createdAt.toISOString() : '',
      item.updatedAt ? item.updatedAt.toISOString() : '',
    ]);

    return buildExcelWorkbook(CSV_HEADER, rows);
  }
}
