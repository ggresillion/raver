import {Injectable, Scope, UnauthorizedException} from '@nestjs/common';
import fetch from 'node-fetch';

@Injectable({scope: Scope.REQUEST})
export class DiscordService {

  private readonly discordApi = 'http://discordapp.com/api';
  private token;

  public setToken(token: string) {
    this.token = token;
  }

  /**
   * Get connected user information from the discord API
   */
  public async getUser() {
    return fetch(`${this.discordApi}/users/@me`, {
      headers: {
        Authorization: `Bearer ${this.token}`,
      },
    }).then(res => {
      if (!res.ok) {
        throw new UnauthorizedException('Failed to contact Discord API');
      }
      return res.json();
    });
  }
}
