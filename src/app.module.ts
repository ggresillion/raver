import {Module} from '@nestjs/common';
import {AppController} from './app.controller';
import {BotModule} from './bot/bot.module';
import {SoundModule} from './sound/sound.module';
import {StorageModule} from './storage/storage.module';
import {TypeOrmModule} from '@nestjs/typeorm';
import {UserModule} from './user/user.module';
import {AuthModule} from './auth/auth.module';
import {ConfigModule} from 'nestjs-config';
import {DiscordModule} from './discord/discord.module';
import { YoutubeModule } from './youtube/youtube.module';
import { CategoryModule } from './category/category.module';
import * as path from 'path';

@Module({
  imports: [
    TypeOrmModule.forRoot(),
    ConfigModule.load(path.resolve(__dirname, 'config', '**/!(*.d).{ts,js}')),
    BotModule, SoundModule, StorageModule, UserModule, AuthModule, DiscordModule, YoutubeModule, CategoryModule,
  ],
  controllers: [AppController],
  providers: [],
})
export class AppModule {
}
