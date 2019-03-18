import {IsEmail, IsNotEmpty, IsString} from 'class-validator';

export class CreateUserDTO {

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
