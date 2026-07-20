import { ApiProperty } from '@nestjs/swagger';

export class TodoItemAttachmentResponseDto {
  @ApiProperty()
  id!: number;

  @ApiProperty()
  todoItemId!: number;

  @ApiProperty()
  fileId!: number;

  @ApiProperty()
  createdAt!: Date;

  @ApiProperty({ nullable: true })
  createdByUserId!: number | null;

  @ApiProperty({ nullable: true })
  updatedAt!: Date | null;

  @ApiProperty({ nullable: true })
  updatedByUserId!: number | null;
}
