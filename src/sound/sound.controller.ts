import {BadRequestException, Body, Controller, Delete, Get, Param, Post, Put, UploadedFile, UseGuards, UseInterceptors} from '@nestjs/common';
import {SoundService} from './sound.service';
import {Sound} from './sound.entity';
import {UserGuard} from 'src/auth/guards/user.guard';
import {ObjectID, UpdateResult} from 'typeorm';
import {FileInterceptor} from '@nestjs/platform-express';
import {SoundDto} from './sound.dto';

@Controller('sounds')
export class SoundController {

  constructor(
    private readonly soundService: SoundService,
  ) {
  }

  @Get()
  public async getSounds(): Promise<Sound[]> {
    return await this.soundService.getSounds();
  }

  @Post()
  @UseInterceptors(FileInterceptor('sound'))
  public async createSound(@UploadedFile()file, @Body() soundData: SoundDto): Promise<Sound> {
    if (!file) {
      throw new BadRequestException('missing file');
    }
    const name = soundData.name && soundData.name !== '' ? soundData.name : file.originalname;
    return await this.soundService.saveSound(name, soundData.categoryId, file.buffer);
  }

  @Delete(':id')
  public async deleteSound(@Param('id')id: number): Promise<Sound> {
    return await this.soundService.deleteSound(id);
  }

  @Put(':id')
  public async editSound(@Param('id')id: number, @Body() soundData: SoundDto): Promise<Sound> {
    return await this.soundService.editSound(id, soundData);
  }

  @Post(':id/play')
  public async playSound(@Param()id: ObjectID) {
    return await this.soundService.playSound(id);
  }
}
