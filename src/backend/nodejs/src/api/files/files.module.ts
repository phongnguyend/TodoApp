import { Module } from '@nestjs/common';
import { FilesController } from './files.controller';
import { FileRepository } from './files.repository';
import { FilesService } from './files.service';

/**
 * FilesModule declares the controller, service, and repository for this feature.
 * Mirrors the pattern of grouping related services in an ASP.NET Core feature folder.
 */
@Module({
  controllers: [FilesController],
  providers: [FilesService, FileRepository],
})
export class FilesModule {}
