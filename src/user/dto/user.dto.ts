export class UserDTO {

  public readonly id: string;
  public readonly email: string;
  public readonly username: string;

  constructor(id: string, email: string, username: string) {
    this.id = id;
    this.email = email;
    this.username = username;
  }
}
