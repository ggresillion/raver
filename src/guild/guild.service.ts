import {Injectable} from '@nestjs/common';
import {BotService} from '../bot/bot.service';
import {UserDTO} from '../user/dto/user.dto';
import {GuildDTO} from './dto/guild.dto';

@Injectable()
export class GuildService {

  constructor(private readonly botService: BotService) {
  }

  public getAvailableGuildsForUser(user: UserDTO): GuildDTO[] {
    return this.botService.getGuildsForUser(user);
  }
}
