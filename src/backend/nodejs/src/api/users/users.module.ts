import { Module } from '@nestjs/common';
import { ConfigModule, ConfigService } from '@nestjs/config';
import { JwtModule } from '@nestjs/jwt';
import { PassportModule } from '@nestjs/passport';
import { UsersController } from './users.controller';
import { UsersRepository } from './users.repository';
import { UserAuthGuard } from './users.security';
import { UsersService } from './users.service';
import { TokensController } from './tokens.controller';
import { JwtStrategy } from './jwt.strategy';
import { APP_GUARD } from '@nestjs/core';

@Module({
  imports: [
    ConfigModule,
    PassportModule.register({ defaultStrategy: 'jwt' }),
    JwtModule.registerAsync({
      imports: [ConfigModule],
      inject: [ConfigService],
      useFactory: (config: ConfigService) => ({
        secret: config.get<string>('JWT_SECRET_KEY') ?? 'change-me-use-at-least-32-bytes-long',
        signOptions: { algorithm: 'HS256' },
      }),
    }),
  ],
  controllers: [UsersController, TokensController],
  providers: [
    UsersService, UsersRepository, JwtStrategy,
    { provide: APP_GUARD, useClass: UserAuthGuard },
  ],
})
export class UsersModule {}
