import { Injectable } from '@nestjs/common';
import { EmailLog, Prisma, User } from '@prisma/client';
import { PrismaService } from '../../shared/prisma/prisma.service';

@Injectable()
export class UsersRepository {
  constructor(private readonly prisma: PrismaService) {}

  async findAll(skip: number, take: number): Promise<{ items: User[]; total: number }> {
    const [items, total] = await this.prisma.$transaction([
      this.prisma.user.findMany({ skip, take, orderBy: { createdAt: 'desc' } }),
      this.prisma.user.count(),
    ]);
    return { items, total };
  }

  findById(id: number): Promise<User | null> {
    return this.prisma.user.findUnique({ where: { id } });
  }

  findByEmail(email: string): Promise<User | null> {
    return this.prisma.user.findUnique({ where: { email } });
  }

  async usernameExists(username: string, excludingId?: number): Promise<boolean> {
    return (await this.prisma.user.count({
      where: { username: { equals: username }, ...(excludingId ? { id: { not: excludingId } } : {}) },
    })) > 0;
  }

  async emailExists(email: string, excludingId?: number): Promise<boolean> {
    return (await this.prisma.user.count({
      where: { email: { equals: email }, ...(excludingId ? { id: { not: excludingId } } : {}) },
    })) > 0;
  }

  create(data: Prisma.UserUncheckedCreateInput): Promise<User> {
    return this.prisma.user.create({ data });
  }

  update(id: number, data: Prisma.UserUncheckedUpdateInput): Promise<User> {
    return this.prisma.user.update({ where: { id }, data });
  }

  createEmailLog(data: Prisma.EmailLogUncheckedCreateInput): Promise<EmailLog> {
    return this.prisma.emailLog.create({ data });
  }
}
