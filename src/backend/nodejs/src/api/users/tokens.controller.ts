import { Body, Controller, Header, HttpCode, HttpStatus, Post, Res } from '@nestjs/common';
import { ApiOkResponse, ApiTags } from '@nestjs/swagger';
import { Response } from 'express';
import { TokenRequestDto, TokenResponseDto } from './dto/token.dto';
import { UsersService } from './users.service';
import { PublicEndpoint } from './users.security';

@ApiTags('Tokens')
@Controller('api/tokens')
export class TokensController {
  constructor(private readonly service: UsersService) {}

  @Post()
  @PublicEndpoint()
  @HttpCode(HttpStatus.OK)
  @Header('Cache-Control', 'no-store')
  @Header('Pragma', 'no-cache')
  @ApiOkResponse({ type: TokenResponseDto })
  async create(
    @Body() dto: TokenRequestDto,
    @Res({ passthrough: true }) response: Response,
  ): Promise<TokenResponseDto> {
    try {
      return await this.service.createToken(dto);
    } catch (error) {
      if (error instanceof Error && error.name === 'UnauthorizedException') {
        response.setHeader('WWW-Authenticate', 'Bearer');
      }
      throw error;
    }
  }
}
