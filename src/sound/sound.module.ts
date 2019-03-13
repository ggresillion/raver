import {Module} from '@nestjs/common';
import {SoundService} from './sound.service';
import {SoundController} from './sound.controller';
import {StorageModule} from '../storage/storage.module';
import {TypeOrmModule} from '@nestjs/typeorm';
import {Sound} from './sound.entity';
import {BotModule} from '../bot/bot.module';

@Module({
  imports: [StorageModule, BotModule, TypeOrmModule.forFeature([Sound])],
  providers: [SoundService],
  controllers: [SoundController],
})
export class SoundModule {}
