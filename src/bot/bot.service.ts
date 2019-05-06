import {Injectable, Logger, OnApplicationShutdown} from '@nestjs/common';
import {Client, Message} from 'discord.js';
import {Command} from './command.enum';
import {StorageService} from '../storage/storage.service';

export enum BotStatus {
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
      dispacher.on('start', () => this.logger.debug(`Playing file ` + uuid));
    });
  }

  private bindToEvents() {
    this.client.on('ready', () => {
      this.logger.log('Bot ready !');
    });
    this.client.on('message', (message: Message) => {
      const command = message.content;
      switch (command) {
        case Command.JOIN:
          this.onJoinCommand(message);
          break;
        case Command.LEAVE:
          this.onLeaveCommand(message);
          break;
      }
    });
  }

  private onJoinCommand(message: Message) {
    const channel = message.member.voice.channel;
    if (channel) {
      channel.join()
        .then(() => {
          this.logger.debug(`Bot connected in channel ${channel.name} (${channel.id})`);
        })
        .catch((err) => this.logger.error(err.message));
    } else {
      message.reply('You need to join a voice channel first!');
    }
  }

  private onLeaveCommand(message: Message) {
    const voice = message.guild.voiceConnection;
    if (voice) {
      this.logger.debug(`Bot disconnected from channel ${voice.channel.name} (${voice.channel.id})`);
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
