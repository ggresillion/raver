import {HttpModule, Module} from '@nestjs/common';
import { YoutubeService } from './youtube.service';
import { YoutubeController } from './youtube.controller';
import {StorageModule} from '../storage/storage.module';
import {SoundModule} from '../sound/sound.module';
import {YoutubeGateway} from './youtube-gateway';

@Module({
  imports: [SoundModule, StorageModule, HttpModule],
  providers: [YoutubeService, YoutubeGateway],
  controllers: [YoutubeController],
})
export class YoutubeModule {}
