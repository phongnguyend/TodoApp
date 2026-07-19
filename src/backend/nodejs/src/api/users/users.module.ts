import { Module } from '@nestjs/common';
import { UsersController } from './users.controller';
import { UsersRepository } from './users.repository';
import { UserAuthGuard } from './users.security';
import { UsersService } from './users.service';
import { TokensController } from './tokens.controller';

@Module({
  controllers: [UsersController, TokensController],
  providers: [UsersService, UsersRepository, UserAuthGuard],
})
export class UsersModule {}
