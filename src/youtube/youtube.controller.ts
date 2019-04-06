import {Controller, Get, Query} from '@nestjs/common';
import {YoutubeService} from './youtube.service';
import {Info} from 'youtube-dl';

@Controller('youtube')
export class YoutubeController {

  constructor(
    private readonly youtubeService: YoutubeService,
  ) {
  }

  @Get('search')
  public async getFromYoutube(@Query('url') url: string): Promise<Info> {
    return await this.youtubeService.searchYoutube(url);
  }
}
