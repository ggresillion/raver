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
import { Sound } from '../sound/entity/sound.entity';
import { BotStatus } from '../bot/dto/bot-status.enum';
import { UploadDto } from './dto/upload.dto';
import { Bucket } from '../storage/bucket.enum';
import youtube from './scraper';
import { TrackInfos } from './dto/track-infos';
import { Video } from 'scrape-youtube/lib/interface';

@Injectable()
export class YoutubeService {

  private playlist: Map<string, TrackInfos[]> = new Map();
  private status: Map<string, PlayerStatus> = new Map();
  private totalLengthSeconds: number;

  constructor(
    private readonly soundService: SoundService,
    private readonly storageService: StorageService,
    private readonly botService: BotService,
    @Inject(forwardRef(() => YoutubeGateway))
    private readonly youtubeGateway: YoutubeGateway
  ) {
    this.botService.onBotStatusUpdate((guildId, status) => this.onBotStatusUpdate(guildId, status));
    this.youtubeGateway.onAddToPlaylist((guildid, track) => this.onAddToPlaylist(guildid, track));
  }

  public async getVideoInfos(id: string): Promise<TrackInfos> {
    const v = (await youtube.search(id)).videos[0];
    return ({
      title: v.title,
      link: v.link,
      thumbnail: v.thumbnail,
      author: {
        name: v.channel.name,
        ref: v.channel.link,
        verified: v.channel.verified
      },
      description: v.description,
      views: v.views,
      duration: v.duration
    });
  }

  public async searchVideos(q: string): Promise<TrackInfos[]> {
    const search = await youtube.search(q);
    return search.videos
      .map((v: Video) => ({
        title: v.title,
        link: v.link,
        thumbnail: v.thumbnail,
        author: {
          name: v.channel.name,
          ref: v.channel.link,
          verified: v.channel.verified
        },
        description: v.description,
        views: v.views,
        duration: v.duration
      }));
  }

  public async uploadFromYoutube(upload: UploadDto): Promise<Sound> {
    const sound = await new Promise<Sound>((resolve, reject) => {
      const sound = this.soundService.createNewSoundEntity(upload.name, upload.categoryId, upload.guildId);
      const stream = ytdl(upload.url, { quality: 'lowestaudio' });
      FFmpeg({ source: stream })
        .audioBitrate(128)
        .withNoVideo()
        .toFormat('opus')
        .save(this.storageService.getUploadDir(Bucket.SOUNDS) + '/' + sound.uuid)
        .on('end', () => {
          resolve(this.soundService.saveSoundEntity(sound));
        })
        .on('error', err => reject(err));
    });

    return sound;
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
      (stop) => {
        if (stop) {
          this.status.set(guildId, PlayerStatus.IDLE);
          this.youtubeGateway.sendStatusUpdate(guildId, PlayerStatus.IDLE);
        } else {
          this.playlist.get(guildId).splice(0, 1);
          if (this.playlist.get(guildId).length === 0) {
            return;
          }
          this.status.set(guildId, PlayerStatus.LOADING);
          this.youtubeGateway.sendStatusUpdate(guildId, PlayerStatus.LOADING);
          this.playSoundFromYoutube(guildId, this.playlist.get(guildId)[0].link);
          this.propagateState(guildId);
        }
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
    this.playlist.get(guildId).splice(0, 1);
    this.propagateState(guildId);
    this.play(guildId);
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

  private onAddToPlaylist(guildId: string, track: TrackInfos) {
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
