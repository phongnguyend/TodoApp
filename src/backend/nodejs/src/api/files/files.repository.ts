import { Injectable } from '@nestjs/common';
import { File, Prisma } from '@prisma/client';
import { PrismaService } from '../../shared/prisma/prisma.service';

export interface PaginatedResult<T> {
  items: T[];
  total: number;
}

/**
 * FileRepository encapsulates all Prisma data-access logic for uploaded files.
 * Mirrors a repository implementation backed by EF Core's DbSet<File>.
 */
@Injectable()
export class FileRepository {
  constructor(private readonly prisma: PrismaService) {}

  async findAll(skip: number, take: number): Promise<PaginatedResult<File>> {
    const [items, total] = await this.prisma.$transaction([
      this.prisma.file.findMany({ skip, take, orderBy: { createdAt: 'desc' } }),
      this.prisma.file.count(),
    ]);
    return { items, total };
  }

  async findById(id: number): Promise<File | null> {
    return this.prisma.file.findUnique({ where: { id } });
  }

  async create(data: Prisma.FileCreateInput): Promise<File> {
    return this.prisma.file.create({ data });
  }

  async delete(id: number): Promise<void> {
    await this.prisma.file.delete({ where: { id } });
  }
}
