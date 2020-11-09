import { Module } from '@nestjs/common';
import { BotService } from './bot.service';
import { StorageModule } from '../storage/storage.module';
import { BotController } from './bot.controller';
import { BotGateway } from './bot.gateway';
import { AuthModule } from '../auth/auth.module';

@Module({
  imports: [
    StorageModule,
    AuthModule
  ],
  providers: [BotService, BotGateway],
  exports: [BotService],
  controllers: [BotController],
})
export class BotModule { }
