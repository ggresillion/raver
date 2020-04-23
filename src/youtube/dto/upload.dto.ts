import {IsNotEmpty, IsNumber, IsOptional, IsString} from 'class-validator';

export class UploadDto {
  @IsString()
  @IsNotEmpty()
  public url: string;
  @IsNumber()
  @IsOptional()
  public categoryId: number;
  @IsString()
  @IsNotEmpty()
  public name: string;
  @IsString()
  @IsNotEmpty()
  public guildId: string;
}
