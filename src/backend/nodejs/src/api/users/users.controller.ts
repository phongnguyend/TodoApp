import { Body, Controller, Get, HttpCode, HttpStatus, Param, ParseIntPipe, Patch, Post, Put, Query, Req, UseGuards } from '@nestjs/common';
import { ApiAcceptedResponse, ApiBearerAuth, ApiCreatedResponse, ApiNoContentResponse, ApiOkResponse, ApiQuery, ApiTags } from '@nestjs/swagger';
import { PaginatedResponseDto } from '../../shared/common/dto/paginated-response.dto';
import { CreateUserDto, SignUpDto } from './dto/create-user.dto';
import { ChangePasswordDto, ConfirmPasswordResetDto, ResetPasswordDto } from './dto/password.dto';
import { UpdateProfileDto, UpdateUserDto } from './dto/update-user.dto';
import { UserResponseDto } from './dto/user-response.dto';
import { AuthenticatedRequest, UserAuthGuard } from './users.security';
import { UsersService } from './users.service';

@ApiTags('Users')
@Controller('api/users')
export class UsersController {
  constructor(private readonly service: UsersService) {}

  @Post('signup')
  @HttpCode(HttpStatus.CREATED)
  @ApiCreatedResponse({ type: UserResponseDto })
  signup(@Body() dto: SignUpDto): Promise<UserResponseDto> {
    return this.service.signup(dto);
  }

  @Post('password/reset')
  @HttpCode(HttpStatus.ACCEPTED)
  @ApiAcceptedResponse({ description: 'Password reset request accepted' })
  async requestPasswordReset(@Body() dto: ResetPasswordDto): Promise<{ message: string }> {
    await this.service.requestPasswordReset(dto);
    return { message: 'If the account exists, a password reset email has been queued.' };
  }

  @Post('password/confirm')
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiNoContentResponse()
  confirmPasswordReset(@Body() dto: ConfirmPasswordResetDto): Promise<void> {
    return this.service.confirmPasswordReset(dto);
  }

  @Post('password/change')
  @UseGuards(UserAuthGuard)
  @ApiBearerAuth()
  @HttpCode(HttpStatus.NO_CONTENT)
  @ApiNoContentResponse()
  changePassword(@Req() request: AuthenticatedRequest, @Body() dto: ChangePasswordDto): Promise<void> {
    return this.service.changePassword(request.user.userId, dto);
  }

  @Get('profile')
  @UseGuards(UserAuthGuard)
  @ApiBearerAuth()
  @ApiOkResponse({ type: UserResponseDto })
  getProfile(@Req() request: AuthenticatedRequest): Promise<UserResponseDto> {
    return this.service.getProfile(request.user.userId);
  }

  @Put('profile')
  @UseGuards(UserAuthGuard)
  @ApiBearerAuth()
  @ApiOkResponse({ type: UserResponseDto })
  updateProfile(@Req() request: AuthenticatedRequest, @Body() dto: UpdateProfileDto): Promise<UserResponseDto> {
    return this.service.updateProfile(request.user.userId, dto);
  }

  @Get()
  @ApiQuery({ name: 'page', required: false, type: Number })
  @ApiQuery({ name: 'pageSize', required: false, type: Number })
  @ApiOkResponse({ description: 'Paginated list of users' })
  getAll(
    @Query('page', new ParseIntPipe({ optional: true })) page = 1,
    @Query('pageSize', new ParseIntPipe({ optional: true })) pageSize = 20,
  ): Promise<PaginatedResponseDto<UserResponseDto>> {
    return this.service.getAll(page, pageSize);
  }

  @Post()
  @HttpCode(HttpStatus.CREATED)
  @ApiCreatedResponse({ type: UserResponseDto })
  create(@Body() dto: CreateUserDto): Promise<UserResponseDto> {
    return this.service.create(dto);
  }

  @Get(':id')
  @ApiOkResponse({ type: UserResponseDto })
  getById(@Param('id', ParseIntPipe) id: number): Promise<UserResponseDto> {
    return this.service.getById(id);
  }

  @Put(':id')
  @ApiOkResponse({ type: UserResponseDto })
  update(@Param('id', ParseIntPipe) id: number, @Body() dto: UpdateUserDto): Promise<UserResponseDto> {
    return this.service.update(id, dto);
  }

  @Patch(':id/activate')
  @ApiOkResponse({ type: UserResponseDto })
  activate(@Param('id', ParseIntPipe) id: number): Promise<UserResponseDto> {
    return this.service.setActive(id, true);
  }

  @Patch(':id/deactivate')
  @ApiOkResponse({ type: UserResponseDto })
  deactivate(@Param('id', ParseIntPipe) id: number): Promise<UserResponseDto> {
    return this.service.setActive(id, false);
  }
}
