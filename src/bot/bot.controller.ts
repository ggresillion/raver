import {Controller, Get, Post, UseGuards} from '@nestjs/common';
import {BotService} from './bot.service';
import { UserGuard } from '../auth/guards/user.guard';
import { ConnectedUser } from '../auth/decorators/connected-user.decorator';
import { UserDTO } from '../user/dto/user.dto';

@Controller('bot')
@UseGuards(UserGuard)
export class BotController {

  constructor(
    private readonly botService: BotService,
  ) {
  }

  @Get('infos')
  public async getBotInfos() {
    // return this.botService.getInfos();
  }

  @Post('join')
  public async joinMyChannel(@ConnectedUser() user: UserDTO) {
    await this.botService.joinMyChannel(user);
  }
}
