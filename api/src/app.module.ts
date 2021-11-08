import { Module } from '@nestjs/common';
import { AppController } from './app.controller';
import { BotModule } from './bot/bot.module';
import { SoundModule } from './sound/sound.module';
import { StorageModule } from './storage/storage.module';
import { TypeOrmModule } from '@nestjs/typeorm';
import { UserModule } from './user/user.module';
import { AuthModule } from './auth/auth.module';
import { ConfigModule } from '@nestjs/config';
import { DiscordModule } from './discord/discord.module';
import { CategoryModule } from './category/category.module';
import configuration from './config/configuration';
import { YoutubeModule } from './youtube/youtube.module';
import { GuildModule } from './guild/guild.module';
import { ConfigService } from '@nestjs/config';
import { ImageModule } from './image/image.module';

@Module({
  imports: [
    ConfigModule.forRoot({
      load: [configuration],
      isGlobal: true
    }),
    TypeOrmModule.forRootAsync({
      imports: [ConfigModule],
      useFactory: (config: ConfigService) => config.get('database'),
      inject: [ConfigService],
    }),
    BotModule,
    SoundModule,
    StorageModule,
    UserModule,
    AuthModule,
    DiscordModule,
    YoutubeModule,
    CategoryModule,
    GuildModule,
    ImageModule,
  ],
  controllers: [AppController],
  providers: [],
})
export class AppModule {
}
