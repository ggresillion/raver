import { Injectable, Logger, OnApplicationShutdown, UnprocessableEntityException } from '@nestjs/common';
import { Client, Message, StreamDispatcher, VoiceChannel, GuildMember, Snowflake, Collection, StreamType } from 'discord.js';
import { Command } from './command.enum';
import { StorageService } from '../storage/storage.service';
import { BotGateway } from './bot.gateway';
import { Readable, Stream } from 'stream';
import { UserDTO } from '../user/dto/user.dto';
import { GuildDTO } from '../guild/dto/guild.dto';
import { BotStatus } from './dto/bot-status.enum';
import { createReadStream } from 'fs';
import { Bucket } from '../storage/bucket.enum';

@Injectable()
export class BotService implements OnApplicationShutdown {

  private readonly PROGRESS_POLLING_INTERVAL = 500;
  private readonly DEFAULT_VOLUME = 0.5;

  private readonly logger = new Logger(BotService.name);
  private readonly token = process.env.BOT_TOKEN;

  private client: Client;
  private dispatchers: Map<string, StreamDispatcher> = new Map<string, StreamDispatcher>();
  private volumes: Map<string, number> = new Map();
  private onStatusChangeListeners: ((guildId: string, status: BotStatus) => void)[] = [];

  constructor(
    private readonly storageService: StorageService,
    private readonly botGateway: BotGateway,
  ) {
    this.client = new Client();
    this.connect();
    this.bindToEvents();
  }

  public onApplicationShutdown(signal?: string): any {
    this.client.voice.connections.forEach(co => co.disconnect());
  }

  public playFile(uuid: string, guildId: string) {
    const path = this.storageService.getPathFromUUID(Bucket.SOUNDS, uuid);
    this.play(guildId, createReadStream(path), { resumeOnInterupt: false, type: 'ogg/opus' },
      () => {
        this.logger.log(`Playing file ` + uuid);
      },
      () => { });
  }

  public playFromStream(guildId: string, stream: Readable, onStart: () => void, onProgress: (progress: number) => void, onEnd: (stop: boolean) => void) {

    let progressPolling;
    this.play(guildId, stream, { resumeOnInterupt: true, type: 'opus' },
      (dispatcher) => {
        progressPolling = setInterval(() => {
          onProgress(dispatcher.streamTime);
        }, this.PROGRESS_POLLING_INTERVAL);
        onStart();
      },
      (stop) => {
        onProgress(0);
        clearInterval(progressPolling);
        onEnd(stop);
      });
  }

  public pauseStream(guildId: string) {
    this.dispatchers.get(guildId).pause();
    this.logger.log('Paused stream');
  }

  public resumeStream(guildId: string) {
    this.dispatchers.get(guildId).resume();
    this.logger.log('Resumed stream');
  }

  public stopStream(guildId: string) {
    if (!this.dispatchers.has(guildId)) {
      return;
    }
    this.dispatchers.get(guildId).emit('stop');
    this.logger.log('Stopped stream');
  }

  public isPlaying(guildId: string): boolean {
    return this.dispatchers.has(guildId);
  }

  public onBotStatusUpdate(cb: (guildId: string, status: BotStatus) => void) {
    this.onStatusChangeListeners.push(cb);
  }

  public setIsBotInGuild(user: UserDTO, guilds: GuildDTO[]): GuildDTO[] {
    return guilds.map(guild => {
      return {
        ...guild,
        isBotInGuild: this.client.guilds.cache.array().some(g => g.id === guild.id),
      };
    });
  }

  public async joinMyChannel(user: UserDTO): Promise<void> {
    for (let channel of this.client.channels.cache.array()
      .filter(c => c instanceof VoiceChannel)) {
      const members = <Collection<Snowflake, GuildMember>>channel['members'];
      if (members.array().find(m => m.id === user.id)) {
        await this.joinChannel(<VoiceChannel>channel);
        return;
      }
    }
  }

  public getStatus(guildId: string): BotStatus {
    return this.client.voice.connections.some(c => c.channel.guild.id === guildId)
      ? BotStatus.IN_VOICE_CHANNEL
      : BotStatus.CONNECTED;
  }

