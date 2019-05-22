import {HttpService, Injectable} from '@nestjs/common';
import {exec} from 'youtube-dl';
import {StorageService} from '../storage/storage.service';
import {SoundService} from '../sound/sound.service';
import {BotService} from '../bot/bot.service';
import {YoutubeGateway} from './youtube-gateway';
import {Status} from './model/status';
import * as ytdl from 'ytdl-core';
import {TrackInfos} from './dto/track-infos';

@Injectable()
export class YoutubeService {

  constructor(
    private readonly soundService: SoundService,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
    private readonly youtubeGateway: YoutubeGateway,
    private readonly httpService: HttpService,
  ) {
  }

  public async searchYoutube(url: string): Promise<TrackInfos> {
    const infos = (await this.httpService.get<any>(`https://noembed.com/embed?url=${url}`).toPromise()).data;
    return ({url: infos.url, title: infos.title, thumbnail: infos.thumbnail_url});
  }

  public async uploadFromYoutube(url: string, name: string, categoryId: number): Promise<any> {
    return new Promise((resolve, reject) => {
      const sound = this.soundService.createNewSoundEntity(name, categoryId);
      exec(
        url,
        ['-x', '--audio-format', 'mp3', '-o', sound.uuid],
        {cwd: this.storageService.getUploadDir()},
        async (err) => {
          if (err) {
            return reject(err);
          }
          resolve(this.soundService.saveNewSoundEntity(sound));
        },
      );
    });
  }

  public async playSoundFromYoutube(url: string): Promise<TrackInfos> {
    const infos = await this.searchYoutube(url);
    this.botService.playFromStream(ytdl(
      url,
      {filter: 'audioonly'}),
      () => this.youtubeGateway.sendStatusUpdate({status: Status.PLAYING, track: infos}),
      () => this.youtubeGateway.sendStatusUpdate({status: Status.PAUSED, track: infos}));
    return infos;
  }
}
