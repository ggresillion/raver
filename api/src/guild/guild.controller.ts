import {Controller, Get, UseGuards} from '@nestjs/common';
import {UserGuard} from '../auth/guards/user.guard';
import {ConnectedUser} from '../auth/decorators/connected-user.decorator';
import {GuildService} from './guild.service';
import {GuildDTO} from './dto/guild.dto';

@Controller('guilds')
@UseGuards(UserGuard)
export class GuildController {

  constructor(private readonly guildService: GuildService) {
  }

  @Get('/')
  public async getAvailableGuilds(@ConnectedUser() user): Promise<GuildDTO[]> {
    return await this.guildService.getAvailableGuildsForUser(user);
  }
}
