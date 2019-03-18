import {Controller, FileInterceptor, Get, Param, Post, UploadedFile, UseGuards, UseInterceptors} from '@nestjs/common';
import {SoundService} from './sound.service';
import {Sound} from './sound.entity';
import {AuthGuard} from '@nestjs/passport';

@Controller('sounds')
@UseGuards(AuthGuard())
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
  public async createSound(@UploadedFile()file): Promise<Sound> {
    return await this.soundService.saveSound(file.originalname, file.buffer);
  }

  @Post('/:id/play')
  public async playSound(@Param()id: number) {
    return await this.soundService.playSound(id);
  }
}
