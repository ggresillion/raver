import {Injectable} from '@nestjs/common';
import {UserService} from '../user/user.service';
import {JwtService} from '@nestjs/jwt';
import {JwtPayload} from './interfaces/jwt-payload.interface';
import {LoginDto} from './dto/login.dto';

@Injectable()
export class AuthService {

  constructor(
    private readonly userService: UserService,
    private readonly jwtService: JwtService,
  ) {
  }

  public async signIn(email: string, password: string): Promise<LoginDto> {
    return {token: this.jwtService.sign({}), expiresIn: 3600, user: null};
  }

  public async validateUser(payload: JwtPayload): Promise<any> {
    return await this.userService.findOneByEmail(payload.email);
  }
}
