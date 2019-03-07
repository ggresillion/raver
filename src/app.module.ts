import {Module} from '@nestjs/common';
import {AppController} from './app.controller';
import {AppService} from './app.service';
import {BotModule} from './bot/bot.module';
import {SoundModule} from './sound/sound.module';
import {StorageModule} from './storage/storage.module';

@Module({
  imports: [BotModule, SoundModule, StorageModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {
}
