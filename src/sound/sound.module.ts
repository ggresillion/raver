import { Module } from '@nestjs/common';
import { SoundService } from './sound.service';
import { SoundController } from './sound.controller';
import {StorageModule} from '../storage/storage.module';

@Module({
  imports: [StorageModule],
  providers: [SoundService],
  controllers: [SoundController],
})
export class SoundModule {}