  public setVolume(guildId: string, volume: number) {
    this.volumes.set(guildId, volume);
    if (this.dispatchers.has(guildId)) {
      this.dispatchers.get(guildId).setVolume(volume);
    }
    this.logger.debug(`(${guildId}) Changed volume to: ${volume}`)
    this.propagateState(guildId);
  }

  public getState(guildId: string) {
    return {
      guilds: this.client.guilds.cache.map(g => ({
        id: g.id,
        status: this.client.voice.connections.some(c => c.channel.guild.id === g.id)
          ? BotStatus.IN_VOICE_CHANNEL
          : BotStatus.CONNECTED,
      })),
      volume: this.getVolume(guildId)
    };
  }

  private play(guildId: string,
    stream: Readable,
    settings: { resumeOnInterupt: boolean, type: StreamType },
    onStart: (progress: StreamDispatcher) => void,
    onEnd: (stop?: boolean) => void) {
    const connection = this.client.voice.connections
      .find((_, id) => id === guildId);
    if (!connection) {
      throw new UnprocessableEntityException('guild id not found');
    }
    if (this.dispatchers.has(guildId)) {
      this.stopStream(guildId);
    }
    const dispatcher = connection.play(stream, { type: settings.type, volume: this.getVolume(guildId), highWaterMark: 50 });
    dispatcher.on('debug', debug => this.logger.debug(debug));
    dispatcher.on('error', error => this.logger.error(error));
    dispatcher.on('start', () => {
      this.logger.log('Playing from stream');
      onStart(dispatcher);
    });
    dispatcher.on('finish', () => {
      this.dispatchers.delete(guildId);
      this.botStatusUpdate(guildId, BotStatus.IN_VOICE_CHANNEL);
      onEnd();
    });
    dispatcher.on('stop', () => {
      dispatcher.destroy();
      this.dispatchers.delete(guildId);
      onEnd(true);
    });
    this.dispatchers.set(guildId, dispatcher);
  }

  private bindToEvents() {
    this.client.on('ready', () => {
      this.logger.log('Bot ready !');
      this.client.guilds.cache.forEach(g => {
        this.botStatusUpdate(g.id, BotStatus.CONNECTED);
      });
    });
    this.client.on('debug', (m) => {
      this.logger.debug(m);
    });
    this.client.on('error', (m) => {
      this.logger.error(m);
    });
    this.client.on('message', async (message: Message) => {
      if (!message.guild) return;
      const command = message.content;
      switch (command) {
        case Command.JOIN:
          await this.onJoinCommand(message);
          break;
        case Command.LEAVE:
          this.onLeaveCommand(message);
          break;
      }
    });
  }

  private botStatusUpdate(guildId: string, status: BotStatus): void {
    this.logger.log(`Bot status: ${status}`);
    this.onStatusChangeListeners.forEach(cb => cb(guildId, status));
    this.propagateState(guildId);
  }

  private async onJoinCommand(message: Message) {
    const channel = message.member.voice.channel;
    if (channel) {
      await this.joinChannel(channel);
    } else {
      await message.reply('You need to join a voice channel first!');
    }
  }

  private async joinChannel(channel: VoiceChannel): Promise<void> {
    if (channel) {
      try {
        await channel.join();
        this.logger.log(`Bot connected in channel ${channel.name} (${channel.id})`);
        this.botStatusUpdate(channel.guild.id, BotStatus.IN_VOICE_CHANNEL);
      } catch (e) {
        this.logger.error(e.message);
      }
    }
  }

  private onLeaveCommand(message: Message) {
    const voice = message.guild.voice;
    if (voice) {
      this.logger.debug(`Bot disconnected from channel ${voice.channel.name} (${voice.channel.id})`);
      this.botStatusUpdate(message.guild.id, BotStatus.CONNECTED);
      voice.connection.disconnect();
    }
  }

  private connect() {
    if (!this.token) {
      this.logger.error('Please define BOT_TOKEN env variable');
      return;
    }
    this.client.login(this.token);
  }

  private getVolume(guildId: string): number {
    return this.volumes.has(guildId) ? this.volumes.get(guildId) : this.DEFAULT_VOLUME;
  }

  private propagateState(guildId: string) {
    this.botGateway.sendStateUpdate(guildId, this.getState(guildId));
  }
}
