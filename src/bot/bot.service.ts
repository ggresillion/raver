import {Injectable, Logger, OnApplicationShutdown} from '@nestjs/common';
import {Client, Message} from 'discord.js';
import {Command} from './command.enum';
import {StorageService} from '../storage/storage.service';
import {BotGateway} from './bot.gateway';

export enum BotStatus {
  PLAYING = 'playing',
  IN_VOICE_CHANNEL = 'inVoiceChannel',
  CONNECTED = 'connected',
  DISCONNECTED = 'disconnected',
}

@Injectable()
export class BotService implements OnApplicationShutdown {

  private readonly logger = new Logger(BotService.name);
  private readonly token = process.env.BOT_TOKEN;
  private client: Client;

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

  public getInfos() {
    if (!this.client.user) {
      return {status: BotStatus.DISCONNECTED};
    }
    return {
      status: this.client.voice.connections.size > 0 ? BotStatus.IN_VOICE_CHANNEL : BotStatus.CONNECTED,
      id: this.client.user.id,
      username: this.client.user.username,
      avatar: this.client.user.avatar,
    };
  }

  public playFile(uuid: string) {
    this.client.voice.connections.forEach(co => {
      const dispacher = co.play(this.storageService.getPathFromUUID(uuid));
      dispacher.on('debug', this.logger.debug);
      dispacher.on('error', this.logger.error);
      dispacher.on('start', () => {
        this.logger.debug(`Playing file ` + uuid);
        this.botGateway.sendStatusUpdate(BotStatus.PLAYING);
      });
      dispacher.on('end', () => {
        this.botGateway.sendStatusUpdate(BotStatus.IN_VOICE_CHANNEL);
      });
    });
  }

  public playFromStream(stream: any, onStart: () => void, onEnd: () => void) {
    onStart();
    setTimeout(onEnd, 1000);
  }

  private bindToEvents() {
    this.client.on('ready', () => {
      this.logger.log('Bot ready !');
      this.botGateway.sendStatusUpdate(BotStatus.CONNECTED);
    });
    this.client.on('message', async (message: Message) => {
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

  private async onJoinCommand(message: Message) {
    const channel = message.member.voice.channel;
    if (channel) {
      await channel.join();
      try {
        this.logger.debug(`Bot connected in channel ${channel.name} (${channel.id})`);
        this.botGateway.sendStatusUpdate(BotStatus.IN_VOICE_CHANNEL);
      } catch (e) {
        this.logger.error(e.message);
      }
    } else {
      await message.reply('You need to join a voice channel first!');
    }
  }

  private onLeaveCommand(message: Message) {
    const voice = message.guild.voiceConnection;
    if (voice) {
      this.logger.debug(`Bot disconnected from channel ${voice.channel.name} (${voice.channel.id})`);
      this.botGateway.sendStatusUpdate(BotStatus.CONNECTED);
      voice.disconnect();
    }
  }

  private connect() {
    if (!this.token) {
      this.logger.error('Please define BOT_TOKEN env variable');
    }
    this.client.login(this.token);
  }
}
