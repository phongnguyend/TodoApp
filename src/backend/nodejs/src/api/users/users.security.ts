import { CanActivate, ExecutionContext, Injectable, UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { Request } from 'express';
import { createHmac, timingSafeEqual } from 'crypto';

export type AuthenticatedRequest = Request & { userId: number };

const decodeJson = (value: string): Record<string, unknown> =>
  JSON.parse(Buffer.from(value, 'base64url').toString('utf8')) as Record<string, unknown>;

@Injectable()
export class UserAuthGuard implements CanActivate {
  constructor(private readonly config: ConfigService) {}

  canActivate(context: ExecutionContext): boolean {
    const request = context.switchToHttp().getRequest<AuthenticatedRequest>();
    const [scheme, token] = (request.headers.authorization ?? '').split(' ', 2);
    if (scheme?.toLowerCase() !== 'bearer' || !token) {
      throw new UnauthorizedException('Authentication required.');
    }

    try {
      const [header, payload, signature, extra] = token.split('.');
      if (!header || !payload || !signature || extra || decodeJson(header).alg !== 'HS256') throw new Error();
      const secret = this.config.get<string>('JWT_SECRET_KEY') ?? 'change-me';
      const expected = createHmac('sha256', secret).update(`${header}.${payload}`).digest();
      const supplied = Buffer.from(signature, 'base64url');
      if (supplied.length !== expected.length || !timingSafeEqual(supplied, expected)) throw new Error();
      const claims = decodeJson(payload);
      if (typeof claims.exp === 'number' && claims.exp < Math.floor(Date.now() / 1000)) throw new Error();
      const userId = Number(claims.sub);
      if (!Number.isInteger(userId) || userId < 1) throw new Error();
      request.userId = userId;
      return true;
    } catch {
      throw new UnauthorizedException('Invalid or expired token.');
    }
  }
}
