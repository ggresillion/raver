import {BadRequestException, Controller, Get, Post, Query} from '@nestjs/common';
import {YoutubeService} from './youtube.service';
import {Body} from '@nestjs/common/decorators/http/route-params.decorator';
import {UploadDto} from './dto/upload.dto';
import {TrackInfos} from './dto/track-infos';
import {Sound} from '../sound/sound.entity';
import {YouTubeSearchResults} from 'youtube-search';

@Controller('youtube')
export class YoutubeController {

  constructor(
    private readonly youtubeService: YoutubeService,
  ) {
  }

  @Get('search')
  public async searchOnYoutube(@Query('q') q: string): Promise<YouTubeSearchResults[]> {
    return await this.youtubeService.searchVideos(q);
  }

  @Post('upload')
  public async uploadFromYoutube(@Body()upload: UploadDto): Promise<Sound> {
    return await this.youtubeService.uploadFromYoutube(upload.url, upload.name, upload.categoryId);
  }

  @Post('play')
  public async playSoundFromYoutube(@Query('url') url: string) {
    if (!url || url === '') {
      throw new BadRequestException();
    }
    return await this.youtubeService.playSoundFromYoutube(url);
  }
}
