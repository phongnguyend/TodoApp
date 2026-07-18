import { Body, Controller, Delete, Get, HttpCode, HttpStatus, Param, ParseIntPipe, Post, Put } from '@nestjs/common';
import { ApiCreatedResponse, ApiNoContentResponse, ApiNotFoundResponse, ApiOkResponse, ApiTags } from '@nestjs/swagger';
import { SaveTodoItemAttachmentDto } from './dto/save-todo-item-attachment.dto';
import { TodoItemAttachmentResponseDto } from './dto/todo-item-attachment-response.dto';
import { TodoItemAttachmentsService } from './todo-item-attachments.service';

@ApiTags('Todo Item Attachments')
@Controller('api/todo-items/:id/attachments')
export class TodoItemAttachmentsController {
  constructor(private readonly service: TodoItemAttachmentsService) {}

  @Get()
  @ApiOkResponse({ type: [TodoItemAttachmentResponseDto] })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  getAll(@Param('id', ParseIntPipe) id: number): Promise<TodoItemAttachmentResponseDto[]> {
    return this.service.getAll(id);
  }

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiCreatedResponse({ type: TodoItemAttachmentResponseDto })
  @ApiNotFoundResponse({ description: 'Todo item or file not found' })
  create(
    @Param('id', ParseIntPipe) id: number,
    @Body() dto: SaveTodoItemAttachmentDto,
  ): Promise<TodoItemAttachmentResponseDto> {
    return this.service.create(id, dto);
  }

  @Get(':attachmentId')
  @ApiOkResponse({ type: TodoItemAttachmentResponseDto })
  @ApiNotFoundResponse({ description: 'Todo item or attachment not found' })
  getById(
    @Param('id', ParseIntPipe) id: number,
    @Param('attachmentId', ParseIntPipe) attachmentId: number,
  ): Promise<TodoItemAttachmentResponseDto> {
    return this.service.getById(id, attachmentId);
  }

  @Put(':attachmentId')
  @ApiOkResponse({ type: TodoItemAttachmentResponseDto })
  @ApiNotFoundResponse({ description: 'Todo item, file, or attachment not found' })
  update(
    @Param('id', ParseIntPipe) id: number,
    @Param('attachmentId', ParseIntPipe) attachmentId: number,
    @Body() dto: SaveTodoItemAttachmentDto,
  ): Promise<TodoItemAttachmentResponseDto> {
    return this.service.update(id, attachmentId, dto);
  }

  @Delete(':attachmentId')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiNoContentResponse({ description: 'Attachment reference removed' })
  @ApiNotFoundResponse({ description: 'Todo item or attachment not found' })
  delete(
    @Param('id', ParseIntPipe) id: number,
    @Param('attachmentId', ParseIntPipe) attachmentId: number,
  ): Promise<void> {
    return this.service.delete(id, attachmentId);
  }
}
