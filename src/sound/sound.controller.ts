import {Controller, Get} from '@nestjs/common';
import {SoundService} from './sound.service';

@Controller('sounds')
export class SoundController {

  constructor(
    private readonly soundService: SoundService,
  ) {
  }

  @Get()
  public async getSounds(): Promise<string[]> {
    return await this.soundService.getSounds();
  }
}
