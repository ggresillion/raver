import {Injectable} from '@nestjs/common';
import {UserService} from '../user/user.service';
import {JwtService} from '@nestjs/jwt';
import {LoginDto} from './dto/login.dto';
import {CreateUserDto} from '../user/dto/create-user.dto';
import {GetUserDTO} from '../user/dto/get-user.dto';

@Injectable()
export class AuthService {

  constructor(
    private readonly userService: UserService,
    private readonly jwtService: JwtService,
  ) {
  }

  public async login(email: string, password: string): Promise<LoginDto> {
    const user = await this.userService.findOneByEmailAndPassword(email, password);
    return this.getTokenForUser(user);
  }

  public async register(userDto: CreateUserDto) {
    const user = await this.userService.createUser(userDto);
    return this.getTokenForUser(user);
  }

  private getTokenForUser(user: GetUserDTO) {
    return {token: this.jwtService.sign({user}), expiresIn: 3600, user};
  }
}
