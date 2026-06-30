import {
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
} from '@nestjs/common';
import {
  ApiCreatedResponse,
  ApiNoContentResponse,
  ApiNotFoundResponse,
  ApiOkResponse,
  ApiQuery,
  ApiTags,
} from '@nestjs/swagger';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { CreateTodoItemDto } from './dto/create-todo-item.dto';
import { TodoItemResponseDto } from './dto/todo-item-response.dto';
import { UpdateTodoItemDto } from './dto/update-todo-item.dto';
import { TodoItemsService } from './todo-items.service';

/**
 * TodoItemsController — mirrors a [ApiController] / [Route("api/todo-items")] controller in ASP.NET Core.
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
  create(@Body() dto: CreateTodoItemDto): Promise<TodoItemResponseDto> {
    return this.service.create(dto);
  }

  @Put(':id')
  @ApiOkResponse({ type: TodoItemResponseDto })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  update(
    @Param('id', ParseIntPipe) id: number,
    @Body() dto: UpdateTodoItemDto,
  ): Promise<TodoItemResponseDto> {
    return this.service.update(id, dto);
  }

  @Patch(':id/complete')
  @ApiOkResponse({ type: TodoItemResponseDto, description: 'Mark todo item as complete' })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  markComplete(@Param('id', ParseIntPipe) id: number): Promise<TodoItemResponseDto> {
    return this.service.markComplete(id);
  }

  @Delete(':id')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiNoContentResponse({ description: 'Todo item deleted' })
  @ApiNotFoundResponse({ description: 'Todo item not found' })
  async delete(@Param('id', ParseIntPipe) id: number): Promise<void> {
    return this.service.delete(id);
  }
}
