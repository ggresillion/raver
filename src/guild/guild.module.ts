import {Module} from '@nestjs/common';
import {GuildService} from './guild.service';
import {GuildController} from './guild.controller';
import {AuthModule} from '../auth/auth.module';
import {DiscordModule} from '../discord/discord.module';
import {ConfigModule} from '@nestjs/config';

@Module({
  imports: [
    AuthModule,
    DiscordModule,
    ConfigModule,
  ],
  providers: [GuildService],
  controllers: [GuildController],
  exports: [GuildService],
})
export class GuildModule {
}
