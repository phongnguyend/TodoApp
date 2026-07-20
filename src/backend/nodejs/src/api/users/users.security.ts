import { ExecutionContext, Injectable, SetMetadata } from '@nestjs/common';
import { Reflector } from '@nestjs/core';
import { Request } from 'express';
import { AuthGuard } from '@nestjs/passport';

export type AuthenticatedRequest = Request & { user: { userId: number } };
export type AuditRequest = Request & { user?: { userId: number } };

const PUBLIC_ENDPOINT_KEY = 'publicEndpoint';

export const PublicEndpoint = () => SetMetadata(PUBLIC_ENDPOINT_KEY, true);

@Injectable()
export class UserAuthGuard extends AuthGuard('jwt') {
  constructor(private readonly reflector: Reflector) {
    super();
  }

  canActivate(context: ExecutionContext) {
    const isPublic = this.reflector.getAllAndOverride<boolean>(PUBLIC_ENDPOINT_KEY, [
      context.getHandler(),
      context.getClass(),
    ]);
    return isPublic ? true : super.canActivate(context);
  }
}
