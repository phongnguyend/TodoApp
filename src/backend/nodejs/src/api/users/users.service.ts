import { BadRequestException, ConflictException, Injectable, NotFoundException } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { User } from '@prisma/client';
import { createHash, createHmac, pbkdf2Sync, randomBytes, timingSafeEqual } from 'crypto';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { CreateUserDto, SignUpDto } from './dto/create-user.dto';
import { ChangePasswordDto, ConfirmPasswordResetDto, ResetPasswordDto } from './dto/password.dto';
import { UpdateProfileDto, UpdateUserDto } from './dto/update-user.dto';
import { UserResponseDto } from './dto/user-response.dto';
import { UsersRepository } from './users.repository';

@Injectable()
export class UsersService {
  constructor(private readonly repository: UsersRepository, private readonly config: ConfigService) {}

  private static toDto(user: User): UserResponseDto {
    const { id, username, email, isActive, createdAt, updatedAt } = user;
    return { id, username, email, isActive, createdAt, updatedAt };
  }

  private async getOrThrow(id: number): Promise<User> {
    const user = await this.repository.findById(id);
    if (!user) throw new NotFoundException(`User ${id} not found.`);
    return user;
  }

  private async ensureUnique(username: string, email: string, excludingId?: number): Promise<void> {
    if (await this.repository.usernameExists(username, excludingId)) {
      throw new ConflictException('Username is already in use.');
    }
    if (await this.repository.emailExists(email, excludingId)) {
      throw new ConflictException('Email is already in use.');
    }
  }

  private hashPassword(password: string): string {
    const iterations = Number(this.config.get<string>('PASSWORD_HASH_ITERATIONS') ?? 120000);
    const salt = randomBytes(16);
    const digest = pbkdf2Sync(password, salt, iterations, 32, 'sha256');
    return `pbkdf2_sha256$${iterations}$${salt.toString('hex')}$${digest.toString('hex')}`;
  }

  private verifyPassword(password: string, encoded: string): boolean {
    try {
      const [algorithm, iterationText, saltText, digestText] = encoded.split('$');
      if (algorithm !== 'pbkdf2_sha256') return false;
      const expected = Buffer.from(digestText, 'hex');
      const actual = pbkdf2Sync(password, Buffer.from(saltText, 'hex'), Number(iterationText), expected.length, 'sha256');
      return expected.length > 0 && timingSafeEqual(actual, expected);
    } catch {
      return false;
    }
  }

  async getAll(page: number, pageSize: number): Promise<PaginatedResponseDto<UserResponseDto>> {
    page = Math.max(1, page);
    pageSize = Math.min(100, Math.max(1, pageSize));
    const { items, total } = await this.repository.findAll((page - 1) * pageSize, pageSize);
    return { items: items.map(UsersService.toDto), total, page, pageSize, totalPages: Math.ceil(total / pageSize) };
  }

  async getById(id: number): Promise<UserResponseDto> {
    return UsersService.toDto(await this.getOrThrow(id));
  }

  async create(dto: CreateUserDto): Promise<UserResponseDto> {
    const username = dto.username.trim();
    const email = dto.email.trim().toLowerCase();
    await this.ensureUnique(username, email);
    return UsersService.toDto(await this.repository.create({
      username, email, passwordHash: this.hashPassword(dto.password), isActive: dto.isActive ?? true,
    }));
  }

  async update(id: number, dto: UpdateUserDto): Promise<UserResponseDto> {
    const existing = await this.getOrThrow(id);
    const username = dto.username?.trim() ?? existing.username;
    const email = dto.email?.trim().toLowerCase() ?? existing.email;
    await this.ensureUnique(username, email, id);
    return UsersService.toDto(await this.repository.update(id, {
      username, email, ...(dto.password ? { passwordHash: this.hashPassword(dto.password) } : {}),
    }));
  }

  async setActive(id: number, isActive: boolean): Promise<UserResponseDto> {
    await this.getOrThrow(id);
    return UsersService.toDto(await this.repository.update(id, { isActive }));
  }

  signup(dto: SignUpDto): Promise<UserResponseDto> {
    return this.create({ ...dto, isActive: true });
  }

  getProfile(userId: number): Promise<UserResponseDto> {
    return this.getById(userId);
  }

  updateProfile(userId: number, dto: UpdateProfileDto): Promise<UserResponseDto> {
    return this.update(userId, dto);
  }

  async changePassword(userId: number, dto: ChangePasswordDto): Promise<void> {
    const user = await this.getOrThrow(userId);
    if (!user.isActive) throw new BadRequestException('The user account is inactive.');
    if (!this.verifyPassword(dto.currentPassword, user.passwordHash)) {
      throw new BadRequestException('The current password is incorrect.');
    }
    await this.repository.update(userId, { passwordHash: this.hashPassword(dto.newPassword) });
  }

  async requestPasswordReset(dto: ResetPasswordDto): Promise<void> {
    const user = await this.repository.findByEmail(dto.email.trim().toLowerCase());
    if (!user?.isActive) return;
    const lifetime = Math.max(1, Number(this.config.get<string>('PASSWORD_RESET_TOKEN_LIFETIME_MINUTES') ?? 60));
    const payload = Buffer.from(JSON.stringify({
      sub: user.id,
      exp: Math.floor(Date.now() / 1000) + lifetime * 60,
      password: createHash('sha256').update(user.passwordHash).digest('hex'),
    })).toString('base64url');
    const secret = this.config.get<string>('PASSWORD_RESET_SECRET_KEY') ?? this.config.get<string>('JWT_SECRET_KEY') ?? 'change-me';
    const signature = createHmac('sha256', secret).update(payload).digest('base64url');
    const token = `${payload}.${signature}`;
    const baseUrl = this.config.get<string>('PASSWORD_RESET_CONFIRMATION_URL') ?? '/reset-password';
    const url = `${baseUrl}${baseUrl.includes('?') ? '&' : '?'}token=${encodeURIComponent(token)}`;
    await this.repository.createEmailLog({
      recipient: user.email,
      subject: 'Reset your Todo API password',
      body: `Use this link to reset your password: ${url}\n\nThis link expires in ${lifetime} minutes.`,
      status: 'pending',
    });
  }

  async confirmPasswordReset(dto: ConfirmPasswordResetDto): Promise<void> {
    try {
      const [payloadText, signatureText, extra] = dto.token.split('.');
      if (!payloadText || !signatureText || extra) throw new Error();
      const secret = this.config.get<string>('PASSWORD_RESET_SECRET_KEY') ?? this.config.get<string>('JWT_SECRET_KEY') ?? 'change-me';
      const expected = createHmac('sha256', secret).update(payloadText).digest();
      const supplied = Buffer.from(signatureText, 'base64url');
      if (supplied.length !== expected.length || !timingSafeEqual(supplied, expected)) throw new Error();
      const payload = JSON.parse(Buffer.from(payloadText, 'base64url').toString()) as { sub?: number; exp?: number; password?: string };
      const user = payload.sub ? await this.repository.findById(payload.sub) : null;
      const fingerprint = user ? createHash('sha256').update(user.passwordHash).digest('hex') : '';
      if (!user?.isActive || !payload.exp || payload.exp < Math.floor(Date.now() / 1000) || payload.password !== fingerprint) throw new Error();
      await this.repository.update(user.id, { passwordHash: this.hashPassword(dto.newPassword) });
    } catch {
      throw new BadRequestException('The password reset token is invalid or expired.');
    }
  }
}
