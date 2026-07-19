import { ApiPropertyOptional } from '@nestjs/swagger';
import { IsEmail, IsOptional, IsString, Length, Matches, MaxLength } from 'class-validator';

export class UpdateUserDto {
  @ApiPropertyOptional({ maxLength: 50 })
  @IsOptional()
  @IsString()
  @Matches(/\S/, { message: 'username must not be blank' })
  @MaxLength(50)
  username?: string;

  @ApiPropertyOptional({ maxLength: 255 })
  @IsOptional()
  @IsEmail()
  @MaxLength(255)
  email?: string;

  @ApiPropertyOptional({ minLength: 8, maxLength: 128 })
  @IsOptional()
  @IsString()
  @Length(8, 128)
  password?: string;
}

export class UpdateProfileDto {
  @ApiPropertyOptional({ maxLength: 50 })
  @IsOptional()
  @IsString()
  @Matches(/\S/, { message: 'username must not be blank' })
  @MaxLength(50)
  username?: string;

  @ApiPropertyOptional({ maxLength: 255 })
  @IsOptional()
  @IsEmail()
  @MaxLength(255)
  email?: string;
}
