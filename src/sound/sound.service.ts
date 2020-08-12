import { Injectable, NotFoundException, OnApplicationShutdown, UnprocessableEntityException } from '@nestjs/common';
import { StorageService } from '../storage/storage.service';
import { ObjectID, Repository } from 'typeorm';
import { Sound } from './entity/sound.entity';
import { InjectRepository } from '@nestjs/typeorm';
import { BotService } from '../bot/bot.service';
import { SoundDto } from './dto/sound.dto';
import { Image } from './entity/image.entity';
import { OpusEncoder } from 'node-opus';
import * as FFmpeg from 'fluent-ffmpeg';
import { Duplex } from 'stream';
import uuid = require('uuid');

@Injectable()
export class SoundService {

  constructor(
    @InjectRepository(Sound)
    private readonly soundRepository: Repository<Sound>,
    @InjectRepository(Image)
    private readonly imageRepository: Repository<Image>,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
  ) {
  }

  public async getSounds(guildId: string): Promise<Sound[]> {
    return this.soundRepository.find({ where: { guildId }, relations: ['category'] });
  }

  public async getSoundById(id: ObjectID): Promise<Sound> {
    return this.soundRepository.findOne(id, { relations: ['category'] });
  }

  public async saveSound(name: string, categoryId: number, guildId: string, bSound: Buffer, bImage: Buffer) {
    try {
      let sound;
      if (!!bImage) {
        const image = await this.imageRepository.save(this.imageRepository.create());
        sound = this.soundRepository.create({ name, categoryId, guildId, image });
        await this.storageService.saveFile(image.uuid, bImage);
      } else {
        sound = this.soundRepository.create({ name, categoryId, guildId });
      }
      // await this.storageService.saveFile(sound.uuid, bSound);
      await this.saveToOpus(bSound, sound.uuid);

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
      this.botService.playFile(sound.uuid, sound.guildId);
      return sound;
    });
  }

  public createNewSoundEntity(name: string, categoryId: number, guildId: string): Sound {
    return this.soundRepository.create({ name, categoryId, guildId });
  }

  public async saveSoundEntity(sound: Sound): Promise<Sound> {
    return await this.soundRepository.save(sound);
  }

  public async deleteSound(soundId: number) {
    const sound = await this.soundRepository.findOne(soundId);
    if (!sound) {
      throw new NotFoundException('sound not found');
    }
    await this.storageService.removeFile(sound.uuid);
    return await this.soundRepository.remove(sound);
  }

  private async saveToOpus(buffer: Buffer, uuid: string): Promise<void> {

    await this.storageService.saveFile('tmp/' + uuid, buffer);

    await new Promise((resolve, reject) => {
      FFmpeg({ source: './uploads/tmp/' + uuid })
        .audioBitrate(128)
        .withNoVideo()
        .toFormat('opus')
        .save(this.storageService.getUploadDir() + '/' + uuid)
        .on('error', (err) => {
          reject(err);
        })
        .on('end', () => {
          resolve();
        });
    });

    await this.storageService.removeFile('tmp/' + uuid);
  }
}
