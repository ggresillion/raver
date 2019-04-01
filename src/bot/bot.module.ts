import {Module} from '@nestjs/common';
import {BotService} from './bot.service';
import {StorageModule} from '../storage/storage.module';

@Module({
  imports: [StorageModule],
  providers: [BotService],
  exports: [BotService],
})
export class BotModule {}
