import {BadRequestException, Body, Controller, Get, Param, Post, Put, UploadedFile, UseGuards, UseInterceptors} from '@nestjs/common';
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
    return await this.soundService.saveSound(file.originalname, file.buffer);
  }

  @Put('/:id')
  public async editSound(@Param('id')id: number, @Body() soundData: SoundDto): Promise<Sound> {
    return await this.soundService.editSound(id, soundData);
  }

  @Post('/:id/play')
  public async playSound(@Param()id: ObjectID) {
    return await this.soundService.playSound(id);
  }
}
