import { BadRequestException, ConflictException, NotFoundException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Test } from '@nestjs/testing';
import { User } from '@prisma/client';
import { UsersRepository } from './users.repository';
import { UsersService } from './users.service';

const user = (overrides: Partial<User> = {}): User => ({
  id: 1, username: 'alice', email: 'alice@example.com', passwordHash: 'hash', isActive: true,
  createdAt: new Date('2026-01-01T00:00:00Z'), updatedAt: null, ...overrides,
});

describe('UsersService', () => {
  let service: UsersService;
  let repository: jest.Mocked<UsersRepository>;

  beforeEach(async () => {
    repository = {
      findAll: jest.fn(), findById: jest.fn(), findByEmail: jest.fn(), usernameExists: jest.fn(),
      emailExists: jest.fn(), create: jest.fn(), update: jest.fn(), createEmailLog: jest.fn(),
    } as unknown as jest.Mocked<UsersRepository>;
    const config = { get: jest.fn((key: string) => ({
      PASSWORD_HASH_ITERATIONS: '1', JWT_SECRET_KEY: 'test-secret',
      PASSWORD_RESET_TOKEN_LIFETIME_MINUTES: '60', PASSWORD_RESET_CONFIRMATION_URL: '/reset-password',
    })[key]) };
    const module = await Test.createTestingModule({ providers: [
      UsersService, { provide: UsersRepository, useValue: repository }, { provide: ConfigService, useValue: config },
    ] }).compile();
    service = module.get(UsersService);
  });

  it('lists users with bounded pagination and without password hashes', async () => {
    repository.findAll.mockResolvedValue({ items: [user()], total: 1 });
    const result = await service.getAll(0, 1000);
    expect(repository.findAll).toHaveBeenCalledWith(0, 100);
    expect(result).toMatchObject({ total: 1, page: 1, pageSize: 100, totalPages: 1 });
    expect(result.items[0]).not.toHaveProperty('passwordHash');
  });

  it('throws when a user is missing', async () => {
    repository.findById.mockResolvedValue(null);
    await expect(service.getById(99)).rejects.toThrow(NotFoundException);
  });

  it('normalizes input, hashes passwords, and detects conflicts', async () => {
    repository.create.mockImplementation(async (data) => user({
      username: data.username as string, email: data.email as string, passwordHash: data.passwordHash as string,
    }));
    const result = await service.create({ username: ' Alice ', email: 'ALICE@EXAMPLE.COM', password: 'password123' });
    expect(repository.create).toHaveBeenCalledWith(expect.objectContaining({ username: 'Alice', email: 'alice@example.com' }));
    expect((repository.create.mock.calls[0][0].passwordHash as string)).toMatch(/^pbkdf2_sha256\$/);
    expect(result).not.toHaveProperty('passwordHash');

    repository.usernameExists.mockResolvedValue(true);
    await expect(service.create({ username: 'Alice', email: 'other@example.com', password: 'password123' }))
      .rejects.toThrow(ConflictException);
  });

  it('activates and deactivates existing users', async () => {
    repository.findById.mockResolvedValue(user());
    repository.update.mockImplementation(async (_id, data) => user({ isActive: data.isActive as boolean }));
    await expect(service.setActive(1, false)).resolves.toMatchObject({ isActive: false });
  });

  it('changes a password only when the current password is valid', async () => {
    repository.create.mockImplementation(async (data) => user({ passwordHash: data.passwordHash as string }));
    await service.create({ username: 'alice', email: 'alice@example.com', password: 'old-password' });
    const stored = repository.create.mock.results[0].value;
    repository.findById.mockResolvedValue(await stored);
    repository.update.mockResolvedValue(user());
    await service.changePassword(1, { currentPassword: 'old-password', newPassword: 'new-password' });
    expect(repository.update).toHaveBeenCalledWith(1, { passwordHash: expect.stringMatching(/^pbkdf2_sha256\$/) });
    await expect(service.changePassword(1, { currentPassword: 'wrong', newPassword: 'new-password' }))
      .rejects.toThrow(BadRequestException);
  });

  it('does not reveal unknown reset accounts and creates reset email for active users', async () => {
    repository.findByEmail.mockResolvedValueOnce(null).mockResolvedValueOnce(user());
    await service.requestPasswordReset({ email: 'missing@example.com' });
    expect(repository.createEmailLog).not.toHaveBeenCalled();
    await service.requestPasswordReset({ email: 'ALICE@EXAMPLE.COM' });
    expect(repository.createEmailLog).toHaveBeenCalledWith(expect.objectContaining({
      recipient: 'alice@example.com', status: 'pending', body: expect.stringContaining('token='),
    }));
  });

  it('accepts a valid reset token once and rejects invalid tokens', async () => {
    repository.findByEmail.mockResolvedValue(user());
    await service.requestPasswordReset({ email: 'alice@example.com' });
    const body = repository.createEmailLog.mock.calls[0][0].body;
    const token = decodeURIComponent(body.match(/token=([^\s]+)/)?.[1] ?? '');
    repository.findById.mockResolvedValue(user());
    repository.update.mockResolvedValue(user());

    await service.confirmPasswordReset({ token, newPassword: 'new-password' });
    expect(repository.update).toHaveBeenCalledWith(1, {
      passwordHash: expect.stringMatching(/^pbkdf2_sha256\$/),
    });
    await expect(service.confirmPasswordReset({ token: 'invalid', newPassword: 'new-password' }))
      .rejects.toThrow(BadRequestException);
  });
});
