import { ApiProperty } from '@nestjs/swagger';

export class UserResponseDto {
  @ApiProperty() id!: number;
  @ApiProperty() username!: string;
  @ApiProperty() email!: string;
  @ApiProperty() isActive!: boolean;
  @ApiProperty() createdAt!: Date;
  @ApiProperty({ nullable: true }) createdByUserId!: number | null;
  @ApiProperty({ nullable: true }) updatedAt!: Date | null;
  @ApiProperty({ nullable: true }) updatedByUserId!: number | null;
}
