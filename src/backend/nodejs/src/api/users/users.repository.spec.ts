import { Test } from '@nestjs/testing';
import { PrismaService } from '../../shared/prisma/prisma.service';
import { UsersRepository } from './users.repository';

describe('UsersRepository', () => {
  let repository: UsersRepository;
  const prisma = {
    user: { findMany: jest.fn(), count: jest.fn(), findUnique: jest.fn(), create: jest.fn(), update: jest.fn() },
    emailLog: { create: jest.fn() },
    $transaction: jest.fn(),
  };

  beforeEach(async () => {
    jest.clearAllMocks();
    const module = await Test.createTestingModule({
      providers: [UsersRepository, { provide: PrismaService, useValue: prisma }],
    }).compile();
    repository = module.get(UsersRepository);
  });

  it('returns paginated users', async () => {
    prisma.user.findMany.mockReturnValue('users-query');
    prisma.user.count.mockReturnValue('count-query');
    prisma.$transaction.mockResolvedValue([[{ id: 1 }], 1]);

    await expect(repository.findAll(10, 5)).resolves.toEqual({ items: [{ id: 1 }], total: 1 });
    expect(prisma.user.findMany).toHaveBeenCalledWith({ skip: 10, take: 5, orderBy: { createdAt: 'desc' } });
  });

  it('excludes the current user from uniqueness checks', async () => {
    prisma.user.count.mockResolvedValue(0);
    await repository.usernameExists('alice', 7);
    await repository.emailExists('alice@example.com', 7);
    expect(prisma.user.count).toHaveBeenNthCalledWith(1, { where: { username: { equals: 'alice' }, id: { not: 7 } } });
    expect(prisma.user.count).toHaveBeenNthCalledWith(2, { where: { email: { equals: 'alice@example.com' }, id: { not: 7 } } });
  });

  it('writes users, updates, and password-reset email logs', async () => {
    prisma.user.create.mockResolvedValue({ id: 1 });
    prisma.user.update.mockResolvedValue({ id: 1 });
    prisma.emailLog.create.mockResolvedValue({ id: 2 });
    await repository.create({ username: 'alice', email: 'a@b.com', passwordHash: 'hash' });
    await repository.update(1, { isActive: false });
    await repository.createEmailLog({ recipient: 'a@b.com', subject: 'Reset', body: 'body' });
    expect(prisma.user.create).toHaveBeenCalled();
    expect(prisma.user.update).toHaveBeenCalledWith({ where: { id: 1 }, data: { isActive: false } });
    expect(prisma.emailLog.create).toHaveBeenCalled();
  });
});
