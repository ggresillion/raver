import {IsInt, IsNotEmpty, IsOptional, IsString} from 'class-validator';

export class SoundDto {
  @IsString()
  @IsNotEmpty()
  @IsOptional()
  public readonly name;
  @IsInt()
  @IsOptional()
  public readonly categoryId;
}
