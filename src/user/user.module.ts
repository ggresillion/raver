import {forwardRef, Module} from '@nestjs/common';
import {UserService} from './user.service';
import {UserController} from './user.controller';
import {AuthModule} from '../auth/auth.module';
import {DiscordModule} from '../discord/discord.module';
import {ConfigModule} from '@nestjs/config';

@Module({
  imports: [
    forwardRef(() => AuthModule),
    DiscordModule,
    ConfigModule,
  ],
  providers: [UserService],
  controllers: [UserController],
  exports: [UserService],
})
export class UserModule {
}
