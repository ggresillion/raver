import {Injectable, UnauthorizedException} from '@nestjs/common';
import fetch from 'node-fetch';

@Injectable()
export class DiscordService {

  private discordApi = 'http://discordapp.com/api';

  /**
   * Get user informations from the discord API
   * @param token
   */
  public async getUser(token: string) {
    return fetch(`${this.discordApi}/users/@me`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    }).then(res => {
      if (!res.ok) {
        throw new UnauthorizedException('Failed to contact Discord API');
      }
      return res.json();
    });
  }
}
