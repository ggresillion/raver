import {Injectable, Logger} from '@nestjs/common';
import {Client, Message, VoiceConnection} from 'discord.js';
import {Command} from './command.enum';
import {StorageService} from '../storage/storage.service';

@Injectable()
export class BotService {

  private readonly logger = new Logger(BotService.name);
  private readonly token = process.env.BOT_TOKEN;
  private client;
  private connections: VoiceConnection[] = [];

  constructor(
    private readonly storageService: StorageService,
  ) {
    this.client = new Client();
    this.connect();
    this.bindToEvents();
  }

  public playFile(uuid: string) {
    this.connections.forEach(co => {
      co.playFile(this.storageService.getPathFromUUID(uuid));
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
          if (message.member.voiceChannel) {
            message.member.voiceChannel.join()
              .then(connection => {
                this.connections.push(connection);
              })
              .catch((err) => this.logger.error(err.message));
          } else {
            message.reply('You need to join a voice channel first!');
          }
          break;
      }
    });
  }

  private connect() {
    if (!this.token) {
      this.logger.error('Please define BOT_TOKEN env variable');
    }
    this.client.login(this.token);
  }
}
