import {forwardRef, Inject, Injectable} from '@nestjs/common';
import {StorageService} from '../storage/storage.service';
import {SoundService} from '../sound/sound.service';
import {BotService, BotStatus} from '../bot/bot.service';
import {YoutubeGateway} from './youtube-gateway';
import {PlayerStatus} from './model/player-status';
import * as ytdlDiscord from 'ytdl-core-discord';
import {TrackInfos} from './dto/track-infos';
import {YoutubeDataAPI} from 'youtube-v3-api';
import * as ytdl from 'ytdl-core';
import * as FFmpeg from 'fluent-ffmpeg';

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
    const res = await this.ytApi.searchVideo(id) as any;
    return {
      id: res.items[0].id,
      title: res.items[0].snippet.title,
      description: res.items[0].snippet.description,
    };

  }

  public async searchVideos(q: string): Promise<TrackInfos[]> {
    const rawVideos = (await this.ytApi.searchAll(q, 10, {part: 'id'}) as any).items;
    return rawVideos.map(rawVideo => ({id: rawVideo.id.videoId, ...rawVideo.snippet}));
  }

  public async uploadFromYoutube(url: string, name: string, categoryId: number): Promise<any> {
    return new Promise((resolve, reject) => {
      const sound = this.soundService.createNewSoundEntity(name, categoryId);
      const stream = ytdl(url, {quality: 'lowestaudio'});
      FFmpeg({source: stream})
        .audioBitrate(128)
        .withNoVideo()
        .toFormat('opus')
        .save(this.storageService.getUploadDir() + '/' + sound.uuid)
        .on('end', () => {
          resolve(this.soundService.saveNewSoundEntity(sound));
        });
    });
  }

  public async playSoundFromYoutube(id: string) {
    ytdl.getBasicInfo(`http://youtube.com/watch?v=${id}`, (err, info) => {
      console.log(info)
    });
    const stream = await ytdlDiscord(
      `http://youtube.com/watch?v=${id}`,
      {highWaterMark: 1024 * 1024 * 10});
    stream.on('data', d => {
      console.log(d);
    });
    this.botService.playFromStream(stream,
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
    this.status = PlayerStatus.PLAYING;
    this.propagateState();
  }

  public pause() {
    this.botService.pauseStream();
    this.status = PlayerStatus.PAUSED;
    this.youtubeGateway.sendStatusUpdate(this.status);
    this.propagateState();
  }

  public stop() {
    this.botService.stopStream();
    this.propagateState();
  }

  public next() {
    this.botService.stopStream();
    this.playlist.splice(0, 1);
    this.playSoundFromYoutube(this.playlist[0].id);
    this.propagateState();
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
    if (this.getStatus() === PlayerStatus.IDLE) {
      this.play();
    }
    this.propagateState();
  }

  private propagateState() {
    this.youtubeGateway.sendStateUpdate({status: this.status, playlist: this.playlist});
  }
}
