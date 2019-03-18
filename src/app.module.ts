import {Module} from '@nestjs/common';
import {AppController} from './app.controller';
import {AppService} from './app.service';
import {BotModule} from './bot/bot.module';
import {SoundModule} from './sound/sound.module';
import {StorageModule} from './storage/storage.module';
import {TypeOrmModule} from '@nestjs/typeorm';
import {UserModule} from './user/user.module';
import {AuthModule} from './auth/auth.module';

@Module({
  imports: [TypeOrmModule.forRoot(), BotModule, SoundModule, StorageModule, UserModule, AuthModule],
  controllers: [AppController],
  providers: [AppService],
})
export class AppModule {
}
