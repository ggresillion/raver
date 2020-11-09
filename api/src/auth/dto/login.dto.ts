import {UserDTO} from '../../user/dto/user.dto';

export class LoginDto {
  token: string;
  user: UserDTO;
  expiresIn: number;
}
