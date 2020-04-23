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
import { join } from 'path';
import { ServeStaticModule } from '@nestjs/serve-static';

@Module({
  imports: [
    TypeOrmModule.forRoot({
      type: "postgres",
      host: "localhost",
      port: 5432,
      database: "dsb",
      username: "dsb",
      password: "dsb",
      autoLoadEntities: true
    }),
    ConfigModule.forRoot({
      load: [configuration],
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
    ServeStaticModule.forRoot({
      rootPath: join(__dirname, 'client'),
    }),
  ],
  controllers: [AppController],
  providers: [],
})
export class AppModule {
}
