import {Injectable} from '@nestjs/common';
import {BotService} from '../bot/bot.service';
import {UserDTO} from '../user/dto/user.dto';
import {GuildDTO} from './dto/guild.dto';
import {DiscordService} from '../discord/discord.service';

@Injectable()
export class GuildService {

  constructor(private readonly discordService: DiscordService) {
  }

  public async getAvailableGuildsForUser(user: UserDTO): Promise<GuildDTO[]> {
    return await this.discordService.getMyGuilds(user);
  }
}
