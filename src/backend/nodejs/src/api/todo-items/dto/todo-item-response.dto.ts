import { ApiProperty } from '@nestjs/swagger';

/** Mirrors a TodoItemDto / view model returned from controllers in C# */
export class TodoItemResponseDto {
  @ApiProperty() id!: number;
  @ApiProperty() title!: string;
  @ApiProperty({ nullable: true }) description!: string | null;
  @ApiProperty() isCompleted!: boolean;
  @ApiProperty() createdAt!: Date;
  @ApiProperty({ nullable: true }) updatedAt!: Date | null;
}
