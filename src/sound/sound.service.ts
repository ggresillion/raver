import {Injectable, NotFoundException, OnApplicationShutdown, UnprocessableEntityException} from '@nestjs/common';
import {StorageService} from '../storage/storage.service';
import {ObjectID, Repository} from 'typeorm';
import {Sound} from './sound.entity';
import {InjectRepository} from '@nestjs/typeorm';
import {BotService} from '../bot/bot.service';
import {SoundDto} from './sound.dto';

@Injectable()
export class SoundService {

  constructor(
    @InjectRepository(Sound)
    private readonly soundRepository: Repository<Sound>,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
  ) {
  }

  public async getSounds(): Promise<Sound[]> {
    return this.soundRepository.find({relations: ['category']});
  }

  public async getSoundById(id: ObjectID): Promise<Sound> {
    return this.soundRepository.findOne(id, {relations: ['category']});
  }

  public async saveSound(name: string, buffer: Buffer) {
    try {
      const sound = this.soundRepository.create({name});
      await this.storageService.saveFile(sound.uuid, buffer);
      return await this.soundRepository.save(sound);
    } catch (e) {
      if (e.code === 11000) {
        throw new UnprocessableEntityException('name already exists');
      }
      throw e;
    }
  }

  public async editSound(id: number, data: SoundDto) {
    const sound = await this.soundRepository.findOne(id);
    if (!sound) {
      throw new NotFoundException('sound not found');
    }
    const updatedSound = Object.assign(sound, data);
    return await this.soundRepository.save(updatedSound);
  }

  public async playSound(id: ObjectID): Promise<Sound> {
    return this.getSoundById(id).then(sound => {
      if (!sound) {
        throw new NotFoundException('sound not found');
      }
      this.botService.playFile(sound.uuid);
      return sound;
    });
  }

  public createNewSoundEntity(name: string, categoryId: number): Sound {
    return this.soundRepository.create({name, categoryId});
  }

  public async saveNewSoundEntity(sound: Sound): Promise<Sound> {
    return await this.soundRepository.save(sound);
  }
}
