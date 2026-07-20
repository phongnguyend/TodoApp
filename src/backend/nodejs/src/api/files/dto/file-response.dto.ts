import { ApiProperty } from '@nestjs/swagger';

/**
 * Mirrors a FileDto / view model returned from controllers in C#.
 * Note: the on-disk `location` is intentionally not exposed to clients; content is
 * retrieved via the dedicated download endpoint instead.
 */
export class FileResponseDto {
  @ApiProperty() id!: number;
  @ApiProperty() name!: string;
  @ApiProperty() extension!: string;
  @ApiProperty() size!: number;
  @ApiProperty({ nullable: true }) contentType!: string | null;
  @ApiProperty() createdAt!: Date;
  @ApiProperty({ nullable: true }) createdByUserId!: number | null;
  @ApiProperty({ nullable: true }) updatedAt!: Date | null;
  @ApiProperty({ nullable: true }) updatedByUserId!: number | null;
}
