import { UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { ExecutionContext } from '@nestjs/common';
import { createHmac } from 'crypto';
import { UserAuthGuard } from './users.security';

const tokenFor = (payload: object, secret = 'test-secret'): string => {
  const header = Buffer.from(JSON.stringify({ alg: 'HS256', typ: 'JWT' })).toString('base64url');
  const body = Buffer.from(JSON.stringify(payload)).toString('base64url');
  const signature = createHmac('sha256', secret).update(`${header}.${body}`).digest('base64url');
  return `${header}.${body}.${signature}`;
};

describe('UserAuthGuard', () => {
  const guard = new UserAuthGuard({ get: jest.fn(() => 'test-secret') } as unknown as ConfigService);
  const contextFor = (authorization?: string): { context: ExecutionContext; request: { headers: { authorization?: string }; userId?: number } } => {
    const request = { headers: { authorization } };
    const context = { switchToHttp: () => ({ getRequest: () => request }) } as unknown as ExecutionContext;
    return { context, request };
  };

  it('accepts a valid HS256 bearer token and exposes its subject', () => {
    const { context, request } = contextFor(`Bearer ${tokenFor({ sub: '42', exp: Math.floor(Date.now() / 1000) + 60 })}`);
    expect(guard.canActivate(context)).toBe(true);
    expect(request.userId).toBe(42);
  });

  it('rejects missing, expired, and incorrectly signed tokens', () => {
    expect(() => guard.canActivate(contextFor().context)).toThrow(UnauthorizedException);
    expect(() => guard.canActivate(contextFor(`Bearer ${tokenFor({ sub: 1, exp: 1 })}`).context)).toThrow(UnauthorizedException);
    expect(() => guard.canActivate(contextFor(`Bearer ${tokenFor({ sub: 1 }, 'wrong')}`).context)).toThrow(UnauthorizedException);
  });
});
