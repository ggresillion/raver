import {Module, forwardRef} from '@nestjs/common';
import {DiscordService} from './discord.service';
import {BotModule} from '../bot/bot.module';

@Module({
  imports: [forwardRef(() => BotModule)],
  providers: [DiscordService],
  exports: [DiscordService],
})
export class DiscordModule {
}
