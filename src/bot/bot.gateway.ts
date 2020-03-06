import {OnGatewayConnection, WebSocketGateway, WebSocketServer} from '@nestjs/websockets';
import {Server} from 'socket.io';
import {ServerEvents} from './dto/server-events.enum';
import {BotStateDTO} from './dto/bot-state.dto';
import {forwardRef, Inject} from '@nestjs/common';
import {YoutubeService} from '../youtube/youtube.service';
import {BotService} from './bot.service';

@WebSocketGateway({namespace: 'bot'})
export class BotGateway implements OnGatewayConnection {

  @WebSocketServer()
  private server: Server;

  constructor(
    @Inject(forwardRef(() => BotService))
    private readonly botService: BotService,
  ) {
  }

  public sendStateUpdate(botState: BotStateDTO) {
    this.server.emit(ServerEvents.STATE_UPDATE.toString(), botState);
  }

  handleConnection(client: any, ...args: any[]): any {
    client.emit(this.botService.getState());
  }
}
