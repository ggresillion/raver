import {BadRequestException, Controller, Get, Post, Query} from '@nestjs/common';
import {YoutubeService} from './youtube.service';
import {Info} from 'youtube-dl';
import {Body} from '@nestjs/common/decorators/http/route-params.decorator';
import {UploadDto} from './dto/upload.dto';

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

  @Post('upload')
  public async uploadFromYoutube(@Body()upload: UploadDto): Promise<any> {
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
