import {GetUserDTO} from '../../user/dto/get-user.dto';

export class LoginDto {
  token: string;
  user: GetUserDTO;
  expiresIn: number;
}
