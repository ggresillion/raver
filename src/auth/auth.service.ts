import {Injectable} from '@nestjs/common';
import fetch from 'node-fetch';

@Injectable()
export class AuthService {

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
}
