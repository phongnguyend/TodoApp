import { Injectable } from '@nestjs/common';
import { Request } from 'express';
import { AuthGuard } from '@nestjs/passport';

export type AuthenticatedRequest = Request & { user: { userId: number } };

@Injectable()
export class UserAuthGuard extends AuthGuard('jwt') {}
