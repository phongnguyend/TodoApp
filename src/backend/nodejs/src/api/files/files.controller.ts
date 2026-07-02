import {
  BadRequestException,
  Controller,
  Delete,
  Get,
  HttpCode,
  HttpStatus,
  Param,
  ParseIntPipe,
  Post,
  Query,
  Res,
  UploadedFile,
  UseInterceptors,
} from '@nestjs/common';
import { FileInterceptor } from '@nestjs/platform-express';
import {
  ApiBody,
  ApiConsumes,
  ApiCreatedResponse,
  ApiNoContentResponse,
  ApiNotFoundResponse,
  ApiOkResponse,
  ApiPayloadTooLargeResponse,
  ApiQuery,
  ApiTags,
} from '@nestjs/swagger';
import { Response } from 'express';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { FileResponseDto } from './dto/file-response.dto';
import { FilesService } from './files.service';

/**
 * FilesController - mirrors a [ApiController] / [Route("api/files")] controller in ASP.NET Core.
 * NestJS decorators map directly to HTTP method + route attributes.
 */
@ApiTags('Files')
@Controller('api/files')
export class FilesController {
  constructor(private readonly service: FilesService) {}

  @Get()
  @ApiQuery({ name: 'page', required: false, type: Number })
  @ApiQuery({ name: 'pageSize', required: false, type: Number })
  @ApiOkResponse({ description: 'Paginated list of uploaded files' })
  getAll(
    @Query('page', new ParseIntPipe({ optional: true })) page = 1,
    @Query('pageSize', new ParseIntPipe({ optional: true })) pageSize = 20,
  ): Promise<PaginatedResponseDto<FileResponseDto>> {
    return this.service.getAll(page, pageSize);
  }

  @Get(':id')
  @ApiOkResponse({ type: FileResponseDto })
  @ApiNotFoundResponse({ description: 'File not found' })
  getById(@Param('id', ParseIntPipe) id: number): Promise<FileResponseDto> {
    return this.service.getById(id);
  }

  @Get(':id/download')
  @ApiOkResponse({ description: "The file's binary content" })
  @ApiNotFoundResponse({ description: 'File or its content not found' })
  async download(@Param('id', ParseIntPipe) id: number, @Res() res: Response): Promise<void> {
    const target = await this.service.getDownloadTarget(id);
    res.download(target.path, target.name, {
      headers: { 'Content-Type': target.contentType },
    });
  }

  @Post()
  @UseInterceptors(FileInterceptor('file'))
  @HttpCode(HttpStatus.CREATED)
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    schema: {
      type: 'object',
      properties: { file: { type: 'string', format: 'binary' } },
    },
  })
  @ApiCreatedResponse({ type: FileResponseDto })
  @ApiPayloadTooLargeResponse({ description: 'File exceeds the maximum allowed size' })
  create(@UploadedFile() file?: Express.Multer.File): Promise<FileResponseDto> {
    if (!file) {
      throw new BadRequestException('file is required.');
    }
    return this.service.upload(file);
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiNoContentResponse({ description: 'File deleted' })
  @ApiNotFoundResponse({ description: 'File not found' })
  async delete(@Param('id', ParseIntPipe) id: number): Promise<void> {
    return this.service.delete(id);
  }
}
