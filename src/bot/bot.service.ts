import {Injectable, Logger} from '@nestjs/common';
import {Client, Message, VoiceConnection} from 'discord.js';
import {Command} from './command.enum';

@Injectable()
export class BotService {

  private readonly logger = new Logger(BotService.name);
  private readonly token = process.env.BOT_SECRET;
  private client;
  private connections: VoiceConnection[] = [];

  constructor() {
    this.client = new Client();
    this.connect();
    this.bindToEvents();
  }

  public playSound(id: number) {
    // TODO
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
              .catch(this.logger.error);
          } else {
            message.reply('You need to join a voice channel first!');
          }
          break;
      }
    });
  }

  private connect() {
    if (!this.token) {
      this.logger.error('Please define BOT_SECRET env variable');
    }
    this.client.login(this.token);
  }
}
