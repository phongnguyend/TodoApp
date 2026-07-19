import { Test } from '@nestjs/testing';
import { ConfigService } from '@nestjs/config';
import { UsersController } from './users.controller';
import { UsersService } from './users.service';

describe('UsersController', () => {
  let controller: UsersController;
  let service: jest.Mocked<UsersService>;

  beforeEach(async () => {
    service = {
      getAll: jest.fn(), getById: jest.fn(), create: jest.fn(), update: jest.fn(), setActive: jest.fn(),
      signup: jest.fn(), getProfile: jest.fn(), updateProfile: jest.fn(), changePassword: jest.fn(),
      requestPasswordReset: jest.fn(), confirmPasswordReset: jest.fn(),
    } as unknown as jest.Mocked<UsersService>;
    const module = await Test.createTestingModule({
      controllers: [UsersController], providers: [
        { provide: UsersService, useValue: service },
        { provide: ConfigService, useValue: { get: jest.fn() } },
      ],
    }).compile();
    controller = module.get(UsersController);
  });

  it('delegates user-management endpoints', async () => {
    await controller.getAll(2, 5);
    await controller.getById(3);
    await controller.create({ username: 'a', email: 'a@b.com', password: 'password' });
    await controller.update(3, { username: 'b' });
    await controller.activate(3);
    await controller.deactivate(3);
    expect(service.getAll).toHaveBeenCalledWith(2, 5);
    expect(service.getById).toHaveBeenCalledWith(3);
    expect(service.setActive).toHaveBeenNthCalledWith(1, 3, true);
    expect(service.setActive).toHaveBeenNthCalledWith(2, 3, false);
  });

  it('uses the authenticated request user for self-management', async () => {
    const request = { user: { userId: 7 } } as never;
    await controller.getProfile(request);
    await controller.updateProfile(request, { username: 'new-name' });
    await controller.changePassword(request, { currentPassword: 'old', newPassword: 'new-password' });
    expect(service.getProfile).toHaveBeenCalledWith(7);
    expect(service.updateProfile).toHaveBeenCalledWith(7, { username: 'new-name' });
    expect(service.changePassword).toHaveBeenCalledWith(7, expect.anything());
  });

  it('returns the enumeration-safe reset response', async () => {
    const response = await controller.requestPasswordReset({ email: 'missing@example.com' });
    expect(response.message).toContain('If the account exists');
    expect(service.requestPasswordReset).toHaveBeenCalled();
  });
});
