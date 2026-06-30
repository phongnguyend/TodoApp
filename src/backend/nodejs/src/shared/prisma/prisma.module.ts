import { Global, Module } from '@nestjs/common';
import { PrismaService } from './prisma.service';

/**
 * @Global makes PrismaService available across all modules without
 * explicit imports — similar to registering DbContext in ASP.NET Core's
 * service collection with a scoped or singleton lifetime.
 */
@Global()
@Module({
  providers: [PrismaService],
  exports: [PrismaService],
})
export class PrismaModule {}
