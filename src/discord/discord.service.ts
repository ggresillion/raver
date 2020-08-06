import {Injectable, Scope, UnauthorizedException} from '@nestjs/common';
import fetch from 'node-fetch';
import {UserDTO} from '../user/dto/user.dto';
import {GuildDTO} from '../guild/dto/guild.dto';
import {BotService} from '../bot/bot.service';

@Injectable({scope: Scope.REQUEST})
export class DiscordService {

  constructor(private readonly botService: BotService) {
  }

  private readonly discordApi = 'http://discordapp.com/api';
  private token;

  public setToken(token: string) {
    this.token = token;
  }

  /**
   * Get connected user information from the discord API
   */
  public async getUser(): Promise<UserDTO> {
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

  /**
   * Return the user's guilds
   * @param user user
   */
  public async getMyGuilds(user: UserDTO): Promise<GuildDTO[]> {
    const res = await fetch(`${this.discordApi}/users/@me/guilds`, {
      headers: {
        Authorization: `Bearer ${this.token}`,
      },
    });
    if (!res.ok) {
      throw new UnauthorizedException('Failed to contact Discord API');
    }
    return this.botService.setIsBotInGuild(user, await res.json() as GuildDTO[]);
  }
}
