import { Injectable, Scope } from '@nestjs/common';
import { UserDTO } from '../user/dto/user.dto';
import { GuildDTO } from '../guild/dto/guild.dto';
import { BotService } from '../bot/bot.service';
import { HttpService } from '@nestjs/axios';
import { firstValueFrom, map } from 'rxjs';

@Injectable({ scope: Scope.REQUEST })
export class DiscordService {

  constructor(private readonly botService: BotService, private readonly httpService: HttpService) {
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
    return firstValueFrom(this.httpService.get(`${this.discordApi}/users/@me`, {
      headers: {
        Authorization: `Bearer ${this.token}`,
      }
    }).pipe(map(res => res.data)));
  }

  /**
   * Return the user's guilds
   * @param user user
   */
  public async getMyGuilds(user: UserDTO): Promise<GuildDTO[]> {
    return firstValueFrom(this.httpService.get(`${this.discordApi}/users/@me/guilds`, {
      headers: {
        Authorization: `Bearer ${this.token}`,
      }
    }).pipe(map(res => this.botService.setIsBotInGuild(user, res.data as GuildDTO[]))));
  }
}
