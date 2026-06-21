import { ApiProperty } from '@nestjs/swagger';
import { IsNotEmpty, IsOptional, IsString, MaxLength } from 'class-validator';

/** Mirrors a CreateTodoItemRequest / command model in C# */
export class CreateTodoItemDto {
  @ApiProperty({ example: 'Buy groceries', maxLength: 200 })
  @IsString()
  @IsNotEmpty()
  @MaxLength(200)
  title!: string;

  @ApiProperty({ example: 'Milk, eggs, bread', required: false, nullable: true })
  @IsString()
  @IsOptional()
  @MaxLength(2000)
  description?: string;
}
