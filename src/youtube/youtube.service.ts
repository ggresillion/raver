import { forwardRef, Inject, Injectable } from '@nestjs/common';
import { StorageService } from '../storage/storage.service';
import { SoundService } from '../sound/sound.service';
import { BotService } from '../bot/bot.service';
import { YoutubeGateway } from './youtube-gateway';
import { PlayerStatus } from './model/player-status';
import * as ytdlDiscord from './util/ytdl-wrapper';
import * as ytdl from 'ytdl-core';
import * as FFmpeg from 'fluent-ffmpeg';
import { PlayerState } from './model/player-state';
import * as youtubeSearch from 'ytsr';
import { Sound } from '../sound/entity/sound.entity';
import { BotStatus } from '../bot/dto/bot-status.enum';
import { UploadDto } from './dto/upload.dto';
import { Bucket } from '../storage/bucket.enum';
import { Video } from 'ytsr';

@Injectable()
export class YoutubeService {

  private playlist: Video[] = [];
  private status: Map<string, PlayerStatus> = new Map();
  private totalLengthSeconds: number;

  constructor(
    private readonly soundService: SoundService,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
    @Inject(forwardRef(() => YoutubeGateway))
    private readonly youtubeGateway: YoutubeGateway,
  ) {
    // const botStatus = this.botService.getInfos().status;
    // this.status = botStatus === BotStatus.CONNECTED || botStatus === BotStatus.IN_VOICE_CHANNEL
    //   ? PlayerStatus.IDLE : PlayerStatus.NA;
    this.botService.onBotStatusUpdate((guildId, status) => this.onBotStatusUpdate(guildId, status));
    this.youtubeGateway.onAddToPlaylist((guildid, track) => this.onAddToPlaylist(guildid, track));
  }

  public async getVideoInfos(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      youtubeSearch(id, (err, result) => {
        if (err) {
          return reject(err);
        }
        if (result.items.length === 0) {
          return reject('did not find any video');
        }
        return resolve(result.items[0]);
      });
    });
  }

  public async searchVideos(q: string): Promise<Video[]> {
    return youtubeSearch.getFilters(q)
      .then(filters => {
        const filter = filters.get('Type').find(o => o.name === 'Video');
        var options = {
          limit: 5,
          nextpageRef: filter.ref,
        }
        return youtubeSearch(null, options);
      }
      ).then(res => <Video[]>res.items);
  }

  public async uploadFromYoutube(upload: UploadDto): Promise<Sound> {
    return new Promise((resolve) => {
      const sound = this.soundService.createNewSoundEntity(upload.name, upload.categoryId, upload.guildId);
      const stream = ytdl(upload.url, { quality: 'lowestaudio' });
      FFmpeg({ source: stream })
        .audioBitrate(128)
        .withNoVideo()
        .toFormat('opus')
        .save(this.storageService.getUploadDir(Bucket.SOUNDS) + '/' + sound.uuid)
        .on('end', () => {
          resolve(this.soundService.saveSoundEntity(sound));
        });
    });
  }

  public async playSoundFromYoutube(guildId: string, link: string) {
    const res = await ytdlDiscord.stream(
      link,
      { highWaterMark: 1024 * 1024 * 10 },
      ((current) => {
        this.youtubeGateway.sendProgressUpdate(guildId, current);
      }));
    this.totalLengthSeconds = res.totalLengthSeconds;
    this.botService.playFromStream(
      guildId,
      res.stream,
      () => {
        this.status.set(guildId, PlayerStatus.PLAYING);
        this.propagateState(guildId);
      },
      () => {
        this.status.set(guildId, PlayerStatus.IDLE);
        this.propagateState(guildId);
      });
  }

  public getPlaylist() {
    return this.playlist;
  }

  public getState(guildId: string): PlayerState {
    return {
      status: this.status.get(guildId),
      playlist: this.playlist,
      totalLengthSeconds: this.totalLengthSeconds,
    };
  }

  public play(guildId: string) {
    if (this.botService.isPlaying(guildId)) {
      this.botService.resumeStream(guildId);
    } else {
      this.playSoundFromYoutube(guildId, this.playlist[0].link);
    }
    this.status.set(guildId, PlayerStatus.PLAYING);
    this.propagateState(guildId);
  }

  public pause(guildId: string) {
    this.botService.pauseStream(guildId);
    this.status.set(guildId, PlayerStatus.PAUSED);
    this.youtubeGateway.sendStatusUpdate(guildId, PlayerStatus.PAUSED);
    this.propagateState(guildId);
  }

  public stop(guildId: string) {
    this.botService.stopStream(guildId);
    this.propagateState(guildId);
  }

  public next(guildId: string) {
    this.botService.stopStream(guildId);
    this.playlist.splice(0, 1);
    this.playSoundFromYoutube(guildId, this.playlist[0].link);
    this.propagateState(guildId);
  }

  private onBotStatusUpdate(guildId: string, status: BotStatus) {
    switch (status) {
      case BotStatus.IN_VOICE_CHANNEL:
        this.status.set(guildId, PlayerStatus.IDLE);
        break;
      case BotStatus.CONNECTED:
        this.status.set(guildId, PlayerStatus.NA);
        break;
      case BotStatus.DISCONNECTED:
        this.status.set(guildId, PlayerStatus.NA);
        break;
    }
    this.youtubeGateway.sendStatusUpdate(guildId, this.status.get(guildId));
  }

  private onAddToPlaylist(guildId: string, track: Video) {
    this.playlist.push(track);
    if (this.status.get(guildId) === PlayerStatus.IDLE) {
      this.play(guildId);
    }
    this.propagateState(guildId);
  }

  private propagateState(guildId: string) {
    this.youtubeGateway.sendStateUpdate(guildId, {
      status: this.status.get(guildId) || PlayerStatus.NA,
      playlist: this.playlist,
      totalLengthSeconds: this.totalLengthSeconds,
    });
  }
}
