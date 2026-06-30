import { Module } from '@nestjs/common';
import { TodoItemsController } from './todo-items.controller';
import { TodoItemRepository } from './todo-items.repository';
import { TodoItemsService } from './todo-items.service';

/**
 * TodoItemsModule declares the controller, service, and repository for this feature.
 * Mirrors the pattern of grouping related services in an ASP.NET Core feature folder.
 */
@Module({
  controllers: [TodoItemsController],
  providers: [TodoItemsService, TodoItemRepository],
})
export class TodoItemsModule {}
