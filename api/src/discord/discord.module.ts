import { Module, forwardRef } from '@nestjs/common';
import { DiscordService } from './discord.service';
import { BotModule } from '../bot/bot.module';
import { HttpModule } from '@nestjs/axios';

@Module({
  imports: [forwardRef(() => BotModule), HttpModule],
  providers: [DiscordService],
  exports: [DiscordService],
})
export class DiscordModule {
}
