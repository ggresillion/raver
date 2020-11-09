import {IsEmail, IsNotEmpty, IsString} from 'class-validator';

export class CredentialsDto {
  @IsEmail()
  @IsNotEmpty()
  readonly email: string;
  @IsString()
  @IsNotEmpty()
  readonly password: string;
}
