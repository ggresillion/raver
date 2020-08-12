import {IsInt, IsNotEmpty, IsOptional, IsString, ValidateIf} from 'class-validator';

export class SoundDto {
  @IsString()
  @IsNotEmpty()
  @IsOptional()
  public readonly name;
  @ValidateIf((el) => !isNaN(el))
  @IsInt()
  @IsOptional()
  public readonly categoryId;
  @IsString()
  @IsOptional()
  public readonly guildId;
}
