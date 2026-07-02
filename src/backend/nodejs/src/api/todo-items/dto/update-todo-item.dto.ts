import { ApiProperty } from '@nestjs/swagger';
import { IsBoolean, IsOptional, IsString, MaxLength } from 'class-validator';

/** Mirrors an UpdateTodoItemRequest in C# - all fields are optional (PATCH semantics) */
export class UpdateTodoItemDto {
  @ApiProperty({ example: 'Buy groceries', maxLength: 200, required: false })
  @IsString()
  @IsOptional()
  @MaxLength(200)
  title?: string;

  @ApiProperty({ example: 'Milk, eggs, bread', required: false, nullable: true })
  @IsString()
  @IsOptional()
  @MaxLength(2000)
  description?: string;

  @ApiProperty({ example: false, required: false })
  @IsBoolean()
  @IsOptional()
  isCompleted?: boolean;
}
