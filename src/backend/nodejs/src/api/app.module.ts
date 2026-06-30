import { Module } from '@nestjs/common';
import { ConfigModule } from '@nestjs/config';
import { PrismaModule } from '../shared/prisma/prisma.module';
import { TodoItemsModule } from './todo-items/todo-items.module';

/**
 * AppModule is the root module — analogous to Startup.cs / Program.cs in ASP.NET Core.
 * It wires together configuration, the Prisma (DbContext) module, and feature modules.
 */
@Module({
  imports: [
    ConfigModule.forRoot({ isGlobal: true }),
    PrismaModule,
    TodoItemsModule,
  ],
})
export class AppModule {}
