import {Injectable} from '@nestjs/common';
import {exec, getInfo, Info} from 'youtube-dl';
import {StorageService} from '../storage/storage.service';
import {SoundService} from '../sound/sound.service';
import {BotService} from '../bot/bot.service';
import youtubedl = require('youtube-dl');

@Injectable()
export class YoutubeService {

  constructor(
    private readonly soundService: SoundService,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
  ) {
  }

  public async searchYoutube(url: string): Promise<Info> {
    return new Promise((resolve, reject) =>
      getInfo(url, (err, info) => {
        if (err) {
          reject(err);
        } else {
          resolve(info);
        }
      }));
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

  public async playSoundFromYoutube(url: string) {
    this.botService.playFromStream(youtubedl(url, ['-f worstaudio'], {}));
  }
}
