import {Controller, Get} from '@nestjs/common';
import {BotService} from './bot.service';

@Controller('bot')
export class BotController {

  constructor(
    private readonly botService: BotService,
  ) {
  }

  @Get('infos')
  public async getBotInfos() {
    // return this.botService.getInfos();
  }
}
