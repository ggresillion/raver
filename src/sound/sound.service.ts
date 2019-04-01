import {Injectable, NotFoundException, UnprocessableEntityException} from '@nestjs/common';
import {StorageService} from '../storage/storage.service';
import {MongoRepository, ObjectID} from 'typeorm';
import {Sound} from './sound.entity';
import {InjectRepository} from '@nestjs/typeorm';
import {BotService} from '../bot/bot.service';

@Injectable()
export class SoundService {

  constructor(
    private readonly storageService: StorageService,
    private readonly botService: BotService,
    @InjectRepository(Sound)
    private readonly soundRepository: MongoRepository<Sound>,
  ) {
  }

  public async getSounds(): Promise<Sound[]> {
    return this.soundRepository.find();
  }

  public async saveSound(name: string, buffer: Buffer) {
    try {
      const sound = await this.soundRepository.save(this.soundRepository.create({name}));
      await this.storageService.saveFile(sound.uuid, buffer);
      return sound;
    } catch (e) {
      if (e.code === 11000) {
        throw new UnprocessableEntityException('name already exists');
      }
      throw e;
    }
  }

  public async playSound(id: number) {
    return this.botService.playSound(id);
  }
}
