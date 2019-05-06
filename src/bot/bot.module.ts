import {Module} from '@nestjs/common';
import {BotService} from './bot.service';
import {StorageModule} from '../storage/storage.module';
import {BotController} from './bot.controller';

@Module({
  imports: [StorageModule],
  providers: [BotService],
  exports: [BotService],
  controllers: [BotController],
})
export class BotModule {}
