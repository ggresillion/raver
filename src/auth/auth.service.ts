import {Injectable} from '@nestjs/common';
import {UserService} from '../user/user.service';
import {JwtService} from '@nestjs/jwt';
import fetch from 'node-fetch';

@Injectable()
export class AuthService {

  constructor(
    private readonly userService: UserService,
    private readonly jwtService: JwtService,
  ) {
  }

  /**
   * Get user informations from the discord API
   * @param token
   */
  public static async getDiscordUser(token: string) {
    return fetch('http://discordapp.com/api/users/@me', {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }).then(res => res.json());
  }

  public async generateToken(discordToken: string, discordRefreshToken: string) {
    const user = AuthService.getDiscordUser(discordToken);
    return {
      token: this.jwtService.sign({
        token: discordToken,
        refreshToken: discordRefreshToken,
        ...user,
      }),
      expiresIn: 3600,
    };
  }
}
