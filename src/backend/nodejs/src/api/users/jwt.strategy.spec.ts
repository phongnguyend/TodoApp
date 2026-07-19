import { ConfigService } from '@nestjs/config';
import { UnauthorizedException } from '@nestjs/common';
import { JwtStrategy } from './jwt.strategy';

describe('JwtStrategy', () => {
  const strategy = new JwtStrategy({ get: jest.fn(() => 'test-secret') } as unknown as ConfigService);

  it('maps a valid subject to the authenticated user', () => {
    expect(strategy.validate({ sub: '42' })).toEqual({ userId: 42 });
  });

  it('rejects an invalid subject', () => {
    expect(() => strategy.validate({ sub: 'invalid' })).toThrow(UnauthorizedException);
  });
});
