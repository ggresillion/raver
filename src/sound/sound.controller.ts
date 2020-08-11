import {
  BadRequestException,
  Body,
  Controller,
  Delete,
  Get,
  Param,
  Post,
  Put,
  Query,
  UploadedFile,
  UploadedFiles,
  UseInterceptors
} from '@nestjs/common';
import { SoundService } from './sound.service';
import { Sound } from './entity/sound.entity';
import { ObjectID } from 'typeorm';
import { FileFieldsInterceptor } from '@nestjs/platform-express';
import { SoundDto } from './dto/sound.dto';

@Controller('sounds')
export class SoundController {

  constructor(
    private readonly soundService: SoundService,
  ) {
  }

  @Get()
  public async getSounds(@Query('guildId') guildId: string): Promise<Sound[]> {
    return await this.soundService.getSounds(guildId);
  }

  @Post()
  @UseInterceptors(FileFieldsInterceptor([
    { name: 'sound', maxCount: 1 },
    { name: 'image', maxCount: 1 },
  ]))
  public async createSound(@UploadedFiles() files,
    @Body() soundData: SoundDto): Promise<Sound> {
    console.log(files);
    const sound = files.sound ? files.sound[0] : null;
    const image = files.image ? files.image[0] : null;
    if (!sound) {
      throw new BadRequestException('missing file');
    }
    const name = soundData.name && soundData.name !== '' ? soundData.name : sound.originalname;
    return await this.soundService.saveSound(name, soundData.categoryId, soundData.guildId, sound.buffer, image ? image.buffer : null);
  }

  @Delete(':id')
  public async deleteSound(@Param('id') id: number): Promise<Sound> {
    return await this.soundService.deleteSound(id);
  }

  @Put(':id')
  public async editSound(@Param('id') id: number, @Body() soundData: SoundDto): Promise<Sound> {
    return await this.soundService.editSound(id, soundData);
  }

  @Post(':id/play')
  public async playSound(@Param('id') id: ObjectID) {
    return await this.soundService.playSound(id);
  }
}
