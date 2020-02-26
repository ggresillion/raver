import {forwardRef, Inject, Injectable} from '@nestjs/common';
// import {exec} from 'youtube-dl';
import {StorageService} from '../storage/storage.service';
import {SoundService} from '../sound/sound.service';
import {BotService, BotStatus} from '../bot/bot.service';
import {YoutubeGateway} from './youtube-gateway';
import {PlayerStatus} from './model/player-status';
import * as ytdlDiscord from 'ytdl-core-discord';
import {TrackInfos} from './dto/track-infos';
import {YoutubeDataAPI} from 'youtube-v3-api';

@Injectable()
export class YoutubeService {

  private ytApi = new YoutubeDataAPI(process.env.YOUTUBE_API_KEY);
  private playlist: TrackInfos[] = [];
  private status: PlayerStatus;

  constructor(
    private readonly soundService: SoundService,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
    @Inject(forwardRef(() => YoutubeGateway))
    private readonly youtubeGateway: YoutubeGateway,
  ) {
    const botStatus = this.botService.getInfos().status;
    this.status = botStatus === BotStatus.CONNECTED || botStatus === BotStatus.IN_VOICE_CHANNEL
      ? PlayerStatus.IDLE : PlayerStatus.NA;
    this.botService.onBotStatusUpdate(status => this.onBotStatusUpdate(status));
    this.youtubeGateway.onAddToPlaylist(track => this.onAddToPlaylist(track));
  }

  public async fetchVideoInfos(id: string): Promise<TrackInfos> {
    return await this.ytApi.searchVideo(id, 10) as TrackInfos;

  }

  public async searchVideos(q: string): Promise<TrackInfos[]> {
    const rawVideos = (await this.ytApi.searchAll(q, 10, {part: 'id'}) as any).items;
    return rawVideos.map(rawVideo => ({id: rawVideo.id.videoId, ...rawVideo.snippet}));
  }

  public async uploadFromYoutube(url: string, name: string, categoryId: number): Promise<any> {
    return new Promise((resolve, reject) => {
      // const sound = this.soundService.createNewSoundEntity(name, categoryId);
      // youtubeDl.exec(
      //   url,
      //   ['-x', '--audio-format', 'mp3', '-o', sound.uuid],
      //   { cwd: this.storageService.getUploadDir() },
      //   async (err) => {
      //     if (err) {
      //       return reject(err);
      //     }
      //     resolve(this.soundService.saveNewSoundEntity(sound));
      //   },
      // );
    });
  }

  public async playSoundFromYoutube(id: string) {
    this.botService.playFromStream(await ytdlDiscord(
      `http://youtube.com/watch?v=${id}`,
      {highWaterMark: 1024 * 1024 * 10}),
      () => this.youtubeGateway.sendStatusUpdate(PlayerStatus.PLAYING),
      () => this.youtubeGateway.sendStatusUpdate(PlayerStatus.PAUSED));
  }

  public getPlaylist() {
    return this.playlist;
  }

  public getStatus() {
    return this.status;
  }

  public play() {
    if (this.botService.isPlaying()) {
      this.botService.resumeStream();
    } else {
      this.playSoundFromYoutube(this.playlist[0].id);
    }
    this.youtubeGateway.sendStateUpdate({status: PlayerStatus.PLAYING, playlist: this.playlist});
  }

  public pause() {
    this.botService.pauseStream();
    this.youtubeGateway.sendStatusUpdate(PlayerStatus.PAUSED);
  }

  private onBotStatusUpdate(status: BotStatus) {
    switch (status) {
      case BotStatus.IN_VOICE_CHANNEL:
        this.status = PlayerStatus.IDLE;
        break;
      case BotStatus.CONNECTED:
        this.status = PlayerStatus.NA;
        break;
      case BotStatus.DISCONNECTED:
        this.status = PlayerStatus.NA;
        break;
    }
    this.youtubeGateway.sendStatusUpdate(this.status);
  }

  private onAddToPlaylist(track: TrackInfos) {
    this.playlist.push(track);
    this.youtubeGateway.sendAddToPlaylist(track);
  }
}
