import {Injectable} from '@nestjs/common';
import {getInfo, Info} from 'youtube-dl';

@Injectable()
export class YoutubeService {

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
}
