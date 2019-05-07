import {Module} from '@nestjs/common';
import {BotService} from './bot.service';
import {StorageModule} from '../storage/storage.module';
import {BotController} from './bot.controller';
import {BotGateway} from './bot.gateway';

@Module({
  imports: [StorageModule],
  providers: [BotService, BotGateway],
  exports: [BotService],
  controllers: [BotController],
})
export class BotModule {}
