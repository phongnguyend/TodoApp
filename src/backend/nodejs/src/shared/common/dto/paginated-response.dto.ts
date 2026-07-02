import { ApiProperty } from '@nestjs/swagger';

/** Generic paginated response - mirrors a PagedResult<T> in C# */
export class PaginatedResponseDto<T> {
  @ApiProperty({ isArray: true })
  items!: T[];

  @ApiProperty() total!: number;
  @ApiProperty() page!: number;
  @ApiProperty() pageSize!: number;
  @ApiProperty() totalPages!: number;
}
