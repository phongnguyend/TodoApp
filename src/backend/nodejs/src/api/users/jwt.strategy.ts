import { Injectable, UnauthorizedException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { PassportStrategy } from '@nestjs/passport';
import { ExtractJwt, Strategy } from 'passport-jwt';

@Injectable()
export class JwtStrategy extends PassportStrategy(Strategy) {
  constructor(config: ConfigService) {
    super({
      jwtFromRequest: ExtractJwt.fromAuthHeaderAsBearerToken(),
      ignoreExpiration: false,
      secretOrKey: config.get<string>('JWT_SECRET_KEY') ?? 'change-me-use-at-least-32-bytes-long',
      algorithms: ['HS256'],
    });
  }

  validate(payload: { sub?: string | number }): { userId: number } {
    const userId = Number(payload.sub);
    if (!Number.isInteger(userId) || userId < 1) {
      throw new UnauthorizedException('Invalid authentication token.');
    }
    return { userId };
  }
}
