import {IsEmail, IsNotEmpty, IsString} from 'class-validator';

export class CreateUserDto {

  @IsEmail()
  @IsNotEmpty()
  public readonly email: string;
  @IsString()
  @IsNotEmpty()
  public readonly username: string;
  @IsString()
  @IsNotEmpty()
  public readonly password: string;

}
