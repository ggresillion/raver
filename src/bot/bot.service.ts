import {Injectable, Logger, OnApplicationShutdown} from '@nestjs/common';
import {Client, Message, StreamDispatcher} from 'discord.js';
import {Command} from './command.enum';
import {StorageService} from '../storage/storage.service';
import {BotGateway} from './bot.gateway';
import {Readable} from 'stream';
import {UserDTO} from '../user/dto/user.dto';
import {GuildDTO} from '../guild/dto/guild.dto';
import {BotStateDTO} from './dto/bot-state.dto';
import {BotStatus} from './dto/bot-status.enum';

@Injectable()
export class BotService implements OnApplicationShutdown {

  private readonly logger = new Logger(BotService.name);
  private readonly token = process.env.BOT_TOKEN;
  private client: Client;
  private dispatcher: StreamDispatcher;
  private onStatusChangeListeners = [];

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
    const connection = this.client.voice.connections
      .find(c => c.channel.guild.id === guildId);
    this.dispatcher = connection.play(this.storageService.getPathFromUUID(uuid));
    // Workaround to prevent stream to end unexpectedly
    (this.dispatcher as any)._setSpeaking(1);
    // tslint:disable-next-line:no-empty
    (this.dispatcher as any)._setSpeaking = () => {
    };
    // --
    this.dispatcher.on('debug', this.logger.debug);
    this.dispatcher.on('error', this.logger.error);
    this.dispatcher.on('start', () => {
      this.logger.log(`Playing file ` + uuid);
      this.botStatusUpdate(BotStatus.PLAYING);
    });
    this.dispatcher.on('end', () => {
      this.dispatcher = null;
      this.botStatusUpdate(BotStatus.IN_VOICE_CHANNEL);
    });
  }

  public playFromStream(stream: Readable, onStart: () => void, onEnd: () => void) {
    this.client.voice.connections.forEach(co => {
      if (this.dispatcher) {
        this.dispatcher.end();
        this.dispatcher = null;
      }
      this.dispatcher = co.play(stream, {type: 'opus'});
      // Workaround to prevent stream to end unexpectedly
      (this.dispatcher as any)._setSpeaking(1);
      // tslint:disable-next-line:no-empty
      (this.dispatcher as any)._setSpeaking = () => {
      };
      // --
      this.dispatcher.on('debug', this.logger.debug);
      this.dispatcher.on('error', this.logger.error);
      this.dispatcher.on('start', () => {
        this.logger.log('Playing from stream');
        onStart();
      });
      this.dispatcher.on('end', () => {
        this.botStatusUpdate(BotStatus.IN_VOICE_CHANNEL);
        this.dispatcher = null;
        onEnd();
      });
    });
  }

  public pauseStream() {
    this.dispatcher.pause();
    this.logger.log('Paused stream');
  }

  public resumeStream() {
    this.dispatcher.resume();
    this.logger.log('Resumed stream');
  }

  public stopStream() {
    this.dispatcher.end();
    this.dispatcher = null;
    this.logger.log('Stopped stream');
  }

  public isPlaying(): boolean {
    return !!this.dispatcher;
  }

  public onBotStatusUpdate(guildId: string, cb: (status: BotStatus) => void) {
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

  private botStatusUpdate(status: BotStatus): void {
    this.logger.log(`Bot status: ${status}`);
    this.onStatusChangeListeners.forEach(cb => cb(status));
    const state: BotStateDTO = {
      guilds: this.client.guilds.cache.map(g => ({
        id: g.id,
        status: this.client.voice.connections.some(c => c.channel.guild.id === g.id)
          ? BotStatus.IN_VOICE_CHANNEL
          : BotStatus.CONNECTED,
      })),
    };
    this.botGateway.sendStateUpdate(state);
  }

  public getState(): BotStateDTO {
    return {
      guilds: this.client.guilds.cache.map(g => ({
        id: g.id,
        status: this.client.voice.connections.some(c => c.channel.guild.id === g.id)
          ? BotStatus.IN_VOICE_CHANNEL
          : BotStatus.CONNECTED,
      })),
    };
  }

  private bindToEvents() {
    this.client.on('ready', () => {
      this.logger.log('Bot ready !');
      this.botStatusUpdate(BotStatus.CONNECTED);
    });
    this.client.on('debug', (m) => {
      this.logger.debug(m);
    });
    this.client.on('error', (m) => {
      this.logger.error(m);
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
      try {
        await channel.join();
        this.logger.log(`Bot connected in channel ${channel.name} (${channel.id})`);
        this.botStatusUpdate(BotStatus.IN_VOICE_CHANNEL);
      } catch (e) {
        this.logger.error(e.message);
      }
    } else {
      await message.reply('You need to join a voice channel first!');
    }
  }

  private onLeaveCommand(message: Message) {
    const voice = message.guild.voice;
    if (voice) {
      this.logger.debug(`Bot disconnected from channel ${voice.channel.name} (${voice.channel.id})`);
      this.botStatusUpdate(BotStatus.CONNECTED);
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
}
