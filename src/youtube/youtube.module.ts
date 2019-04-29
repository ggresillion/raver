import { Module } from '@nestjs/common';
import { YoutubeService } from './youtube.service';
import { YoutubeController } from './youtube.controller';
import {StorageModule} from '../storage/storage.module';
import {SoundModule} from '../sound/sound.module';

@Module({
  imports: [SoundModule, StorageModule],
  providers: [YoutubeService],
  controllers: [YoutubeController],
})
export class YoutubeModule {}
