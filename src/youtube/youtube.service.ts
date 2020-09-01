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
import { GuildDTO } from '../guild/dto/guild.dto';

@Injectable()
export class YoutubeService {

  private playlist: Map<string, Video[]> = new Map();
  private status: Map<string, PlayerStatus> = new Map();
  private shouldNotPlayNext: Map<string, boolean> = new Map();
  private totalLengthSeconds: number;

  constructor(
    private readonly soundService: SoundService,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
    @Inject(forwardRef(() => YoutubeGateway))
    private readonly youtubeGateway: YoutubeGateway,
  ) {
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
          limit: 10,
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
    const res = await ytdlDiscord.stream(link);
    this.totalLengthSeconds = res.totalLengthSeconds;
    this.botService.playFromStream(
      guildId,
      res.stream,
      () => {
        this.status.set(guildId, PlayerStatus.PLAYING);
        this.propagateState(guildId);
      },
      (progress) => {
        this.youtubeGateway.sendProgressUpdate(guildId, progress / 1000);
      },
      () => {
        if (!this.shouldNotPlayNext.get(guildId)) {
          this.playlist.get(guildId).splice(0, 1);
          if (this.playlist.get(guildId).length === 0) {
            return;
          }
          this.status.set(guildId, PlayerStatus.LOADING);
          this.youtubeGateway.sendStatusUpdate(guildId, PlayerStatus.LOADING);
          this.playSoundFromYoutube(guildId, this.playlist.get(guildId)[0].link);
          this.propagateState(guildId);
        } else {
          this.status.set(guildId, PlayerStatus.PAUSED);
          this.youtubeGateway.sendStatusUpdate(guildId, PlayerStatus.PAUSED);
        }
        this.shouldNotPlayNext.delete(guildId);
      });
  }

  public getPlaylist(guildId: string) {
    return this.playlist.get(guildId);
  }

  public getState(guildId: string): PlayerState {
    return {
      status: this.status.get(guildId),
      playlist: this.playlist.get(guildId),
      totalLengthSeconds: this.totalLengthSeconds,
    };
  }

  public play(guildId: string) {
    this.setStatus(guildId, PlayerStatus.LOADING);
    if (this.botService.isPlaying(guildId)) {
      this.botService.resumeStream(guildId);
      this.setStatus(guildId, PlayerStatus.PLAYING);
    } else {
      this.playSoundFromYoutube(guildId, this.playlist.get(guildId)[0].link);
    }
  }

  public pause(guildId: string) {
    this.botService.pauseStream(guildId);
    this.status.set(guildId, PlayerStatus.PAUSED);
    this.youtubeGateway.sendStatusUpdate(guildId, PlayerStatus.PAUSED);
  }

  public stop(guildId: string) {
    this.shouldNotPlayNext.set(guildId, true);
    this.botService.stopStream(guildId);
  }

  public next(guildId: string) {
    if (!this.playlist.has(guildId)) {
      return;
    }
    if (this.playlist.get(guildId).length <= 1) {
      return;
    }
    this.botService.stopStream(guildId);
  }

  public moveUpwards(guildId: string, index: number): void {
    const playlist = this.playlist.get(guildId);
    if (!playlist) {
      return;
    }
    if (playlist.length < index || index < 2) {
      return;
    }
    playlist.splice(index - 1, 0, playlist.splice(index, 1)[0]);
    this.playlist.set(guildId, playlist);
    this.propagateState(guildId);
  }

  public moveDownwards(guildId: string, index: number): void {
    const playlist = this.playlist.get(guildId);
    if (!playlist) {
      return;
    }
    if (playlist.length <= index) {
      return;
    }
    playlist.splice(index + 1, 0, playlist.splice(index, 1)[0]);
    this.playlist.set(guildId, playlist);
    this.propagateState(guildId);
  }

  public removeFromPlaylist(guildId: string, index: number): void {
    const playlist = this.playlist.get(guildId);
    if (!playlist) {
      return;
    }
    if (playlist.length <= index || index === 0) {
      return;
    }
    playlist.splice(index, 1);
    this.playlist.set(guildId, playlist);
    this.propagateState(guildId);
  }

  private onBotStatusUpdate(guildId: string, status: BotStatus) {
    switch (status) {
      case BotStatus.IN_VOICE_CHANNEL:
        if (this.status.get(guildId) !== PlayerStatus.NA) {
          return;
        }
        this.setStatus(guildId, PlayerStatus.IDLE);
        break;
      case BotStatus.CONNECTED:
        this.setStatus(guildId, PlayerStatus.NA);
        break;
      case BotStatus.DISCONNECTED:
        this.setStatus(guildId, PlayerStatus.NA);
        break;
    }
  }

  private setStatus(guildId: string, status: PlayerStatus): void {
    this.status.set(guildId, status);
    this.youtubeGateway.sendStatusUpdate(guildId, this.status.get(guildId));
  }

  private onAddToPlaylist(guildId: string, track: Video) {
    if (!this.playlist.has(guildId)) {
      this.playlist.set(guildId, []);
    }
    this.playlist.get(guildId).push(track);
    if (this.status.get(guildId) === PlayerStatus.IDLE) {
      this.play(guildId);
    }
    this.propagateState(guildId);
  }

  private propagateState(guildId: string) {
    this.youtubeGateway.sendStateUpdate(guildId, {
      status: this.status.get(guildId) || PlayerStatus.NA,
      playlist: this.playlist.get(guildId),
      totalLengthSeconds: this.totalLengthSeconds,
    });
  }
}
