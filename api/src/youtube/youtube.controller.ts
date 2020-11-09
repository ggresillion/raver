import {BadRequestException, Controller, Get, Post, Query} from '@nestjs/common';
import {YoutubeService} from './youtube.service';
import {Body} from '@nestjs/common/decorators/http/route-params.decorator';
import {UploadDto} from './dto/upload.dto';
import { Sound } from '../sound/entity/sound.entity';

@Controller('youtube')
export class YoutubeController {

  constructor(
    private readonly youtubeService: YoutubeService,
  ) {
  }

  @Get('infos')
  public async getVideoInfos(@Query('url') url: string): Promise<any> {
    const id = new URLSearchParams(url.split('?')[1]).get('v');
    return await this.youtubeService.getVideoInfos(id);
  }

  @Get('search')
  public async searchOnYoutube(@Query('q') q: string): Promise<any[]> {
    return await this.youtubeService.searchVideos(q);
  }

  @Post('upload')
  public async uploadFromYoutube(@Body()upload: UploadDto): Promise<Sound> {
    return await this.youtubeService.uploadFromYoutube(upload);
  }

  @Post('play')
  public async playSoundFromYoutube(@Query('id') id: string) {
    if (!id || id === '') {
      throw new BadRequestException();
    }
    return await this.youtubeService.playSoundFromYoutube("guildId", id);
  }
}
