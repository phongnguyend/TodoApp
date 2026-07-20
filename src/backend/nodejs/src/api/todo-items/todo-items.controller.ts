import {
  BadRequestException,
  Body,
  Controller,
  Delete,
  Get,
  HttpCode,
  HttpStatus,
  Param,
  ParseIntPipe,
  Patch,
  Post,
  Put,
  Query,
  Req,
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
  ApiQuery,
  ApiTags,
} from '@nestjs/swagger';
import { Response } from 'express';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { CreateTodoItemDto } from './dto/create-todo-item.dto';
import { ImportResultDto } from './dto/import-result.dto';
import { TodoItemResponseDto } from './dto/todo-item-response.dto';
import { UpdateTodoItemDto } from './dto/update-todo-item.dto';
import { TodoItemsService } from './todo-items.service';
import { AuditRequest } from '../users/users.security';

/**
 * TodoItemsController - mirrors a [ApiController] / [Route("api/todo-items")] controller in ASP.NET Core.
 * NestJS decorators map directly to HTTP method + route attributes.
 */
@ApiTags('Todo Items')
@Controller('api/todo-items')
export class TodoItemsController {
  constructor(private readonly service: TodoItemsService) {}

  @Get()
  @ApiQuery({ name: 'page', required: false, type: Number })
  @ApiQuery({ name: 'pageSize', required: false, type: Number })
  @ApiOkResponse({ description: 'Paginated list of all todo items' })
  getAll(
    @Query('page', new ParseIntPipe({ optional: true })) page = 1,
    @Query('pageSize', new ParseIntPipe({ optional: true })) pageSize = 20,
  ): Promise<PaginatedResponseDto<TodoItemResponseDto>> {
    return this.service.getAll(page, pageSize);
  }

  @Get('incomplete')
  @ApiQuery({ name: 'page', required: false, type: Number })
  @ApiQuery({ name: 'pageSize', required: false, type: Number })
  @ApiOkResponse({ description: 'Paginated list of incomplete todo items' })
  getIncomplete(
    @Query('page', new ParseIntPipe({ optional: true })) page = 1,
    @Query('pageSize', new ParseIntPipe({ optional: true })) pageSize = 20,
  ): Promise<PaginatedResponseDto<TodoItemResponseDto>> {
    return this.service.getIncomplete(page, pageSize);
  }

  @Get(':id')
  @ApiOkResponse({ type: TodoItemResponseDto })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  getById(@Param('id', ParseIntPipe) id: number): Promise<TodoItemResponseDto> {
    return this.service.getById(id);
  }

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiCreatedResponse({ type: TodoItemResponseDto })
  create(@Body() dto: CreateTodoItemDto, @Req() request?: AuditRequest): Promise<TodoItemResponseDto> {
    const actor = request?.user?.userId;
    return actor === undefined ? this.service.create(dto) : this.service.create(dto, actor);
  }

  @Put(':id')
  @ApiOkResponse({ type: TodoItemResponseDto })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  update(
    @Param('id', ParseIntPipe) id: number,
    @Body() dto: UpdateTodoItemDto,
    @Req() request?: AuditRequest,
  ): Promise<TodoItemResponseDto> {
    const actor = request?.user?.userId;
    return actor === undefined ? this.service.update(id, dto) : this.service.update(id, dto, actor);
  }

  @Patch(':id/complete')
  @ApiOkResponse({ type: TodoItemResponseDto, description: 'Mark todo item as complete' })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  markComplete(@Param('id', ParseIntPipe) id: number, @Req() request?: AuditRequest): Promise<TodoItemResponseDto> {
    const actor = request?.user?.userId;
    return actor === undefined ? this.service.markComplete(id) : this.service.markComplete(id, actor);
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiNoContentResponse({ description: 'Todo item deleted' })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  async delete(@Param('id', ParseIntPipe) id: number): Promise<void> {
    return this.service.delete(id);
  }

  @Post('import/csv')
  @UseInterceptors(FileInterceptor('file'))
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    schema: {
      type: 'object',
      properties: { file: { type: 'string', format: 'binary' } },
    },
  })
  @ApiOkResponse({ type: ImportResultDto, description: 'Import result summary' })
  importCsv(@UploadedFile() file?: Express.Multer.File, @Req() request?: AuditRequest): Promise<ImportResultDto> {
    if (!file) {
      throw new BadRequestException('file is required.');
    }
    const actor = request?.user?.userId;
    return actor === undefined ? this.service.importCsv(file.buffer) : this.service.importCsv(file.buffer, actor);
  }

  @Get('export/csv')
  @ApiOkResponse({ description: 'CSV file containing all todo items' })
  async exportCsv(@Res() res: Response): Promise<void> {
    const content = await this.service.exportCsv();
    res.set({
      'Content-Type': 'text/csv',
      'Content-Disposition': 'attachment; filename="todo_items.csv"',
    });
    res.send(content);
  }

  @Post('import/excel')
  @UseInterceptors(FileInterceptor('file'))
  @ApiConsumes('multipart/form-data')
  @ApiBody({
    schema: {
      type: 'object',
      properties: { file: { type: 'string', format: 'binary' } },
    },
  })
  @ApiOkResponse({ type: ImportResultDto, description: 'Import result summary' })
  importExcel(@UploadedFile() file?: Express.Multer.File, @Req() request?: AuditRequest): Promise<ImportResultDto> {
    if (!file) {
      throw new BadRequestException('file is required.');
    }
    const actor = request?.user?.userId;
    return actor === undefined ? this.service.importExcel(file.buffer) : this.service.importExcel(file.buffer, actor);
  }

  @Get('export/excel')
  @ApiOkResponse({ description: 'Excel file containing all todo items' })
  async exportExcel(@Res() res: Response): Promise<void> {
    const content = await this.service.exportExcel();
    res.set({
      'Content-Type': 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet',
      'Content-Disposition': 'attachment; filename="todo_items.xlsx"',
    });
    res.send(content);
  }
}
