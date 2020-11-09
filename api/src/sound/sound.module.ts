import {Module} from '@nestjs/common';
import {SoundService} from './sound.service';
import {SoundController} from './sound.controller';
import {StorageModule} from '../storage/storage.module';
import {TypeOrmModule} from '@nestjs/typeorm';
import {Sound} from './entity/sound.entity';
import {Image} from '../image/entity/image.entity';
import {AuthModule} from '../auth/auth.module';
import {BotModule} from '../bot/bot.module';

@Module({
  imports: [StorageModule, BotModule, TypeOrmModule.forFeature([Sound, Image]), AuthModule],
  providers: [SoundService],
  controllers: [SoundController],
  exports: [SoundService],
})
export class SoundModule {}
