import {User} from '../user.entity';

export class GetUserDTO {

  public readonly id: number;
  public readonly email: string;
  public readonly username: string;

  constructor(id: number, email: string, username: string) {
    this.id = id;
    this.email = email;
    this.username = username;
  }

  public static toDTO(user: User) {
    return new GetUserDTO(user.id, user.email, user.username);
  }
}
