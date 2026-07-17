import { ApiProperty } from '@nestjs/swagger';

/** A single failed row encountered while importing todo items. */
export class ImportRowErrorDto {
  @ApiProperty() row!: number;
  @ApiProperty() error!: string;
}

/** Result summary returned by the CSV/Excel import endpoints. */
export class ImportResultDto {
  @ApiProperty() imported!: number;
  @ApiProperty() failed!: number;
  @ApiProperty({ type: [ImportRowErrorDto] }) errors!: ImportRowErrorDto[];
}
