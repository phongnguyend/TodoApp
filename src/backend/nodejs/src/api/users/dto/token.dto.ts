import { ApiProperty } from '@nestjs/swagger';
import { IsEmail, IsNotEmpty, IsString, MaxLength } from 'class-validator';

export class TokenRequestDto {
  @ApiProperty({ maxLength: 255 })
  @IsEmail()
  @MaxLength(255)
  email!: string;

  @ApiProperty({ maxLength: 128 })
  @IsString()
  @IsNotEmpty()
  @MaxLength(128)
  password!: string;
}

export class TokenResponseDto {
  @ApiProperty()
  access_token!: string;

  @ApiProperty({ example: 'Bearer' })
  token_type!: string;

  @ApiProperty({ example: 3600 })
  expires_in!: number;
}
