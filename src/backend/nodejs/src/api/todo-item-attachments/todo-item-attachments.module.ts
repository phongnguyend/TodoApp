import { Module } from '@nestjs/common';
import { FileRepository } from '../files/files.repository';
import { TodoItemRepository } from '../todo-items/todo-items.repository';
import { TodoItemAttachmentsController } from './todo-item-attachments.controller';
import { TodoItemAttachmentRepository } from './todo-item-attachments.repository';
import { TodoItemAttachmentsService } from './todo-item-attachments.service';

@Module({
  controllers: [TodoItemAttachmentsController],
  providers: [
    TodoItemAttachmentsService,
    TodoItemAttachmentRepository,
    TodoItemRepository,
    FileRepository,
  ],
})
export class TodoItemAttachmentsModule {}
