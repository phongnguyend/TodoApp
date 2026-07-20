import { Injectable, NotFoundException, PayloadTooLargeException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { File } from '@prisma/client';
import { randomUUID } from 'crypto';
import { promises as fs } from 'fs';
import * as path from 'path';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { FileResponseDto } from './dto/file-response.dto';
import { FileRepository } from './files.repository';

export interface DownloadTarget {
  path: string;
  name: string;
  contentType: string;
}

/**
 * FilesService contains all business logic for uploading, listing, downloading,
 * and deleting files. Mirrors a service class registered via DI in ASP.NET Core.
 */
@Injectable()
export class FilesService {
  private readonly storageDir: string;
  private readonly maxUploadSizeBytes: number;

  constructor(
    private readonly repository: FileRepository,
    private readonly config: ConfigService,
  ) {
    this.storageDir = this.config.get<string>('FILE_STORAGE_PATH') ?? './uploads';
    this.maxUploadSizeBytes = Number(this.config.get<string>('MAX_UPLOAD_SIZE_BYTES') ?? 10 * 1024 * 1024);
  }

  // ── Helpers ──────────────────────────────────────────────────────────────────

  private static toDto(file: File): FileResponseDto {
    return {
      id: file.id,
      name: file.name,
      extension: file.extension,
      size: file.size,
      contentType: file.contentType,
      createdAt: file.createdAt,
      createdByUserId: file.createdByUserId,
      updatedAt: file.updatedAt,
      updatedByUserId: file.updatedByUserId,
    };
  }

  private static toPaginated(
    items: File[],
    total: number,
    page: number,
    pageSize: number,
  ): PaginatedResponseDto<FileResponseDto> {
    return {
      items: items.map(FilesService.toDto),
      total,
      page,
      pageSize,
      totalPages: Math.ceil(total / pageSize),
    };
  }

  private async getOrThrow(id: number): Promise<File> {
    const file = await this.repository.findById(id);
    if (!file) {
      throw new NotFoundException(`File ${id} not found.`);
    }
    return file;
  }

  // ── Queries ───────────────────────────────────────────────────────────────────

  async getAll(page: number, pageSize: number): Promise<PaginatedResponseDto<FileResponseDto>> {
    const skip = (page - 1) * pageSize;
    const { items, total } = await this.repository.findAll(skip, pageSize);
    return FilesService.toPaginated(items, total, page, pageSize);
  }

  async getById(id: number): Promise<FileResponseDto> {
    const file = await this.getOrThrow(id);
    return FilesService.toDto(file);
  }

  async getDownloadTarget(id: number): Promise<DownloadTarget> {
    const file = await this.getOrThrow(id);
    try {
      await fs.access(file.location);
    } catch {
      throw new NotFoundException(`File ${id} content not found on disk.`);
    }
    return {
      path: file.location,
      name: file.name,
      contentType: file.contentType ?? 'application/octet-stream',
    };
  }

  // ── Commands ──────────────────────────────────────────────────────────────────

  async upload(uploadedFile: Express.Multer.File, actorUserId?: number): Promise<FileResponseDto> {
    if (uploadedFile.size > this.maxUploadSizeBytes) {
      throw new PayloadTooLargeException(
        `File exceeds the maximum allowed size of ${this.maxUploadSizeBytes} bytes.`,
      );
    }

    // Strip any directory components from the client-supplied name to prevent path traversal.
    const originalName = path.basename(uploadedFile.originalname);
    const extension = path.extname(originalName).replace(/^\./, '').toLowerCase();

    await fs.mkdir(this.storageDir, { recursive: true });

    // A random prefix avoids collisions/overwrites between uploads that share a name.
    const storedName = `${randomUUID()}_${originalName}`;
    const location = path.join(this.storageDir, storedName);
    await fs.writeFile(location, uploadedFile.buffer);

    const file = await this.repository.create({
      name: originalName,
      extension,
      size: uploadedFile.size,
      contentType: uploadedFile.mimetype,
      location,
      ...(actorUserId !== undefined ? { createdByUserId: actorUserId } : {}),
    });
    return FilesService.toDto(file);
  }

  async delete(id: number): Promise<void> {
    const file = await this.getOrThrow(id);
    await this.repository.delete(id);
    try {
      await fs.unlink(file.location);
    } catch {
      // Content already missing on disk - nothing left to clean up.
    }
  }
}
