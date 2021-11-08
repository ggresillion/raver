import {Module} from '@nestjs/common';
import {AuthService} from './auth.service';
import {AuthController} from './auth.controller';
import {DiscordModule} from '../discord/discord.module';
import {ConfigModule} from '@nestjs/config';
import { HttpModule } from '@nestjs/axios';

@Module({
  imports: [
    DiscordModule,
    ConfigModule,
    HttpModule
  ],
  providers: [AuthService],
  controllers: [AuthController],
  exports: [
    AuthService,
    DiscordModule,
    ConfigModule
  ],
})
export class AuthModule {
}
